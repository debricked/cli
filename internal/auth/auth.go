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
	SecretClient ISecretClient
	OAuthConfig  *oauth2.Config
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
		SecretClient: DebrickedSecretClient{
			User: "DebrickedCLI",
		},
		OAuthConfig: &oauth2.Config{
			ClientID:     "01919462-7d6e-78e8-aa24-ba779213c90f",
			ClientSecret: "",
			Endpoint: oauth2.Endpoint{
				AuthURL:  client.Host() + "/app/oauth/authorize",
				TokenURL: client.Host() + "/app/oauth/token",
			},
			RedirectURL: "http://localhost:9096/callback",
			Scopes:      []string{"select", "profile", "basicRepo"},
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

func (a Authenticator) callback(state string) string {
	code := make(chan string)
	defer close(code)
	server := &http.Server{Addr: ":9096"} // Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	defer server.Shutdown(
		context.Background(),
	) // Ensure the server is shut down when we're done
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		code <- r.URL.Query().Get("code")
		fmt.Fprintf(w, "Authentication successful! You can close this window now.")
	})
	authCode := <-code // Wait for the authorization code

	return authCode
}

func (a Authenticator) exchange(authCode, codeVerifier string) (*oauth2.Token, error) {
	return a.OAuthConfig.Exchange(
		context.Background(),
		authCode,
		oauth2.SetAuthURLParam("client_id", a.OAuthConfig.ClientID),
		oauth2.VerifierOption(codeVerifier),
	)

}

func (a Authenticator) Authenticate() error {
	state := oauth2.GenerateVerifier()
	codeVerifier := oauth2.GenerateVerifier()

	authURL := a.OAuthConfig.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(codeVerifier),
	)

	err := openBrowser(authURL)
	if err != nil {
		log.Fatal("Could not open browser:", err)
	}

	authCode := a.callback(state)
	token, err := a.exchange(authCode, codeVerifier)
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
