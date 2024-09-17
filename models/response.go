package models

type ServerStatus struct {
	Details string `json:"details"`
	Time    string `json:"time"`
}

type ServerVersion struct {
	Version string `json:"version"`
}
