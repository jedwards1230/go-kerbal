package module

type ModuleVersion struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Identifier  string `json:"identifier"`
	Abstract    string `json:"abstract"`
	License     string `json:"license"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
}

type Module struct {
	Id       string
	Versions []ModuleVersion
}
