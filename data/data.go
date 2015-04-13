package data

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _resources_ddl_mysql_select_nearest_stations_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x74\x8d\x51\x0a\xc2\x30\x10\x44\xaf\xb2\x3f\xfe\x89\x27\x50\xaf\x22\x4b\x3b\x8d\xc1\x64\x02\x49\x44\xeb\xe9\xad\x15\xb7\x34\xd2\xbf\x79\xfb\x98\x9d\x82\x80\xae\x8a\x3a\xb0\x1b\x2f\x37\x8c\xfb\x5f\xf6\xbd\x45\x6a\x84\xc1\x3d\x07\xcb\xd5\x47\xbc\x12\x17\x19\x94\xce\x20\x7a\x4e\x87\xba\xb0\x3e\xd7\xfc\xf1\x89\x6b\x9f\x28\x43\x4e\x51\x5c\x1d\xca\x61\x16\x1e\x45\x1e\x57\x64\x34\x7f\xe5\x78\x92\x5d\x11\x65\xdf\x0c\xc8\xf9\x4f\x7c\x97\x36\x1a\x93\x98\x1b\xef\x00\x00\x00\xff\xff\xea\x44\x32\x64\x0c\x01\x00\x00")

func resources_ddl_mysql_select_nearest_stations_sql_bytes() ([]byte, error) {
	return bindata_read(
		_resources_ddl_mysql_select_nearest_stations_sql,
		"resources/ddl/mysql/select-nearest-stations.sql",
	)
}

