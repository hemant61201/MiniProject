package types

type Device struct {
	Name       string `json:"name"`
	DeviceType string `json:"deviceType"`
	IpAddress  string `json:"ipAddress"`
	Status     string `json:"status,omitempty"`
	OsType     string `json:"osType"`
}

type UpdateDeviceInput struct {
	Name       *string `json:"name,omitempty"`
	DeviceType *string `json:"device_type,omitempty"`
	IPAddress  *string `json:"ip_address,omitempty"`
	Status     *string `json:"status,omitempty"`
	OSType     *string `json:"os_type,omitempty"`
}

type DeviceInfo struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	DeviceType string `json:"deviceType"`
	IpAddress  string `json:"ipAddress"`
	Status     string `json:"status"`
	OsType     string `json:"osType"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	Metadata   struct {
		CPU    string `json:"cpu"`
		Memory string `json:"memory"`
		Disk   string `json:"disk"`
	} `json:"metadata,omitempty"`
}
