package datacollector

type Ckan struct {
	SpecVersion string `json:"spec_version"`
	Identifier  string `json:"identifier"`
	Name        string `json:"name"`
	Abstract    string `json:"abstract"`
	Author      string `json:"author"`
	Download    string `json:"download"`
	License     string `json:"license"`
	Version     string `json:"version"`
}
