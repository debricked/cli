package auth

import (
	"context"
	"fmt"
	"log"
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
		fmt.Fprintf(w, "Authentication successful! You can close this window now.")
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
	return browser.OpenURL(authURL)
}
