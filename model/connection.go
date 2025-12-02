package model

import "time"

type Connection struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Host          string    `json:"host"`
	Domain        string    `json:"domain"`
	Username      string    `json:"username"`
	StorePassword bool      `json:"storePassword"`
	Password      string    `json:"password,omitempty"`
	ExtraArgs     []string  `json:"extraArgs,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
