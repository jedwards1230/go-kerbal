package datacollector

// CKAN Spec: https://github.com/KSP-CKAN/CKAN/blob/master/Spec.md

type Ckan struct {
	SpecVersion   string                 `json:"spec_version"`
	Identifier    string                 `json:"identifier"`
	Name          string                 `json:"name"`
	Abstract      string                 `json:"abstract"`
	Author        string                 `json:"author"`
	Download      string                 `json:"download"`
	License       string                 `json:"license"`
	Version       string                 `json:"version"`
	Epoch         string                 `json:"epoch"`
	VersionKspMax string                 `json:"ksp_version_max"`
	VersionKspMin string                 `json:"ksp_version_min"`
	Resources     resource               `json:"resources"`
	Tags          map[string]interface{} `json:"tags"`
	Depends       map[string]interface{} `json:"depends"`
	Conflicts     map[string]interface{} `json:"conflicts"`
}

type resource struct {
	Homepage    string `json:"homepage"`
	Spacedock   string `json:"spacedock"`
	Repository  string `json:"repository"`
	XScreenshot string `json:"x_screenshot"`
}
