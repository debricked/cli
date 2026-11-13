package pnpm

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/debricked/cli/internal/sanitize/pm/testdata"
	"github.com/debricked/cli/internal/sanitize/pm/util"
)

const pnpmV6Before = `lockfileVersion: '6.0'

settings:
  autoInstallPeers: true
  excludeLinksFromLockfile: false

importers:

  .:
    dependencies:
      '@acme/auth-client':
        specifier: git+ssh://git@git.acme.internal/platform/auth-client.git#v2.1.0
        version: git.acme.internal/platform/auth-client/9f3c1ab
      '@acme/logger':
        specifier: ^1.4.0
        version: 1.4.2
      '@acme/shared-utils':
        specifier: file:../shared-utils
        version: file:../shared-utils
      express:
        specifier: ^4.18.2
        version: 4.18.2
    devDependencies:
      jest:
        specifier: ^29.7.0
        version: 29.7.0

  packages/billing-worker:
    dependencies:
      qs:
        specifier: 6.11.0
        version: 6.11.0

packages:

  /@acme/logger@1.4.2:
    resolution: {integrity: sha512-AAAA, tarball: https://artifactory.acme.internal/api/npm/npm-local/@acme/logger/-/logger-1.4.2.tgz}
    dependencies:
      qs: 6.11.0
    dev: false

  /express@4.18.2:
    resolution: {integrity: sha512-BBBB}
    engines: {node: '>= 0.10.0'}
    dependencies:
      qs: 6.11.0
    dev: false

  /jest@29.7.0:
    resolution: {integrity: sha512-CCCC}
    hasBin: true
    dev: true

  /qs@6.11.0:
    resolution: {integrity: sha512-DDDD}
    engines: {node: '>=0.6'}
    dev: false

  git.acme.internal/platform/auth-client/9f3c1ab:
    resolution: {commit: 9f3c1ab, repo: 'git@git.acme.internal:platform/auth-client.git', type: git}
    name: '@acme/auth-client'
    version: 2.1.0
    dev: false

  file:../shared-utils:
    resolution: {directory: ../shared-utils, type: directory}
    name: '@acme/shared-utils'
    version: 1.0.0
    dev: false
`

const pnpmV6After = `lockfileVersion: '6.0'
importers:
  .:
    dependencies:
      '@acme/auth-client': {}
      '@acme/logger': {}
      '@acme/shared-utils': {}
      express: {}
    devDependencies:
      jest: {}
  redacted-1:
    dependencies:
      qs: {}
packages:
  /@acme/logger@1.4.2:
    dependencies:
      qs: 6.11.0
  /express@4.18.2:
    dependencies:
      qs: 6.11.0
  /jest@29.7.0: {}
  /qs@6.11.0: {}
  redacted-1:
    name: '@acme/auth-client'
    version: 2.1.0
  'file:./redacted-2':
    name: '@acme/shared-utils'
    version: 1.0.0
`

const pnpmV9Before = `lockfileVersion: '9.0'

settings:
  autoInstallPeers: true
  excludeLinksFromLockfile: false

importers:

  .:
    dependencies:
      '@acme/auth-client':
        specifier: git+ssh://git@git.acme.internal/platform/auth-client.git#v2.1.0
        version: git+ssh://git@git.acme.internal/platform/auth-client.git#9f3c1ab
      '@acme/logger':
        specifier: ^1.4.0
        version: 1.4.2
      '@acme/shared-utils':
        specifier: file:../shared-utils
        version: file:../shared-utils
      express:
        specifier: ^4.18.2
        version: 4.18.2
    devDependencies:
      jest:
        specifier: ^29.7.0
        version: 29.7.0

packages:

  '@acme/auth-client@git+ssh://git@git.acme.internal/platform/auth-client.git#9f3c1ab':
    resolution: {commit: 9f3c1ab, repo: 'git@git.acme.internal:platform/auth-client.git', type: git}
    version: 2.1.0

  '@acme/logger@1.4.2':
    resolution: {integrity: sha512-AAAA, tarball: https://artifactory.acme.internal/api/npm/npm-local/@acme/logger/-/logger-1.4.2.tgz}

  '@acme/shared-utils@file:../shared-utils':
    resolution: {directory: ../shared-utils, type: directory}

  express@4.18.2:
    resolution: {integrity: sha512-BBBB}
    engines: {node: '>= 0.10.0'}

  jest@29.7.0:
    resolution: {integrity: sha512-CCCC}
    hasBin: true

  qs@6.11.0:
    resolution: {integrity: sha512-DDDD}
    engines: {node: '>=0.6'}

snapshots:

  '@acme/auth-client@git+ssh://git@git.acme.internal/platform/auth-client.git#9f3c1ab': {}

  '@acme/logger@1.4.2':
    dependencies:
      qs: 6.11.0

  '@acme/shared-utils@file:../shared-utils': {}

  express@4.18.2:
    dependencies:
      qs: 6.11.0

  jest@29.7.0: {}

  qs@6.11.0: {}
`

