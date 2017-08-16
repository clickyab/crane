package mysql

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const nullString = "null"

// Int64Slice is simple slice to handle it for json field
type Int64Slice []int64

// Int64Array is used to handle real array in database
type Int64Array []int64

// NullTime is null-time for json in null
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// NullInt64 is null int64 for json in null
type NullInt64 struct {
	Int64 int64
	Valid bool // Valid is true if Int64 is not NULL
}

// NullString is the json friendly null string
type NullString struct {
	Valid  bool
	String string
}

type slaveStatus struct {
	SlaveIOState              sql.NullString `db:"slave_io_state"`
	MasterHost                sql.NullString `db:"master_host"`
	MasterUser                sql.NullString `db:"Master_User"`
	MasterPort                sql.NullString `db:"master_port"`
	ConnectRetry              sql.NullString `db:"Connect_Retry"`
	MasterLogFile             sql.NullString `db:"Master_Log_File"`
	ReadMasterLogPos          sql.NullString `db:"Read_Master_Log_Pos"`
	RelayLogFile              sql.NullString `db:"Relay_Log_File"`
	RelayLogPos               sql.NullString `db:"Relay_Log_Pos"`
	RelayMasterLogFile        sql.NullString `db:"Relay_Master_Log_File"`
	SlaveIORunning            string         `db:"Slave_IO_Running"`
	SlaveSQLRunning           string         `db:"Slave_SQL_Running"`
	ReplicateDoDB             sql.NullString `db:"Replicate_Do_DB"`
	ReplicateIgnoreDB         sql.NullString `db:"Replicate_Ignore_DB"`
	ReplicateDoTable          sql.NullString `db:"Replicate_Do_Table"`
	ReplicateIgnoreTable      sql.NullString `db:"Replicate_Ignore_Table"`
	ReplicateWildDoTable      sql.NullString `db:"Replicate_Wild_Do_Table"`
	ReplicateWildIgnoreTable  sql.NullString `db:"Replicate_Wild_Ignore_Table"`
	LastErrno                 int            `db:"Last_Errno"`
	LastError                 sql.NullString `db:"Last_Error"`
	SkipCounter               sql.NullString `db:"Skip_Counter"`
	ExecMasterLogPos          sql.NullString `db:"Exec_Master_Log_Pos"`
	RelayLogSpace             sql.NullString `db:"Relay_Log_Space"`
	UntilCondition            sql.NullString `db:"Until_Condition"`
	UntilLogFile              sql.NullString `db:"Until_Log_File"`
	UntilLogPos               sql.NullString `db:"Until_Log_Pos"`
	MasterSSLAllowed          sql.NullString `db:"Master_SSL_Allowed"`
	MasterSSLCAFile           sql.NullString `db:"Master_SSL_CA_File"`
	MasterSSLCAPath           sql.NullString `db:"Master_SSL_CA_Path"`
	MasterSSLCert             sql.NullString `db:"Master_SSL_Cert"`
	MasterSSLCipher           sql.NullString `db:"Master_SSL_Cipher"`
	MasterSSLKey              sql.NullString `db:"Master_SSL_Key"`
	SecondsBehindMaster       sql.NullInt64  `db:"Seconds_Behind_Master"`
	MasterSSLVerifyServerCert sql.NullString `db:"Master_SSL_Verify_Server_Cert"`
	LastIOErrno               sql.NullString `db:"Last_IO_Errno"`
	LastIOError               sql.NullString `db:"Last_IO_Error"`
	LastSQLErrno              sql.NullString `db:"Last_SQL_Errno"`
	LastSQLError              sql.NullString `db:"Last_SQL_Error"`
	ReplicateIgnoreServerIDs  sql.NullString `db:"Replicate_Ignore_Server_Ids"`
	MasterServerID            sql.NullString `db:"Master_Server_Id"`
	MasterSSLCrl              sql.NullString `db:"Master_SSL_Crl"`
	MasterSSLCrlpath          sql.NullString `db:"Master_SSL_Crlpath"`
	UsingGtid                 sql.NullString `db:"Using_Gtid"`
	GtidIOPos                 sql.NullString `db:"Gtid_IO_Pos"`
}

