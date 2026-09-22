package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseReleaseContent(t *testing.T) {
	body := "## What's Changed\n* SPM support by @octocat in https://github.com/debricked/cli/pull/336\n\n**Full Changelog**: https://github.com/debricked/cli/compare/v26.3.4...v26.3.5\n"

	prs := parsePullRequests(body)
	if len(prs) != 1 {
		t.Fatalf("expected one pull request, got %d", len(prs))
	}
	if prs[0].Title != "SPM support" || prs[0].URL != "https://github.com/debricked/cli/pull/336" {
		t.Fatalf("unexpected pull request: %#v", prs[0])
	}
	if got := parseChangelogURL(body); got != "https://github.com/debricked/cli/compare/v26.3.4...v26.3.5" {
		t.Fatalf("unexpected changelog URL: %s", got)
	}
}

func TestAssetPayloadsIncludeCountAndSizes(t *testing.T) {
	assets := parseAssets([]releaseAsset{
		{Name: "checksums.txt", BrowserDownloadURL: "https://example.test/checksums.txt", Size: 612},
		{Name: "cli_linux_amd64.apk", BrowserDownloadURL: "https://example.test/cli_linux_amd64.apk", Size: 31000000},
	})
	if assets[0].Size != "612 Bytes" || assets[1].Size != "29.56 MB" {
		t.Fatalf("unexpected asset sizes: %#v", assets)
	}

	slack := slackMessage("v26.3.5", "https://example.test/release", nil, "", assets)
	if !strings.Contains(slack.Blocks[2].Text.Text, "*Assets (2)*") {
		t.Fatalf("Slack payload omits asset count: %s", slack.Blocks[2].Text.Text)
	}
	if _, err := json.Marshal(slack); err != nil {
		t.Fatalf("Slack payload does not serialize: %v", err)
	}

	teams := teamsMessage("v26.3.5", "https://example.test/release", nil, "", assets)
	body := teams.Attachments[0].Content.Body
	assetsHeading := body[4].(map[string]interface{})["text"]
	if assetsHeading != "Assets (2)" {
		t.Fatalf("Teams payload has unexpected asset heading: %v", assetsHeading)
	}
	encoded, err := json.Marshal(teams)
	if err != nil {
		t.Fatalf("Teams payload does not serialize: %v", err)
	}
	if strings.Contains(string(encoded), "<at>") || strings.Contains(string(encoded), "msteams") {
		t.Fatalf("Teams payload contains an unsupported placeholder mention: %s", encoded)
	}
}
