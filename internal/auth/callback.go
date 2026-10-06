package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/pkg/browser"
)

type IAuthWebHelper interface {
	Callback(string) string
	OpenURL(string) error
}

type AuthWebHelper struct {
	ServeMux *http.ServeMux
	openURL  func(string) error
}

func NewAuthWebHelper() AuthWebHelper {
	mux := http.NewServeMux()

	return AuthWebHelper{
		ServeMux: mux,
	}
}

func (awh AuthWebHelper) Callback(state string) string {
	code := make(chan string)
	defer close(code)

	awh.ServeMux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "Invalid state", http.StatusBadRequest)

			return
		}

		code <- r.URL.Query().Get("code")
		_, _ = fmt.Fprintf(w, "Authentication successful! You can close this window now.")
	})

	server := &http.Server{
		Addr:              ":9096",
		ReadHeaderTimeout: time.Minute,
		Handler:           awh.ServeMux,
	}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()
	defer func() {
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("HTTP server shutdown error: %v", err)
		}
	}()
	authCode := <-code // Wait for the authorization code

	return authCode
}

func (awh AuthWebHelper) OpenURL(authURL string) error {
	if awh.openURL != nil {
		return awh.openURL(authURL)
	}
	return browser.OpenURL(authURL)
}

func (awh AuthWebHelper) Login(ctx context.Context, authURL, state string) (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:9096")
	if err != nil {
		return "", fmt.Errorf("start OAuth callback listener: %w", err)
	}
	defer listener.Close()
	codes := make(chan string, 1)
	failures := make(chan error, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}
		if r.URL.Query().Get("error") != "" {
			select {
			case failures <- errors.New("browser authorization was denied"):
			default:
			}
			http.Error(w, "Authorization was denied", http.StatusBadRequest)
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Authorization code is missing", http.StatusBadRequest)
			return
		}
		select {
		case codes <- code:
			_, _ = fmt.Fprint(w, "Authorization received. You can close this window.")
		default:
			http.Error(w, "Authorization already received", http.StatusConflict)
		}
	})
	server := &http.Server{Handler: mux, ReadHeaderTimeout: time.Minute}
	defer server.Close()
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			select {
			case failures <- err:
			default:
			}
		}
	}()
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if err := awh.OpenURL(authURL); err != nil {
		return "", fmt.Errorf("open browser: %w", err)
	}
	select {
	case code := <-codes:
		return code, nil
	case err := <-failures:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