const pnpmV9After = `lockfileVersion: '9.0'
importers:
  .:
    dependencies:
      '@acme/auth-client': {}
      '@acme/logger': {}
      '@acme/shared-utils': {}
      express: {}
    devDependencies:
      jest: {}
packages:
  '@acme/auth-client@git+ssh://git@redacted/redacted-repo.git#redacted':
    version: 2.1.0
  '@acme/logger@1.4.2': {}
  '@acme/shared-utils@file:./redacted': {}
  express@4.18.2: {}
  jest@29.7.0: {}
  qs@6.11.0: {}
snapshots:
  '@acme/auth-client@git+ssh://git@redacted/redacted-repo.git#redacted': {}
  '@acme/logger@1.4.2':
    dependencies:
      qs: 6.11.0
  '@acme/shared-utils@file:./redacted': {}
  express@4.18.2:
    dependencies:
      qs: 6.11.0
  jest@29.7.0: {}
  qs@6.11.0: {}
`

func TestSanitizePnpmLock(t *testing.T) {
	cases := map[string]struct{ before, after string }{
		"v6": {pnpmV6Before, pnpmV6After},
		"v9": {pnpmV9Before, pnpmV9After},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			sanitized, err := sanitizePnpmLock([]byte(c.before), newTestContext("pnpm-lock.yaml"))

			assert.NoError(t, err)
			assert.Equal(t, c.after, string(sanitized))
		})
	}
}

func TestSanitizePnpmLockReferencesGetTheSamePlaceholderAsTheirKeys(t *testing.T) {
	v6 := `lockfileVersion: '6.0'
packages:
  /app@1.0.0:
    dependencies:
      '@acme/auth': git.acme.internal/platform/auth/9f3c1ab
      '@acme/utils': file:../utils
      qs: 6.11.0(react@18.2.0)
  git.acme.internal/platform/auth/9f3c1ab:
    name: '@acme/auth'
    version: 1.0.0
  file:../utils:
    name: '@acme/utils'
    version: 1.0.0
`
	sanitized, err := sanitizePnpmLock([]byte(v6), newTestContext("pnpm-lock.yaml"))

	assert.NoError(t, err)
	assert.Equal(t, `lockfileVersion: '6.0'
packages:
  /app@1.0.0:
    dependencies:
      '@acme/auth': redacted-1
      '@acme/utils': 'file:./redacted-2'
      qs: '6.11.0(react@18.2.0)'
  redacted-1:
    name: '@acme/auth'
    version: 1.0.0
  'file:./redacted-2':
    name: '@acme/utils'
    version: 1.0.0
`, string(sanitized))
}

func TestSanitizePnpmLockSingleProjectFromPnpm8(t *testing.T) {
	v6 := `lockfileVersion: '6.0'

settings:
  autoInstallPeers: true

dependencies:
  is-odd:
    specifier: github:acme/is-odd#3.0.1
    version: git@github.com+acme/is-odd/a80ee0d
  qs:
    specifier: ^6.11.0
    version: 6.11.0

devDependencies:
  jest:
    specifier: ^29.7.0
    version: 29.7.0

packages:

  /app@1.0.0:
    dependencies:
      is-odd: git@github.com+acme/is-odd/a80ee0d

  git@github.com+acme/is-odd/a80ee0d:
    resolution: {tarball: https://codeload.github.com/acme/is-odd/tar.gz/a80ee0d}
    name: is-odd
    version: 3.0.1
`
	sanitized, err := sanitizePnpmLock([]byte(v6), newTestContext("pnpm-lock.yaml"))

	assert.NoError(t, err)
	assert.Equal(t, `lockfileVersion: '6.0'
dependencies:
  is-odd: {}
  qs: {}
devDependencies:
  jest: {}
packages:
  /app@1.0.0:
    dependencies:
      is-odd: redacted-1
  redacted-1:
    name: is-odd
    version: 3.0.1
`, string(sanitized))
}

func TestSanitizePnpmLockMasksPeerSuffixes(t *testing.T) {
	v9 := `lockfileVersion: '9.0'
packages:
  plugin@1.0.0:
    peerDependencies:
      '@acme/base': '*'
  '@acme/base@file:../base': {}
snapshots:
  plugin@1.0.0(@acme/base@file:../base):
    dependencies:
      '@acme/base': file:../base
  '@acme/base@file:../base': {}
`
	sanitized, err := sanitizePnpmLock([]byte(v9), newTestContext("pnpm-lock.yaml"))

	assert.NoError(t, err)
	assert.NotContains(t, string(sanitized), "../base")
	assert.Contains(t, string(sanitized), "'plugin@1.0.0(@acme/base@file:./redacted)':")
	assert.Contains(t, string(sanitized), "'@acme/base': 'file:./redacted'")
}

