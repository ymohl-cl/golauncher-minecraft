package main

import "time"

type Profile struct {
	CreationDate time.Time
	Icon         string
	LastUsed     time.Time
	VersionID    string
	Type         string
}

type LauncherProfile struct {
	Profiles map[string]Profile
}

// Exist check if the profil exist in the profiles list
func (l LauncherProfile) Exist(profileName string) bool {
	for key := range l.Profiles {
		if key == profileName {
			return true
		}
	}
	return false
}
