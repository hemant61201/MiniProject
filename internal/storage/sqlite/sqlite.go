package sqlite

import (
	"MiniProject/internal/config"
	"MiniProject/internal/types"
	"database/sql"

	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

func NewSqlite(config *config.Config) (*Sqlite, error) {

	db, err := sql.Open("sqlite", config.StoragePath)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS devices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    deviceType VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    ipAddress VARCHAR(255) NOT NULL,
    osType VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{Db: db}, nil
}

func (sqlite *Sqlite) RegisterDevice(device *types.Device) (int64, error) {

	stmt, err := sqlite.Db.Prepare("INSERT INTO devices (name, deviceType, status, ipAddress, osType) VALUES (?, ?, ?, ?, ?)")

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(device.Name, device.DeviceType, device.Status, device.IpAddress, device.OsType)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (sqlite *Sqlite) GetDeviceList() ([]types.DeviceInfo, error) {

	rows, err := sqlite.Db.Query("SELECT * FROM devices")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var devices []types.DeviceInfo

	for rows.Next() {
		var device types.DeviceInfo

		err := rows.Scan(&device.Id, &device.Name, &device.DeviceType, &device.Status, &device.IpAddress, &device.OsType, &device.CreatedAt, &device.UpdatedAt)

		if err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func (sqlite *Sqlite) GetDevice(id int64) ([]types.DeviceInfo, error) {

	stmt, err := sqlite.Db.Prepare("SELECT * FROM devices WHERE id = ?")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var devices []types.DeviceInfo

	for rows.Next() {
		var device types.DeviceInfo

		err := rows.Scan(&device.Id, &device.Name, &device.DeviceType, &device.Status, &device.IpAddress, &device.OsType, &device.CreatedAt, &device.UpdatedAt)

		if err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func (sqlite *Sqlite) CheckDevice(id int64) (bool, error) {

	stmt, err := sqlite.Db.Prepare("SELECT id FROM devices WHERE id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	var exists int64

	err = stmt.QueryRow(id).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, err
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (sqlite *Sqlite) UpdateDevice(id int64, deviceInput *types.UpdateDeviceInput) (int64, error) {

	var device types.Device

	err := sqlite.Db.QueryRow("SELECT name, deviceType, status, ipAddress, osType FROM devices WHERE id = ?", id).
		Scan(&device.Name, &device.DeviceType, &device.Status, &device.IpAddress, &device.OsType)

	if err == sql.ErrNoRows {
		return 0, err
	} else if err != nil {
		return 0, err
	}

	if deviceInput.Name != nil {
		device.Name = *deviceInput.Name
	}
	if deviceInput.DeviceType != nil {
		device.DeviceType = *deviceInput.DeviceType
	}
	if deviceInput.IPAddress != nil {
		device.IpAddress = *deviceInput.IPAddress
	}
	if deviceInput.Status != nil {
		device.Status = *deviceInput.Status
	}
	if deviceInput.OSType != nil {
		device.OsType = *deviceInput.OSType
	}

	_, err = sqlite.Db.Exec(`UPDATE devices SET name=?, deviceType=?, status=?, ipAddress=?, osType=?, updatedAt=CURRENT_TIMESTAMP WHERE id=?`,
		device.Name, device.DeviceType, device.Status, device.IpAddress, device.OsType, id)

	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (sqlite *Sqlite) DeleteDevice(id int64) (int64, error) {

	result, err := sqlite.Db.Exec("DELETE FROM devices WHERE id = ?", id)

	if err != nil {
		return 0, err
	}

	lastId, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return lastId, nil
}
