package types

import "database/sql"

type Device struct {
	Name       string `json:"name"`
	DeviceType string `json:"deviceType"`
	IpAddress  string `json:"ipAddress"`
	Status     string `json:"status,omitempty"`
	OsType     string `json:"osType"`
}

type DeviceInfo struct {
	Id         int64          `json:"id"`
	Name       string         `json:"name"`
	DeviceType string         `json:"deviceType"`
	IpAddress  string         `json:"ipAddress"`
	Status     string         `json:"status"`
	OsType     string         `json:"osType"`
	CreatedAt  string         `json:"createdAt"`
	UpdatedAt  string         `json:"updatedAt"`
	Metadata   sql.NullString `json:"metadata"`
}
