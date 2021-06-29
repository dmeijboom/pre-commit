package config

type When struct {
	Glob string `json:"glob"`
}

type Check struct {
	Name string `json:"name"`
	Cmd  string `json:"cmd"`
	When []When `json:"when"`
}
