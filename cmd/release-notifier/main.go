package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	githubAPI      = "https://api.github.com"
	maxBlockLength = 2900
)

var (
	pullRequestLine = regexp.MustCompile(`(?m)^\*\s+(.+?)\s+by\s+@[\w-]+\s+in\s+(https://github\.com/\S+/pull/\d+)\s*$`)
	changelogLine   = regexp.MustCompile(`\*\*Full Changelog\*\*:\s*(\S+)`)
)

type release struct {
	Assets []releaseAsset `json:"assets"`
	Body   string         `json:"body"`
}

type releaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

type pullRequest struct {
	Title string
	URL   string
}

type asset struct {
	Name string
	URL  string
	Size string
}

type slackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type slackBlock struct {
	Type string    `json:"type"`
	Text slackText `json:"text"`
}

type slackPayload struct {
	Text   string       `json:"text"`
	Blocks []slackBlock `json:"blocks"`
}

type teamsPayload struct {
	Type        string            `json:"type"`
	Attachments []teamsAttachment `json:"attachments"`
}

type teamsAttachment struct {
	ContentType string    `json:"contentType"`
	ContentURL  *string   `json:"contentUrl"`
	Content     teamsCard `json:"content"`
}

type teamsCard struct {
	Schema  string        `json:"$schema"`
	Type    string        `json:"type"`
	Version string        `json:"version"`
	Body    []interface{} `json:"body"`
}

func main() {
	repo := requiredEnv("GITHUB_REPOSITORY")
	token := requiredEnv("GH_TOKEN")
	tag := requiredEnv("RELEASE_TAG")
	releaseURL := requiredEnv("RELEASE_URL")
	slackWebhook := strings.TrimSpace(os.Getenv("RELEASE_BOT_SLACK_WEBHOOK"))
	teamsWebhook := strings.TrimSpace(os.Getenv("RELEASE_BOT_TEAMS_WEBHOOK"))

	if slackWebhook == "" && teamsWebhook == "" {
		fail("configure the RELEASE_BOT_SLACK_WEBHOOK or RELEASE_BOT_TEAMS_WEBHOOK repository secret")
	}

	release := getRelease(repo, token, tag)
	prs := parsePullRequests(release.Body)
	assets := parseAssets(release.Assets)
	changelogURL := parseChangelogURL(release.Body)
	fmt.Printf("Found %d pull request(s) and %d asset(s) for release %s\n", len(prs), len(assets), tag)

	if slackWebhook != "" {
		postWebhook(slackWebhook, slackMessage(tag, releaseURL, prs, changelogURL, assets))
		fmt.Println("Posted to Slack")
	}
	if teamsWebhook != "" {
		postWebhook(teamsWebhook, teamsMessage(tag, releaseURL, prs, changelogURL, assets))
		fmt.Println("Posted to Teams")
	}
}

func requiredEnv(name string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		fail(name + " is required")
	}

	return value
}

func getRelease(repo, token, tag string) release {
	var release release
	githubRequest(http.MethodGet, fmt.Sprintf("%s/repos/%s/releases/tags/%s", githubAPI, repo, tag), token, nil, &release)

	return release
}

func githubRequest(method, url, token string, body interface{}, target interface{}) {
	var requestBody io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			fail(fmt.Sprintf("encode GitHub request: %v", err))
		}
		requestBody = bytes.NewReader(encoded)
	}

	request, err := http.NewRequest(method, url, requestBody) //nolint:gosec // GitHub API URL is built from the configured repository.
	if err != nil {
		fail(fmt.Sprintf("create GitHub request: %v", err))
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Content-Type", "application/json")

	response := doRequest(request)
	defer func() {
		if err := response.Body.Close(); err != nil {
			fail(fmt.Sprintf("close GitHub response: %v", err))
		}
	}()

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		fail(fmt.Sprintf("decode GitHub response: %v", err))
	}
}

func parsePullRequests(body string) []pullRequest {
	matches := pullRequestLine.FindAllStringSubmatch(body, -1)
	prs := make([]pullRequest, 0, len(matches))
	for _, match := range matches {
		prs = append(prs, pullRequest{Title: match[1], URL: match[2]})
	}

	return prs
}

func parseChangelogURL(body string) string {
	match := changelogLine.FindStringSubmatch(body)
	if len(match) == 0 {
		return ""
	}

	return match[1]
}

func parseAssets(releaseAssets []releaseAsset) []asset {
	assets := make([]asset, 0, len(releaseAssets))
	for _, releaseAsset := range releaseAssets {
		assets = append(assets, asset{
			Name: releaseAsset.Name,
			URL:  releaseAsset.BrowserDownloadURL,
			Size: humanSize(releaseAsset.Size),
		})
	}

	return assets
}

func humanSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d Bytes", size)
	}
	units := []string{"KB", "MB", "GB"}
	value := float64(size)
	for _, unit := range units {
		value /= 1024
		if value < 1024 || unit == "GB" {
			return fmt.Sprintf("%.2f %s", value, unit)
		}
	}

	return fmt.Sprintf("%.2f GB", value)
}