// GenericJSONField is used to handle generic json data in postgres
type GenericJSONField map[string]interface{}

// StringJSONArray is use to handle string to string map in postgres
type StringJSONArray []string

// Scan convert the json array ino string slice
func (is *Int64Slice) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}

	return json.Unmarshal(b, is)
}

// Value try to get the string slice representation in database
func (is Int64Array) Value() (driver.Value, error) {
	b, err := json.Marshal(is)
	if err != nil {
		return nil, err
	}
	// Its time to change [] to {}
	b = bytes.Replace(b, []byte("["), []byte("{"), 1)
	b = bytes.Replace(b, []byte("]"), []byte("}"), 1)

	return b, nil
}

// Scan convert the json array ino string slice
func (is *Int64Array) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	b = bytes.Replace(b, []byte("{"), []byte("["), 1)
	b = bytes.Replace(b, []byte("}"), []byte("]"), 1)

	return json.Unmarshal(b, is)
}

// Value try to get the string slice representation in database
func (is Int64Slice) Value() (driver.Value, error) {
	return json.Marshal(is)
}

// MarshalJSON try to marshaling to json
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return nt.Time.MarshalJSON()
	}

	return []byte(nullString), nil
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// UnmarshalJSON try to unmarshal dae from input
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	text := strings.ToLower(string(b))
	if text == nullString {
		nt.Valid = false
		nt.Time = time.Time{}
		return nil
	}

	err := json.Unmarshal(b, &nt.Time)
	if err != nil {
		return err
	}

	nt.Valid = true
	return nil
}

// MarshalJSON try to marshaling to json
func (nt NullInt64) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return []byte(fmt.Sprintf(`%d`, nt.Int64)), nil
	}

	return []byte(nullString), nil
}

// UnmarshalJSON try to unmarshal dae from input
func (nt *NullInt64) UnmarshalJSON(b []byte) error {
	text := strings.ToLower(string(b))
	if text == nullString {
		nt.Valid = false

		return nil
	}

	err := json.Unmarshal(b, &nt.Int64)
	if err != nil {
		return err
	}

	nt.Valid = true
	return nil
}

// Scan implements the Scanner interface.
func (nt *NullInt64) Scan(value interface{}) error {
	inn := &sql.NullInt64{}
	err := inn.Scan(value)
	if err != nil {
		return err
	}
	nt.Int64 = inn.Int64
	nt.Valid = inn.Valid
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullInt64) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Int64, nil
}

// Scan convert the json array ino string slice
func (gjf *GenericJSONField) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
		return nil
	default:
		return errors.New("unsupported type")
	}

	return json.Unmarshal(b, gjf)
}

// Value try to get the string slice representation in database
func (gjf GenericJSONField) Value() (driver.Value, error) {
	return json.Marshal(gjf)
}

// Scan convert the json array ino string slice
func (ss *StringJSONArray) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}

	return json.Unmarshal(b, ss)
}

// Value try to get the string slice representation in database
func (ss StringJSONArray) Value() (driver.Value, error) {
	return json.Marshal(ss)
}

// Scan implements the Scanner interface.
func (ns *NullString) Scan(value interface{}) error {
	tmp := &sql.NullString{}
	err := tmp.Scan(value)
	if err != nil {
		return err
	}
	ns.Valid = tmp.Valid
	ns.String = tmp.String
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// MarshalJSON try to marshaling to json
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}

	return []byte(nullString), nil
}

// UnmarshalJSON try to unmarshal dae from input
func (ns NullString) UnmarshalJSON(b []byte) error {
	text := strings.ToLower(string(b))
	if text == nullString {
		ns.Valid = false
		ns.String = ""
		return nil
	}

	err := json.Unmarshal(b, &ns.String)
	if err != nil {
		return err
	}

	ns.Valid = true
	return nil
}
