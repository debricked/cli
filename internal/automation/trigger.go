package automation

type TriggerEvent struct {
	Dependency     string   `json:"dependency"`
	DependencyLink string   `json:"dependencyLink"`
	Licenses       []string `json:"licenses"`
	Cve            string   `json:"cve"`
	Cvss2          float32  `json:"cvss2"`
	Cvss3          float32  `json:"cvss3"`
	CveLink        string   `json:"cveLink"`
}
