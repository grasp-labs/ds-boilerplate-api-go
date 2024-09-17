package models

// ServerStatus represents the status of the server
// swagger:model
type ServerStatus struct {
	Details string `json:"details"`
	Time    string `json:"time"`
}

type ServerVersion struct {
	Version string `json:"version"`
}
