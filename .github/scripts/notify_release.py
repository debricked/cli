#!/usr/bin/env python3
"""Posts a "new release" notification to Slack and/or Microsoft Teams for a
just-published debricked/cli GitHub release.

Run by .github/workflows/notify-release.yml on every `release: published`
event. Mirrors the formatting of the internal-tools/deploy-notifier-cli Rust
tool: a "What's Changed" pull request list plus a "Full Changelog" link, both
taken from GitHub's own release-notes generator.
"""

import json
import os
import re
import sys
import urllib.error
import urllib.request

GITHUB_API = "https://api.github.com"

# PR bullet lines look like:
# "* <title> by @<author> in https://github.com/<owner>/<repo>/pull/<number>"
PR_LINE_RE = re.compile(
    r"^\*\s+(?P<title>.+?)\s+by\s+@[\w-]+\s+in\s+"
    r"(?P<url>https://github\.com/\S+/pull/(?P<number>\d+))\s*$",
    re.MULTILINE,
)
CHANGELOG_RE = re.compile(r"\*\*Full Changelog\*\*:\s*(\S+)")


def env(name: str) -> str:
    value = os.environ.get(name, "").strip()
    if not value:
        print(f"::error::{name} is required", file=sys.stderr)
        sys.exit(1)
    return value


def http_json(method: str, url: str, headers: dict, body: dict | None = None) -> dict:
    data = json.dumps(body).encode() if body is not None else None
    request = urllib.request.Request(url, data=data, method=method, headers=headers)
    try:
        with urllib.request.urlopen(request) as response:
            return json.load(response)
    except urllib.error.HTTPError as error:
        detail = error.read().decode(errors="replace")
        print(f"::error::{method} {url} returned {error.code}: {detail}", file=sys.stderr)
        raise


def generate_release_notes(repo: str, token: str, tag: str) -> dict:
    # GitHub auto-detects the previous tag itself; this is the same content
    # shown by the "Generate release notes" button in the GitHub UI.
    return http_json(
        "POST",
        f"{GITHUB_API}/repos/{repo}/releases/generate-notes",
        headers={
            "Accept": "application/vnd.github+json",
            "Authorization": f"Bearer {token}",
            "X-GitHub-Api-Version": "2022-11-28",
        },
        body={"tag_name": tag},
    )


def parse_pull_requests(body: str) -> list[dict]:
    return [
        {"number": int(m["number"]), "title": m["title"], "url": m["url"]}
        for m in (match.groupdict() for match in PR_LINE_RE.finditer(body))
    ]


def parse_changelog_url(body: str) -> str | None:
    match = CHANGELOG_RE.search(body)
    return match.group(1) if match else None


def slack_payload(tag: str, release_url: str, prs: list[dict], changelog_url: str | None) -> dict:
    if prs:
        pr_text = "".join(f"\u2022 <{pr['url']}|{pr['title']}>\n" for pr in prs)
    else:
        pr_text = "_No pull requests in this release_\n"
    if changelog_url:
        pr_text += f"\n<{changelog_url}|Full Changelog>\n"

    return {
        "text": f"New Debricked CLI release: {tag}!",
        "blocks": [
            {
                "type": "section",
                "text": {
                    "type": "mrkdwn",
                    "text": (
                        f":sparkles: *New Debricked CLI release: {tag}* :sparkles:\n"
                        f"<{release_url}|View release on GitHub>\n"
                    ),
                },
            },
            {
                "type": "section",
                "text": {"type": "mrkdwn", "text": f"*What's Changed*\n{pr_text}"},
            },
        ],
    }


def teams_payload(tag: str, release_url: str, prs: list[dict], changelog_url: str | None) -> dict:
    if prs:
        pr_text = "\n".join(f"- [{pr['title']}]({pr['url']})" for pr in prs)
    else:
        pr_text = "_No pull requests in this release_"
    if changelog_url:
        pr_text += f"\n\n[Full Changelog]({changelog_url})"

    pr_word = "pull request" if len(prs) == 1 else "pull requests"

    body = [
        {
            "type": "TextBlock",
            "text": f"\u2728 New Debricked CLI Release: {tag} \u2728",
            "wrap": True,
            "weight": "Bolder",
            "size": "Large",
        },
        {"type": "TextBlock", "text": "<at>channel</at>", "wrap": True},
        {
            "type": "TextBlock",
            "text": f"This release includes {len(prs)} {pr_word}. [View release on GitHub]({release_url})",
            "wrap": True,
            "size": "Small",
            "isSubtle": True,
        },
        {
            "type": "TextBlock",
            "text": "What's Changed",
            "wrap": True,
            "weight": "Bolder",
            "size": "Medium",
            "separator": True,
        },
        {"type": "TextBlock", "text": pr_text, "wrap": True},
    ]

    return {
        "type": "message",
        "attachments": [
            {
                "contentType": "application/vnd.microsoft.card.adaptive",
                "contentUrl": None,
                "content": {
                    "$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
                    "type": "AdaptiveCard",
                    "version": "1.4",
                    "body": body,
                    "msteams": {
                        "entities": [
                            {
                                "type": "mention",
                                "text": "<at>channel</at>",
                                "mentioned": {"id": "channel", "name": "channel"},
                            }
                        ]
                    },
                },
            }
        ],
    }


def post_webhook(url: str, payload: dict) -> None:
    http_json("POST", url, headers={"Content-Type": "application/json"}, body=payload)


def main() -> None:
    repo = env("GITHUB_REPOSITORY")
    token = env("GH_TOKEN")
    tag = env("RELEASE_TAG")
    release_url = env("RELEASE_URL")
    slack_webhook = os.environ.get("RELEASE_BOT_SLACK_WEBHOOK", "").strip()
    teams_webhook = os.environ.get("RELEASE_BOT_TEAMS_WEBHOOK", "").strip()

    if not slack_webhook and not teams_webhook:
        print(
            "::error::configure the RELEASE_BOT_SLACK_WEBHOOK or RELEASE_BOT_TEAMS_WEBHOOK repository secret",
            file=sys.stderr,
        )
        sys.exit(1)

    notes = generate_release_notes(repo, token, tag)
    prs = parse_pull_requests(notes["body"])
    changelog_url = parse_changelog_url(notes["body"])
    print(f"Found {len(prs)} pull request(s) for release {tag}")

    if slack_webhook:
        post_webhook(slack_webhook, slack_payload(tag, release_url, prs, changelog_url))
        print("Posted to Slack")
    if teams_webhook:
        post_webhook(teams_webhook, teams_payload(tag, release_url, prs, changelog_url))
        print("Posted to Teams")


if __name__ == "__main__":
    main()
