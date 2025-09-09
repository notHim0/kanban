package models

type User struct {
	Id string
	Name string
	Password string
}

type Project struct {
	Id string `json:"id,omitempty"`
	UserId string `json:"user,omitempty"`
	Name string `json:"name,omitempty"`
	RepoUrl string `json:"repo_url,omitempty"`
	SiteUrl string `json:"site_url,omitempty"`
	Description string `json:"description,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
	DevDependencies []string `json:"dev_dependencies,omitempty"`
	Status string `json:"status,omitempty"`
}