func TestSanitizePnpmLockKeepsTwoSourcesOfOnePackageApart(t *testing.T) {
	v9 := `lockfileVersion: '9.0'
packages:
  'a@file:../one': {}
  'a@file:../two': {}
  'b@git+ssh://git@git.acme.internal/p/b-c.git#one': {}
  'b@git+ssh://git@git.acme.internal/p/b-c.git#two': {}
snapshots:
  'a@file:../two': {}
  'app@1.0.0':
    dependencies:
      a: file:../two
`
	sanitized, err := sanitizePnpmLock([]byte(v9), newTestContext("pnpm-lock.yaml"))

	assert.NoError(t, err)
	assert.Equal(t, `lockfileVersion: '9.0'
packages:
  'a@file:./redacted': {}
  'a@file:./redacted-2': {}
  'b@git+ssh://git@redacted/redacted-repo.git#redacted': {}
  'b@git+ssh://git@redacted/redacted-repo.git#redacted-2': {}
snapshots:
  'a@file:./redacted-2': {}
  app@1.0.0:
    dependencies:
      a: 'file:./redacted-2'
`, string(sanitized))
}

func TestSanitizePnpmLockKeepsVersionsAsStrings(t *testing.T) {
	v6 := `lockfileVersion: 5.4
packages:
  /a/1.10:
    version: 1.10
    dependencies:
      b: 1.10
      c: 2
      d: "1.2.3"
`
	sanitized, err := sanitizePnpmLock([]byte(v6), newTestContext("pnpm-lock.yaml"))
	assert.NoError(t, err)

	assert.Equal(t, `lockfileVersion: 5.4
packages:
  /a/1.10:
    version: '1.10'
    dependencies:
      b: '1.10'
      c: '2'
      d: 1.2.3
`, string(sanitized))

	var parsed struct {
		Packages map[string]map[string]any `yaml:"packages"`
	}
	assert.NoError(t, yaml.Unmarshal(sanitized, &parsed))
	assert.Equal(t, "1.10", parsed.Packages["/a/1.10"]["version"])
	assert.Equal(t, map[string]any{"b": "1.10", "c": "2", "d": "1.2.3"}, parsed.Packages["/a/1.10"]["dependencies"])
}

func TestSanitizePnpmLockWritesLongKeysAsPlainKeys(t *testing.T) {
	long := "plugin@1.0.0(" + strings.Repeat("peer-with-a-long-name@1.0.0", 10) + ")"
	v9 := "lockfileVersion: '9.0'\nsnapshots:\n  '" + long + "': {}\n"

	sanitized, err := sanitizePnpmLock([]byte(v9), newTestContext("pnpm-lock.yaml"))

	assert.NoError(t, err)
	assert.NotContains(t, string(sanitized), "? ")
	var parsed struct {
		Snapshots map[string]any `yaml:"snapshots"`
	}
	assert.NoError(t, yaml.Unmarshal(sanitized, &parsed))
	assert.Contains(t, parsed.Snapshots, long)
}

func TestSanitizePnpmLockRemovesCommentsAndUnknownKeys(t *testing.T) {
	v9 := `# s3cret comment
lockfileVersion: '9.0' # s3cret
catalogs:
  default:
    s3cret: 1.0.0
packages:
  # s3cret
  qs@6.11.0:
    s3cret: yes
    deprecated: s3cret
`
	sanitized, err := sanitizePnpmLock([]byte(v9), newTestContext("pnpm-lock.yaml"))

	assert.NoError(t, err)
	assert.Equal(t, "lockfileVersion: '9.0'\npackages:\n  qs@6.11.0: {}\n", string(sanitized))
}

func TestSanitizePnpmLockRecordsPathsInTheMap(t *testing.T) {
	context, all := testdata.NewContext("pnpm-lock.yaml")

	_, err := sanitizePnpmLock([]byte(pnpmV6Before), context)

	assert.NoError(t, err)
	assert.Equal(t, []util.MapEntry{
		{Placeholder: "redacted-1", Name: "git.acme.internal/platform/auth-client/9f3c1ab", Path: "pnpm-lock.yaml"},
		{Placeholder: "file:./redacted-2", Name: "file:../shared-utils", Path: "pnpm-lock.yaml"},
		{Placeholder: "redacted-1", Name: "packages/billing-worker", Path: "pnpm-lock.yaml"},
	}, all.Entries())
}

func TestSanitizePnpmLockInvalid(t *testing.T) {
	invalid := map[string]string{
		"empty":            "",
		"not yaml":         "a: b: c: [",
		"not a mapping":    "- a\n",
		"unrelated yaml":   "name: portal\n",
		"packages is list": "packages:\n  - a\n",
		"unbalanced peers": "packages:\n  'a@1.0.0(b@1.0.0': {}\n",
	}
	for name, content := range invalid {
		t.Run(name, func(t *testing.T) {
			_, err := sanitizePnpmLock([]byte(content), newTestContext("pnpm-lock.yaml"))
			assert.ErrorIs(t, err, util.ErrInvalidFormat)
		})
	}
}

func newTestContext(path string) *util.FileContext {
	context, _ := testdata.NewContext(path)

	return context
}
