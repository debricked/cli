package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskValue(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		expected string
	}{
		{"caret range", "^4.18.2", "^4.18.2"},
		{"compound range", ">=1.0.0 <2.0.0 || 3.x", ">=1.0.0 <2.0.0 || 3.x"},
		{"star", "*", "*"},
		{"dist-tag", "latest", "latest"},
		{"berry npm range", "npm:^1.4.0", "npm:^1.4.0"},
		{"berry npm alias", "npm:@scope/other@^1.0.0", "npm:@scope/other@^1.0.0"},

		{"npmjs tarball", "https://registry.npmjs.org/express/-/express-4.18.2.tgz", "https://registry.npmjs.org/express/-/express-4.18.2.tgz"},
		{"yarnpkg tarball", "https://registry.yarnpkg.com/jest/-/jest-29.7.0.tgz#cccc", "https://registry.yarnpkg.com/jest/-/jest-29.7.0.tgz#cccc"},
		{"yarnpkg without scheme", "registry.yarnpkg.com/jest/-/jest-29.7.0.tgz", "registry.yarnpkg.com/jest/-/jest-29.7.0.tgz"},
		{"registry with credentials", "https://ci:s3cret@registry.npmjs.org/x.tgz", MaskedURL},
		{"lookalike registry host", "https://registry.npmjs.org.acme.internal/x.tgz", MaskedURL},

		{"git word-only repo", "git+ssh://git@git.acme.internal/platform/authclient.git#v1.0.0", MaskedGitWord},
		{"git hyphenated repo", "git+ssh://git@git.acme.internal/platform/auth-client.git#v2.1.0", MaskedGitOther},
		{"git dotted repo", "git+ssh://git@git.acme.internal/platform/auth.client.git", MaskedGitOther},
		{"scp style", "git@git.acme.internal:platform/authclient.git", MaskedGitWord},
		{"no git user, word-only repo", "git+https://git.acme.internal/platform/authclient.git", MaskedGitWordNoUser},
		{"no git user, hyphenated repo", "git+https://git.acme.internal/platform/auth-client.git", MaskedGitOtherNoUser},
		{"other ssh user", "ssh://jane@git.acme.internal/platform/authclient.git", MaskedGitWordNoUser},
		{"https git url", "https://git.acme.internal/platform/authclient.git#main", MaskedGitWordNoUser},
		{"github protocol", "github:acme/auth-client", MaskedGitOtherNoUser},
		{"github shorthand", "acme/authclient#feature/PROJ-123", MaskedGitWordNoUser},

		{"file with dot", "file:../shared-utils", MaskedLocalDot},
		{"link with dot", "link:./x", MaskedLocalDot},
		{"portal with dot", "portal:../x", MaskedLocalDot},
		{"exec with dot", "exec:./x", MaskedLocalDot},
		{"file without dot", "file:shared-utils", MaskedLocal},
		{"link without dot", "link:shared-utils", MaskedLocal},
		{"absolute file path", "file:/home/jane/x", MaskedLocal},
		{"bare relative path", "../shared-utils", MaskedLocalDot},

		{"other https url", "https://artifactory.acme.internal/api/logger-1.4.2.tgz", MaskedURL},
		{"other http url", "http://files.acme.internal/logger.tgz", MaskedURL},

		{"workspace star", "workspace:*", "workspace:*"},
		{"workspace caret", "workspace:^", "workspace:^"},
		{"workspace tilde", "workspace:~", "workspace:~"},
		{"workspace dot", "workspace:.", "workspace:."},
		{"workspace path", "workspace:packages/shared-utils", MaskedWorkspace},

		{"patch file", "patch:foo@npm%3A1.0.0#./.yarn/patches/foo.patch", MaskedPatchFile},
		{"yarn builtin patch", "patch:typescript@npm%3A5.0.0#optional!builtin<compat/typescript>", "patch:typescript@npm%3A5.0.0#optional!builtin<compat/typescript>"},
		{"yarn builtin patch naming the project", "patch:x@npm%3A1#optional!builtin<compat/x>::version=1&hash=a&locator=%40acme%2Fportal%40workspace%3A.", "patch:x@npm%3A1#optional!builtin<compat/x>::version=1&hash=a"},
		{"other patch", "patch:foo@npm%3A1.0.0#./internal/fix", MaskedPatch},

		{"unknown protocol", "svn://svn.acme.internal/repo", MaskedUnknown},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, MaskValue(c.value))
		})
	}
}

func TestMatchesUnsupportedProtocol(t *testing.T) {
	assert.True(t, MatchesUnsupportedProtocol("git+ssh://git@git.acme.internal/platform/authclient.git#v1"))
	assert.True(t, MatchesUnsupportedProtocol("file:../x"))
	assert.True(t, MatchesUnsupportedProtocol("patch:foo#./x.patch"))

	assert.False(t, MatchesUnsupportedProtocol("git+ssh://git@git.acme.internal/platform/auth-client.git#v1"))
	assert.False(t, MatchesUnsupportedProtocol("file:x"))
	assert.False(t, MatchesUnsupportedProtocol("https://artifactory.acme.internal/x.tgz"))
	assert.False(t, MatchesUnsupportedProtocol("https://registry.npmjs.org/x/-/x-1.0.0.tgz"))
	assert.False(t, MatchesUnsupportedProtocol("../x"))
}
