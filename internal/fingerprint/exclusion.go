package fingerprint

type DefaultFingerprintExclusionList struct {
	Directories []string
	Files       []string
	Endings     []string
	Extensions  []string
}

var defaultFingerprintExclusions = DefaultFingerprintExclusionList{
	Directories: []string{
		".idea",
		"nbproject",
		"nbbuild",
		"nddist",
		"__pycache__",
		"_yardoc",
		"eggs",
		"wheels",
		"htmlcov",
		"__pypackages__",
		".git",
		"*.egg-info",
		"*venv",
		"*venv3",
		"node_modules",
		"vendor",
		".git",
		"obj",
		"bower_components",
		".m2",                 // Default Maven directory for settings.xml and dependencies
		".debrickedTmpFolder", // temporary debricked data
	},
	Files: []string{
		"gradlew", "gradlew.bat", "mvnw", "mvnw.cmd", "gradle-wrapper.jar", "maven-wrapper.jar",
		"thumbs.db", "babel.config.js", "license.txt", "license.md", "copying.lib", "makefile",
		"\\[content_types\\].xml", "\\[Content_Types\\].xml", "py.typed", "LICENSE.APACHE2", "LICENSE.MIT",
		".nycrc", ".jshintrc", ".tm_properties", "webpack.config.js", "eslint.config.js",
		"eslint.config.mjs", "eslint.config.cjs", ".eslintrc.yaml", ".eslintrc.js", ".eslintrc.mjs",
		".eslintrc.cjs", "jest.config.js", "vite.config.mjs", "vite.config.js",
	},
	Endings: []string{
		"-doc", "changelog", "config", "copying", "license", "authors", "news", "licenses", "notice",
		"readme", "swiftdoc", "texidoc", "todo", "version", "ignore", "manifest", "sqlite", "sqlite3",
		"nycrc", "targ", "eslintrc", "prettierrc", "resx",
	},
	Extensions: []string{
		".1", ".2", ".3", ".4", ".5", ".6", ".7", ".8", ".9", ".ac", ".adoc", ".am",
		".asciidoc", ".bmp", ".build", ".cfg", ".chm", ".cmake", ".cnf",
		".conf", ".config", ".contributors", ".copying", ".crt", ".csproj", ".css",
		".csv", ".dat", ".data", ".doc", ".docx", ".dtd", ".dts", ".iws", ".c9", ".c9revisions",
		".dtsi", ".dump", ".eot", ".eps", ".geojson", ".gdoc", ".gif",
		".glif", ".gmo", ".gradle", ".guess", ".hex", ".htm", ".html", ".ico", ".iml",
		".in", ".inc", ".info", ".ini", ".ipynb", ".jpeg", ".jpg", ".json", ".jsonld", ".lock",
		".log", ".m4", ".map", ".markdown", ".md", ".md5", ".meta", ".mk", ".mxml",
		".o", ".otf", ".out", ".pbtxt", ".pdf", ".pem", ".phtml", ".plist", ".png",
		".po", ".ppt", ".prefs", ".properties", ".pyc", ".qdoc", ".result", ".rgb",
		".rst", ".scss", ".sha", ".sha1", ".sha2", ".sha256", ".sln", ".spec", ".sql",
		".sub", ".svg", ".svn-base", ".tab", ".template", ".test", ".tex", ".tiff",
		".toml", ".transform", ".ttf", ".txt", ".utf-8", ".vim", ".wav", ".woff", ".woff2", ".xht",
		".xhtml", ".xls", ".xlsx", ".xpm", ".xsd", ".xul", ".yaml", ".yml", ".wfp",
		".editorconfig", ".dotcover", ".pid", ".lcov", ".egg", ".manifest", ".cache", ".coverage", ".cover",
		".gem", ".lst", ".pickle", ".pdb", ".gml", ".pot", ".plt", ".pyi",
	},
}

func DefaultExclusionsFingerprint() []string {
	var default_exclusions []string
	for _, excluded_dir := range defaultFingerprintExclusions.Directories {
		default_exclusions = append(default_exclusions, "**/"+excluded_dir+"/**")
	}
	for _, excluded_file := range defaultFingerprintExclusions.Files {
		default_exclusions = append(default_exclusions, "**/"+excluded_file)
	}
	for _, excluded_extension := range defaultFingerprintExclusions.Extensions {
		default_exclusions = append(default_exclusions, "**/*"+excluded_extension)
	}
	for _, excluded_ending := range defaultFingerprintExclusions.Endings {
		default_exclusions = append(default_exclusions, "**/*"+excluded_ending)
	}

	return default_exclusions
}

func DefaultInclusionsFingerprint() []string {

	return []string{
		"package.json",
	}
}