func slackMessage(tag, releaseURL string, prs []pullRequest, changelogURL string, assets []asset) slackPayload {
	changed := "_No pull requests in this release_\n"
	if len(prs) > 0 {
		var lines strings.Builder
		for _, pr := range prs {
			fmt.Fprintf(&lines, "• <%s|%s>\n", pr.URL, pr.Title)
		}
		changed = lines.String()
	}
	if changelogURL != "" {
		changed += fmt.Sprintf("\n<%s|Full Changelog>\n", changelogURL)
	}

	blocks := []slackBlock{
		section(fmt.Sprintf(":sparkles: *New Debricked CLI release: %s* :sparkles:\n<%s|View release on GitHub>\n", tag, releaseURL)),
		section("*What's Changed*\n" + changed),
	}
	blocks = append(blocks, slackAssetBlocks(assets)...)

	return slackPayload{Text: fmt.Sprintf("New Debricked CLI release: %s!", tag), Blocks: blocks}
}

func section(text string) slackBlock {
	return slackBlock{Type: "section", Text: slackText{Type: "mrkdwn", Text: text}}
}

func slackAssetBlocks(assets []asset) []slackBlock {
	if len(assets) == 0 {
		return []slackBlock{section("*Assets (0)*\n_No assets attached to this release_\n")}
	}

	heading := fmt.Sprintf("*Assets (%d)*\n", len(assets))
	current := heading
	blocks := []slackBlock{}
	for _, asset := range assets {
		line := fmt.Sprintf("• <%s|%s> (%s)\n", asset.URL, asset.Name, asset.Size)
		if len(current)+len(line) > maxBlockLength && current != heading {
			blocks = append(blocks, section(current))
			current = ""
		}
		current += line
	}

	return append(blocks, section(current))
}

func teamsMessage(tag, releaseURL string, prs []pullRequest, changelogURL string, assets []asset) teamsPayload {
	changed := "_No pull requests in this release_"
	if len(prs) > 0 {
		lines := make([]string, 0, len(prs))
		for _, pr := range prs {
			lines = append(lines, fmt.Sprintf("- [%s](%s)", pr.Title, pr.URL))
		}
		changed = strings.Join(lines, "\n")
	}
	if changelogURL != "" {
		changed += fmt.Sprintf("\n\n[Full Changelog](%s)", changelogURL)
	}

	assetText := "_No assets attached to this release_"
	if len(assets) > 0 {
		lines := make([]string, 0, len(assets))
		for _, asset := range assets {
			lines = append(lines, fmt.Sprintf("- [%s](%s) (%s)", asset.Name, asset.URL, asset.Size))
		}
		assetText = strings.Join(lines, "\n")
	}

	prWord := "pull requests"
	if len(prs) == 1 {
		prWord = "pull request"
	}
	body := []interface{}{
		textBlock(fmt.Sprintf("✨ New Debricked CLI Release: %s ✨", tag), "Bolder", "Large", false),
		map[string]interface{}{"type": "TextBlock", "text": fmt.Sprintf("This release includes %d %s. [View release on GitHub](%s)", len(prs), prWord, releaseURL), "wrap": true, "size": "Small", "isSubtle": true},
		textBlock("What's Changed", "Bolder", "Medium", true),
		textBlock(changed, "", "", false),
		textBlock(fmt.Sprintf("Assets (%d)", len(assets)), "Bolder", "Medium", true),
		textBlock(assetText, "", "", false),
	}

	return teamsPayload{
		Type: "message",
		Attachments: []teamsAttachment{{
			ContentType: "application/vnd.microsoft.card.adaptive",
			Content: teamsCard{
				Schema:  "http://adaptivecards.io/schemas/adaptive-card.json",
				Type:    "AdaptiveCard",
				Version: "1.4",
				Body:    body,
			},
		}},
	}
}

func textBlock(text, weight, size string, separator bool) map[string]interface{} {
	block := map[string]interface{}{"type": "TextBlock", "text": text, "wrap": true}
	if weight != "" {
		block["weight"] = weight
	}
	if size != "" {
		block["size"] = size
	}
	if separator {
		block["separator"] = true
	}

	return block
}

func postWebhook(url string, payload interface{}) {
	body, err := json.Marshal(payload)
	if err != nil {
		fail(fmt.Sprintf("encode webhook payload: %v", err))
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body)) //nolint:gosec // Webhook URL is supplied through a repository secret.
	if err != nil {
		fail(fmt.Sprintf("create webhook request: %v", err))
	}
	request.Header.Set("Content-Type", "application/json")
	response := doRequest(request)
	if err := response.Body.Close(); err != nil {
		fail(fmt.Sprintf("close webhook response: %v", err))
	}
}

func doRequest(request *http.Request) *http.Response {
	response, err := http.DefaultClient.Do(request) //nolint:gosec // URLs are the configured GitHub API or webhook endpoints.
	if err != nil {
		fail(fmt.Sprintf("%s %s failed: %v", request.Method, request.URL, err))
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(response.Body)
		if err := response.Body.Close(); err != nil {
			fail(fmt.Sprintf("close error response: %v", err))
		}

		fail(fmt.Sprintf("%s %s returned HTTP %d: %s", request.Method, request.URL, response.StatusCode, string(body)))
	}

	return response
}

func fail(message string) {
	fmt.Fprintln(os.Stderr, "::error::"+message)
	os.Exit(1)
}
