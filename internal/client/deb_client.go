package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/debricked/cli/internal/auth"

	"github.com/fatih/color"
)

const DefaultDebrickedUri = "https://debricked.com"
const DefaultTimeout = 15
const enterpriseCheckUri = "/api/1.0/open/user-profile/get-billing-info"

type IDebClient interface {
	// Post makes a POST request to one of Debricked's API endpoints
	Post(uri string, contentType string, body *bytes.Buffer, timeout int) (*http.Response, error)
	// Get makes a GET request to one of Debricked's API endpoints
	Get(uri string, format string) (*http.Response, error)
	SetAccessToken(accessToken *string)
	IsEnterpriseCustomer(silent bool) bool
	Host() string
	Authenticator() auth.IAuthenticator
}

type DebClient struct {
	host          *string
	httpClient    IClient
	accessToken   *string
	jwtToken      string
	authenticator auth.IAuthenticator
}

func NewDebClient(accessToken *string, httpClient IClient) *DebClient {
	host := os.Getenv("DEBRICKED_URI")
	if len(host) == 0 {
		host = DefaultDebrickedUri
	}
	authenticator := auth.NewDebrickedAuthenticator(host)

	return &DebClient{
		host:          &host,
		httpClient:    httpClient,
		accessToken:   initAccessToken(accessToken),
		jwtToken:      "",
		authenticator: authenticator,
	}
}

func (debClient *DebClient) Host() string {
	return *debClient.host
}

func (debClient *DebClient) Post(uri string, contentType string, body *bytes.Buffer, timeout int) (*http.Response, error) {
	if timeout > 0 {
		return postWithTimeout(uri, debClient, contentType, body, true, timeout)
	}

	return post(uri, debClient, contentType, body, true)
}

func (debClient *DebClient) Get(uri string, format string) (*http.Response, error) {
	return get(uri, debClient, true, format)
}

func (debClient *DebClient) SetAccessToken(accessToken *string) {
	debClient.accessToken = initAccessToken(accessToken)
}

func (debClient *DebClient) Authenticator() auth.IAuthenticator {
	return debClient.authenticator
}

func initAccessToken(accessToken *string) *string {
	if accessToken == nil {
		accessToken = new(string)
	}
	if len(*accessToken) == 0 {
		*accessToken = os.Getenv("DEBRICKED_TOKEN")
	}

	if len(*accessToken) == 0 {
		return nil
	}

	return accessToken
}

type BillingPlan struct {
	SCA    string `json:"sca"`
	Select string `json:"select"`
}

func printNonEnterpriseMessage(specificError string, finalMessage string, silent bool) {
	if !silent {
		fmt.Print(
			color.YellowString("⚠️"),
			" Could not validate enterprise billing plan due to ",
			specificError,
			". File fingerprint will not be run or analyzed, since ",
			"it requires a verified enterprise SCA subscription.",
			finalMessage,
		)
	}
}

func (debClient *DebClient) IsEnterpriseCustomer(silent bool) bool {
	res, err := debClient.Get(enterpriseCheckUri, "application/json")
	if err != nil {
		printNonEnterpriseMessage("connection error", "If this issue persists please create an issue on github: https://github.com/debricked/cli/issues\n", silent)

		return false
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		printNonEnterpriseMessage("HTTP error", "If this issue persists please create an issue on github: https://github.com/debricked/cli/issues\n", silent)

		return false
	}

	billingPlanJSON, err := io.ReadAll(res.Body)
	if err != nil {
		printNonEnterpriseMessage("response JSON formatting error", "If this issue persists please create an issue on github: https://github.com/debricked/cli/issues\n", silent)

		return false
	}

	var billingPlan BillingPlan

	err = json.Unmarshal(billingPlanJSON, &billingPlan)
	if err != nil {
		printNonEnterpriseMessage("malformed response", "If this issue persists please create an issue on github: https://github.com/debricked/cli/issues\n", silent)

		return false
	}

	if billingPlan.SCA != "enterprise" {
		response := "billing plan currently being \"" + string(billingPlan.SCA) + "\""
		final := "To upgrade your plan visit: " + color.BlueString("https://debricked.com/app/en/repositories?billingModal=enterprise,free") + "\n"
		printNonEnterpriseMessage(response, final, silent)

		return false
	}

	return true
}