func resources_ddl_mysql_select_nearest_stations_sql() (*asset, error) {
	bytes, err := resources_ddl_mysql_select_nearest_stations_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "resources/ddl/mysql/select-nearest-stations.sql", size: 268, mode: os.FileMode(420), modTime: time.Unix(1428879790, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _resources_ddl_mysql_select_stop_times_by_calendar_dates_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x64\x91\x41\x72\xc3\x20\x0c\x45\xf7\x39\x05\x1b\x4f\x36\x9d\xdc\xc0\x67\x61\x34\x48\x6e\xd4\x21\x92\x2b\xc9\x69\x7b\xfb\x9a\x78\x51\x4c\xd9\xf1\xf4\xf9\xfc\x0f\x4e\x95\x4a\xa4\x4b\xda\x97\xc7\x72\xf3\xd0\x35\x33\xbe\x9d\x81\xc0\x83\x06\x84\xe4\x65\x40\x15\x62\x24\x2a\x7f\xa4\x6a\x81\x60\x95\x1c\x3f\x6b\xe7\x06\x66\xfc\x84\x9a\x83\xfb\x3b\x90\x56\xb0\xd8\x8c\x06\xfe\xb2\x75\xfa\xdc\x48\x4a\x2f\x67\xdb\x6b\x34\xf3\x3e\xbb\xe9\x16\x94\xfd\xae\x16\x43\x85\x63\x72\x0e\x72\xb0\xa2\x55\xed\x9f\x90\xbe\x63\x9c\x84\x71\x7b\xa9\xcb\x62\xfa\x78\xc1\xf7\x58\x3c\x4f\x7e\x44\x6c\xa9\x3d\x2f\x5b\xad\x4d\x9c\x58\x84\x2c\x7d\x28\xcb\x49\x5a\xa0\x92\x20\x58\x46\x08\xf2\x54\x30\xa9\x1c\x35\xc9\x9e\x5c\x68\xf7\x9f\x0b\x76\xbb\xcb\xd7\x9d\x8c\xc6\xdf\x9a\x27\x4c\x20\xb8\x9f\xbf\x35\xa3\x34\xa7\xeb\xe4\xd7\xdf\x00\x00\x00\xff\xff\x89\xca\x20\x83\xdb\x01\x00\x00")

func resources_ddl_mysql_select_stop_times_by_calendar_dates_sql_bytes() ([]byte, error) {
	return bindata_read(
		_resources_ddl_mysql_select_stop_times_by_calendar_dates_sql,
		"resources/ddl/mysql/select-stop-times-by-calendar-dates.sql",
	)
}

func resources_ddl_mysql_select_stop_times_by_calendar_dates_sql() (*asset, error) {
	bytes, err := resources_ddl_mysql_select_stop_times_by_calendar_dates_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "resources/ddl/mysql/select-stop-times-by-calendar-dates.sql", size: 475, mode: os.FileMode(420), modTime: time.Unix(1428881643, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _resources_ddl_mysql_select_stop_times_by_calendars_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x64\x91\xc1\x4e\xc4\x30\x0c\x44\xef\xfb\x15\xbe\x54\x7b\x41\x48\xdc\x29\xbf\x12\x59\x89\xcb\x06\x65\xed\x62\xbb\x0b\xfc\x3d\xe9\xe6\xd2\xa4\xbd\xf5\xcd\x74\xec\x71\x8d\x0a\x45\xbf\x40\x7d\xcc\x97\x57\x73\x59\x43\x4e\x2f\xd0\x13\xc6\x3b\x8d\x2c\x91\xc5\x91\x15\xf4\x13\x12\x3e\xa0\x22\x11\x3d\x0b\x07\xff\x5b\x8f\x89\xa8\x9a\x1f\x58\x82\xe7\x6e\x50\xa2\x15\xd5\x37\xa5\x51\x78\x46\x1b\x7d\x6f\xc4\xb1\xfb\x20\x6b\xad\xb3\x0f\xe8\x4a\xa8\x6c\x4e\xc1\x6e\xa2\x3e\x76\x69\xd2\xb0\x4e\x83\x51\x8a\xe8\xd9\x4a\xbf\x7e\x92\x5c\xf3\x7e\x37\xb8\x2c\x2a\xf7\x27\xfd\xf4\xc5\xc2\x64\x6d\xd3\x7d\x7b\x0b\xcb\x56\xca\xee\x86\xcc\x4c\x0a\x5f\x92\x19\x3a\x6f\xc4\x42\x9c\x50\x0d\x22\x08\xb7\xa2\xa4\x8f\x1c\xa9\x66\xcf\xf1\xf0\x02\x97\x9f\x1b\x29\x8d\xff\x6d\x9e\x12\x20\xa7\x96\x5a\xfd\x5e\xaf\x17\x12\x3a\xc1\xfb\x0c\xd7\xc9\xae\x47\xb5\x8e\x6a\xda\xc7\x41\x9b\x6c\x7e\xfb\x0f\x00\x00\xff\xff\x34\x36\x86\x83\x14\x02\x00\x00")

func resources_ddl_mysql_select_stop_times_by_calendars_sql_bytes() ([]byte, error) {
	return bindata_read(
		_resources_ddl_mysql_select_stop_times_by_calendars_sql,
		"resources/ddl/mysql/select-stop-times-by-calendars.sql",
	)
}

func resources_ddl_mysql_select_stop_times_by_calendars_sql() (*asset, error) {
	bytes, err := resources_ddl_mysql_select_stop_times_by_calendars_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "resources/ddl/mysql/select-stop-times-by-calendars.sql", size: 532, mode: os.FileMode(420), modTime: time.Unix(1428881643, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _resources_ddl_mysql_select_stops_by_date_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x94\x8d\xcd\x0a\x02\x31\x0c\x84\x5f\x25\x97\x85\x5d\x59\x04\x0f\x1e\x04\xdf\xa5\xc4\x36\x6a\xa1\xdb\x94\x4e\x40\xf6\xed\xad\x7f\x15\xbc\x99\xd3\x4c\xf2\x65\x06\x92\xc4\x1b\x61\x0b\xd3\xe2\x62\x98\x3f\x32\xf3\x22\xdd\x04\x81\xef\x26\xb1\x7d\xb5\xe6\x87\x4e\xea\xd9\xa2\x66\x67\x6b\x69\x5f\xbb\x36\x87\x3d\x6d\x08\xe6\x42\x84\x71\xf6\x32\x16\x8d\xd9\xc6\x01\x33\x0d\x98\x7a\xc0\x45\x74\x22\x06\xbd\x6a\xde\x2c\x9d\xab\x2e\x0d\x7b\x22\xed\x46\xb7\xab\x54\xf9\x3b\xf6\xd8\x56\xa4\x35\x48\xa5\xd3\xfa\xd3\xc0\xf0\xf7\x00\x00\x00\xff\xff\x58\x83\x04\xe4\xfb\x00\x00\x00")

func resources_ddl_mysql_select_stops_by_date_sql_bytes() ([]byte, error) {
	return bindata_read(
		_resources_ddl_mysql_select_stops_by_date_sql,
		"resources/ddl/mysql/select-stops-by-date.sql",
	)
}

func resources_ddl_mysql_select_stops_by_date_sql() (*asset, error) {
	bytes, err := resources_ddl_mysql_select_stops_by_date_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "resources/ddl/mysql/select-stops-by-date.sql", size: 251, mode: os.FileMode(420), modTime: time.Unix(1428881302, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _resources_ddl_postgres_select_nearest_stations_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x74\x8d\x51\x0a\xc2\x30\x10\x44\xaf\xb2\x3f\xfe\x89\x27\x50\xaf\x22\x4b\x3b\x8d\xc1\x64\x02\x49\x44\xeb\xe9\xad\x15\xb7\x34\xd2\xbf\x79\xfb\x98\x9d\x82\x80\xae\x8a\x3a\xb0\x1b\x2f\x37\x8c\xfb\x5f\xf6\xbd\x45\x6a\x84\xc1\x3d\x07\xcb\xd5\x47\xbc\x12\x17\x19\x94\xce\x20\x7a\x4e\x87\xba\xb0\x3e\xd7\xfc\xf1\x89\x6b\x9f\x28\x43\x4e\x51\x5c\x1d\xca\x61\x16\x1e\x45\x1e\x57\x64\x34\x7f\xe5\x78\x92\x5d\x11\x65\xdf\x0c\xc8\xf9\x4f\x7c\x97\x36\x1a\x93\x98\x1b\xef\x00\x00\x00\xff\xff\xea\x44\x32\x64\x0c\x01\x00\x00")

func resources_ddl_postgres_select_nearest_stations_sql_bytes() ([]byte, error) {
	return bindata_read(
		_resources_ddl_postgres_select_nearest_stations_sql,
		"resources/ddl/postgres/select-nearest-stations.sql",
	)
}

func resources_ddl_postgres_select_nearest_stations_sql() (*asset, error) {
	bytes, err := resources_ddl_postgres_select_nearest_stations_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "resources/ddl/postgres/select-nearest-stations.sql", size: 268, mode: os.FileMode(420), modTime: time.Unix(1428879790, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _resources_ddl_postgres_select_stop_times_by_calendar_dates_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x74\x91\xc1\x4e\xc4\x30\x0c\x44\xef\xfd\x8a\x5c\xaa\x82\x84\xf6\x80\x38\x21\xf5\xbe\x1c\xf8\x86\xc8\x4a\x5c\x1a\x94\xb5\x8b\xed\x2e\xf0\xf7\x34\x14\xd1\x6e\x56\xdb\x5b\x9f\x27\xe3\x99\x44\x31\x63\x30\xd7\xb8\xe5\x53\x1b\x0e\x6a\x3c\xf9\x14\x1f\x2e\x01\xc1\x09\x2b\x14\x51\x43\x85\x32\x58\x4d\x98\x36\x92\x39\x80\x25\x26\x6f\xdf\xd3\x9f\x9b\xb1\x0f\x23\xc8\x5d\x19\x83\x48\x3a\x43\xf6\x96\x96\x5d\xae\x3b\x1e\x1f\x9f\x9e\x5f\x5f\xba\xfb\x6b\x65\xc4\x09\xc4\x66\xc1\x1b\xda\xff\xf5\x8a\x1f\x33\x52\xd8\x45\x8f\x49\x96\xba\x25\xc4\xbe\xa3\xf0\x6c\xe8\x75\x64\xb1\xaa\xea\x3a\xd9\x02\x6f\x2c\x70\x66\xb9\x12\xe2\x97\xd5\x13\x93\x54\x6e\xb4\x19\x84\x4f\xbf\xf0\xcd\x06\xf5\xad\xae\x11\x4b\x03\xf5\xc3\x9c\x73\x11\xbb\x44\x84\xe2\xde\x39\xd1\x85\x34\x40\x46\x8a\x20\x3e\x82\xa1\xba\x10\x1d\xd3\x5a\x13\xe5\x9c\x02\x2e\xfe\x7d\x88\xbb\xbf\xe6\x73\x44\xc1\xfa\x55\xfb\x36\x3a\xa0\xb8\x9c\x3f\x14\x23\xd7\xbb\xae\xd5\xee\x27\x00\x00\xff\xff\xea\x46\x4a\x6c\x03\x02\x00\x00")

func resources_ddl_postgres_select_stop_times_by_calendar_dates_sql_bytes() ([]byte, error) {
	return bindata_read(
		_resources_ddl_postgres_select_stop_times_by_calendar_dates_sql,
		"resources/ddl/postgres/select-stop-times-by-calendar-dates.sql",
	)
}

func resources_ddl_postgres_select_stop_times_by_calendar_dates_sql() (*asset, error) {
	bytes, err := resources_ddl_postgres_select_stop_times_by_calendar_dates_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "resources/ddl/postgres/select-stop-times-by-calendar-dates.sql", size: 515, mode: os.FileMode(420), modTime: time.Unix(1428923142, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _resources_ddl_postgres_select_stop_times_by_calendars_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x74\x91\xc1\x4e\xc4\x30\x0c\x44\xef\xfb\x15\xbe\x54\x05\x69\xc5\x01\x71\x42\x94\xf3\x72\xe0\x1b\x22\x2b\x71\xd9\xa0\x6c\xbc\xd8\xee\x02\x7f\x4f\x42\x05\xa4\xad\xe8\xad\x6f\x46\xf6\x8c\xa3\x94\xc8\xdb\x0e\xca\xa7\x36\xde\xa8\xf1\xd9\xc5\xb0\x87\x25\xc9\x78\xa2\x35\x0b\xa4\x7e\xcd\x12\xda\x06\x71\x6e\x50\x62\x8f\x16\x39\x3b\xfb\x3c\xff\x4c\x34\x76\xfe\x88\x72\x55\x75\x14\x89\x17\x4c\xce\x62\x5d\xd8\x1f\x0e\xb7\x77\xf7\xcf\x4f\xfd\xf5\x7e\xe3\x0c\x74\x46\xb1\x49\xe8\x1f\xef\x6f\x00\xa5\xb7\x89\xb2\x6f\xf3\x87\x28\xa5\x74\x8d\xb1\xa8\x2a\x3c\x19\x39\x3d\xb2\xd8\xba\xf1\x2c\x35\xa1\xff\xa0\xe7\xc4\xb2\xb5\xd2\x87\x6d\x24\x93\x58\xaf\x0b\xbb\x51\xf8\xf4\x4d\x5f\x6c\x54\xd7\xe9\x9c\xb4\x16\x51\x37\x4e\x29\x55\x37\xc4\x9c\x49\xe0\x95\x63\x86\x85\xd7\x63\xa2\x1c\x50\x14\x3c\x70\x9e\x8b\x92\x5c\xa2\xa7\x32\x7b\xf0\xcd\x0f\xec\xde\x8f\x24\xb4\x7e\xdd\xa1\x0b\x80\x39\xcc\x53\x8b\xdf\xca\x21\x5d\x40\x23\x78\x18\xa0\xef\xb4\x6f\xd5\xb2\x6a\xd6\x1e\x1b\xad\xd3\xc1\x64\xa2\xaf\x00\x00\x00\xff\xff\x3c\xd2\x85\xcc\x3d\x02\x00\x00")

func resources_ddl_postgres_select_stop_times_by_calendars_sql_bytes() ([]byte, error) {
	return bindata_read(
		_resources_ddl_postgres_select_stop_times_by_calendars_sql,
		"resources/ddl/postgres/select-stop-times-by-calendars.sql",
	)
}

func resources_ddl_postgres_select_stop_times_by_calendars_sql() (*asset, error) {
	bytes, err := resources_ddl_postgres_select_stop_times_by_calendars_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "resources/ddl/postgres/select-stop-times-by-calendars.sql", size: 573, mode: os.FileMode(420), modTime: time.Unix(1428923142, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _resources_ddl_postgres_select_stops_by_date_sql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xac\x8d\xcf\xca\x02\x31\x0c\xc4\x5f\x25\x97\x65\x77\x3f\x3e\x04\x0f\x1e\x04\xdf\xa5\xc4\x36\x6a\xa1\x6d\x4a\x13\xd0\x7d\x7b\xb3\xfe\xa9\xe0\xd9\x9c\x7e\x93\xcc\x64\x84\x12\x79\x05\xd9\x88\x72\x75\x31\xfc\xbf\xb1\x60\xa6\x2e\x02\x89\xef\x22\xa1\x7e\x98\xcb\xca\x89\x3d\x6a\xe4\xe2\x74\xa9\x96\xda\xda\xec\x77\xf0\x07\xa2\x2e\x44\x51\x2c\x9e\x26\xe3\x33\x71\x3e\x35\xce\x4a\x37\x9d\xc6\xca\xb1\xe8\x34\x08\x0c\x32\x8f\x73\x7f\x69\xa6\x19\x50\xe0\x59\xfc\x4a\xc3\x1a\x33\xe3\xc3\x62\x37\xb8\x5e\xa8\xd1\x0f\x8a\x0e\xb6\x04\x6e\x81\x1a\x1c\x97\xaf\x4e\x14\x7f\x0f\x00\x00\xff\xff\xe9\xc3\x56\xee\x1f\x01\x00\x00")

func resources_ddl_postgres_select_stops_by_date_sql_bytes() ([]byte, error) {
	return bindata_read(
		_resources_ddl_postgres_select_stops_by_date_sql,
		"resources/ddl/postgres/select-stops-by-date.sql",
	)
}

func resources_ddl_postgres_select_stops_by_date_sql() (*asset, error) {
	bytes, err := resources_ddl_postgres_select_stops_by_date_sql_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "resources/ddl/postgres/select-stops-by-date.sql", size: 287, mode: os.FileMode(420), modTime: time.Unix(1428885456, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"resources/ddl/mysql/select-nearest-stations.sql": resources_ddl_mysql_select_nearest_stations_sql,
	"resources/ddl/mysql/select-stop-times-by-calendar-dates.sql": resources_ddl_mysql_select_stop_times_by_calendar_dates_sql,
	"resources/ddl/mysql/select-stop-times-by-calendars.sql": resources_ddl_mysql_select_stop_times_by_calendars_sql,
	"resources/ddl/mysql/select-stops-by-date.sql": resources_ddl_mysql_select_stops_by_date_sql,
	"resources/ddl/postgres/select-nearest-stations.sql": resources_ddl_postgres_select_nearest_stations_sql,
	"resources/ddl/postgres/select-stop-times-by-calendar-dates.sql": resources_ddl_postgres_select_stop_times_by_calendar_dates_sql,
	"resources/ddl/postgres/select-stop-times-by-calendars.sql": resources_ddl_postgres_select_stop_times_by_calendars_sql,
	"resources/ddl/postgres/select-stops-by-date.sql": resources_ddl_postgres_select_stops_by_date_sql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"resources": &_bintree_t{nil, map[string]*_bintree_t{
		"ddl": &_bintree_t{nil, map[string]*_bintree_t{
			"mysql": &_bintree_t{nil, map[string]*_bintree_t{
				"select-nearest-stations.sql": &_bintree_t{resources_ddl_mysql_select_nearest_stations_sql, map[string]*_bintree_t{
				}},
				"select-stop-times-by-calendar-dates.sql": &_bintree_t{resources_ddl_mysql_select_stop_times_by_calendar_dates_sql, map[string]*_bintree_t{
				}},
				"select-stop-times-by-calendars.sql": &_bintree_t{resources_ddl_mysql_select_stop_times_by_calendars_sql, map[string]*_bintree_t{
				}},
				"select-stops-by-date.sql": &_bintree_t{resources_ddl_mysql_select_stops_by_date_sql, map[string]*_bintree_t{
				}},
			}},
			"postgres": &_bintree_t{nil, map[string]*_bintree_t{
				"select-nearest-stations.sql": &_bintree_t{resources_ddl_postgres_select_nearest_stations_sql, map[string]*_bintree_t{
				}},
				"select-stop-times-by-calendar-dates.sql": &_bintree_t{resources_ddl_postgres_select_stop_times_by_calendar_dates_sql, map[string]*_bintree_t{
				}},
				"select-stop-times-by-calendars.sql": &_bintree_t{resources_ddl_postgres_select_stop_times_by_calendars_sql, map[string]*_bintree_t{
				}},
				"select-stops-by-date.sql": &_bintree_t{resources_ddl_postgres_select_stops_by_date_sql, map[string]*_bintree_t{
				}},
			}},
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

