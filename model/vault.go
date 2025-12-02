package model

type Vault struct {
	Version     int          `json:"version"`
	Connections []Connection `json:"connections"`
}
