package installer

// ConfigForge json
type ConfigForge struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

// ConfigServer json
type ConfigServer struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

// ConfigFTP json
type ConfigFTP struct {
	URL      string `json:"url"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Config json to installer
type Config struct {
	GameVersion string       `json:"minecraft_version"`
	Forge       ConfigForge  `json:"forge"`
	Server      ConfigServer `json:"server"`
	FTP         ConfigFTP    `json:"ftp"`
}
