package sqlx_client

import (
	"time"
)

type Device struct {
	Product  string `db:"product"`
	DeviceId string `db:"device_id"`
}

type DeviceDataByDay struct {
	DeviceId   string    `db:"device_id"`
	PowerQ     string    `db:"switches_power_q"`
	EnergyPt   string    `db:"switches_energy_pt"`
	RmsI       string    `db:"switches_rms_i"`
	RecordTime time.Time `db:"record_time"`
}

//查阅用户关联设备列表
func GetDeviceListByUsername(name string) ([]Device, error) {
	var ds []Device
	err := db.Select(&ds, "select product,device_id from user_bind_device where username=?", name)
	return ds, err
}

//按天查设备的数据明细
func GetDeviceDataByDay(deviceId string, year, month, day int) (*DeviceDataByDay, error) {
	var ddbd DeviceDataByDay
	err := db.Get(&ddbd, "select device_id,switches_power_q,switches_energy_pt,switches_rms_i,record_time from device_data where record_year=? and record_month=? and record_day=?", year, month, day)
	return &ddbd, err
}

//按月查设备的数据明细，按时间戳升序
func GetDeviceDataByMonth(deviceId string, year, month int) ([]DeviceDataByDay, error) {
	var ddbds []DeviceDataByDay
	err := db.Select(&ddbds, "select device_id,switches_power_q,switches_energy_pt,switches_rms_i,record_time from device_data where record_year=? and record_month=? order by record_time", year, month)
	return ddbds, err
}

//按年查设备的数据明细，按时间戳升序
func GetDeviceDataByYear(deviceId string, year int) ([]DeviceDataByDay, error) {
	var ddbds []DeviceDataByDay
	err := db.Select(&ddbds, "select device_id,switches_power_q,switches_energy_pt,switches_rms_i,record_time from device_data where record_year=? order by record_time", year)
	return ddbds, err
}
