package storage

import "MiniProject/internal/types"

type Storage interface {
	RegisterDevice(device types.Device) (int64, error)
	UpdateDevice(device types.Device) (int64, error)
	GetDevice(id int64) ([]types.DeviceInfo, error)
	GetDeviceList() ([]types.DeviceInfo, error)
	DeleteDevice(id int64) (int64, error)
}
