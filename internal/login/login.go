package login

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type IAuthenticator interface {
	Authenticate() (string, error)
}

type Authenticator struct {
	ClientID string
	Scopes   []string
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func generateRandomString(length int) string {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func createCodeChallenge(codeVerifier string) string {
	// Create a SHA-256 hash of the code verifier
	hash := sha256.Sum256([]byte(codeVerifier))

	// Encode the hash to base64
	encoded := base64.StdEncoding.EncodeToString(hash[:])

	// Make it URL safe
	encoded = strings.TrimRight(encoded, "=")
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")

	return encoded
}

func (a Authenticator) Authenticate() (string, error) {
	// Set up OAuth2 configuration
	config := &oauth2.Config{
		ClientID:     a.ClientID,
		ClientSecret: "",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://debricked.com/app/oauth/authorize",
			TokenURL: "https://debricked.com/app/oauth/token",
		},
		RedirectURL: "http://localhost:9096/callback",
		Scopes:      a.Scopes,
	}

	// Create a random state
	state := generateRandomString(8)
	codeVerifier := generateRandomString(64)

	// Generate the authorization URL
	authURL := config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge", createCodeChallenge(codeVerifier)),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	// Start a temporary HTTP server to handle the callback
	code := make(chan string)
	defer close(code)
	server := &http.Server{Addr: ":9096"}
	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Ensure the server is shut down when we're done
	defer server.Shutdown(context.Background())

	// Open the browser for the user to log in
	err := openBrowser(authURL)
	if err != nil {
		log.Fatal("Could not open browser:", err)
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		code <- r.URL.Query().Get("code")
		fmt.Fprintf(w, "Authentication successful! You can close this window now.")
	})

	// Wait for the authorization code
	authCode := <-code

	// Exchange the authorization code for a token
	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
