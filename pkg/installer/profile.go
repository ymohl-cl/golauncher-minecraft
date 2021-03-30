package installer

import "time"

// Profile json data
type Profile struct {
	Icon      string    `json:"icon"`
	LastUsed  time.Time `json:"lastUsed"`
	VersionID string    `json:"lastVersionId"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
}

// LauncherProfile json data
type LauncherProfile struct {
	Profiles map[string]Profile `json:"profiles"`
}
