package auth

import (
	"context"
	"fmt"
	"github.com/debricked/cli/internal/client"
	"github.com/zalando/go-keyring"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
)

type IAuthenticator interface {
	Authenticate() error
	Logout() error
	Token() (*oauth2.Token, error)
}

type Authenticator struct {
	ClientID     string
	Scopes       []string
	Client       client.IDebClient
	SecretClient ISecretClient
}

type ISecretClient interface {
	Set(string, string) error
	Get(string) (string, error)
	Delete(string) error
}

type DebrickedSecretClient struct {
	User string
}

func (dsc DebrickedSecretClient) Set(service, secret string) error {
	return keyring.Set(service, dsc.User, secret)
}

func (dsc DebrickedSecretClient) Get(service string) (string, error) {
	return keyring.Get(service, dsc.User)
}

func (dsc DebrickedSecretClient) Delete(service string) error {
	return keyring.Delete(service, dsc.User)
}

func NewDebrickedAuthenticator(client client.IDebClient) Authenticator {
	return Authenticator{
		ClientID: "01919462-7d6e-78e8-aa24-ba779213c90f",
		Scopes:   []string{"select", "profile", "basicRepo"},
		Client:   client,
		SecretClient: DebrickedSecretClient{
			User: "DebrickedCLI",
		},
	}
}

func (a Authenticator) Logout() error {
	err := a.SecretClient.Delete("DebrickedRefreshToken")
	if err != nil {
		return err
	}
	err = a.SecretClient.Delete("DebrickedAccessToken")
	return err
}

func (a Authenticator) Token() (*oauth2.Token, error) {
	refreshToken, err := a.SecretClient.Get("DebrickedRefreshToken")
	if err != nil {
		return nil, err
	}
	accessToken, err := a.SecretClient.Get("DebrickedAccessToken")
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{
		RefreshToken: refreshToken,
		TokenType:    "jwt",
		AccessToken:  accessToken,
	}, nil
}

func (a Authenticator) Authenticate() error {
	// Set up OAuth2 configuration
	config := &oauth2.Config{
		ClientID:     a.ClientID,
		ClientSecret: "",
		Endpoint: oauth2.Endpoint{
			AuthURL:  a.Client.Host() + "/app/oauth/authorize",
			TokenURL: a.Client.Host() + "/app/oauth/token",
		},
		RedirectURL: "http://localhost:9096/callback",
		Scopes:      a.Scopes,
	}

	// Create a random state
	state := oauth2.GenerateVerifier()
	codeVerifier := oauth2.GenerateVerifier()

	// Generate the authorization URL
	authURL := config.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(codeVerifier),
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
	token, err := config.Exchange(
		context.Background(),
		authCode,
		oauth2.SetAuthURLParam("client_id", a.ClientID),
		oauth2.VerifierOption(codeVerifier),
	)
	if err != nil {
		return err
	}
	a.SecretClient.Set("DebrickedRefreshToken", token.RefreshToken)
	a.SecretClient.Set("DebrickedAccessToken", token.AccessToken)
	return nil
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
