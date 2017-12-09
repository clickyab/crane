package ip2location

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"net"
	"strconv"
)

// metaData
type metaData struct {
	databaseType      uint8
	databaseColumn    uint8
	databaseDay       uint8
	databaseMonth     uint8
	databaseYear      uint8
	ipv4databaseCount uint32
	ipv4databaseAddr  uint32
	ipv6databaseCount uint32
	ipv6databaseAddr  uint32
	ipv4indexBaseAddr uint32
	ipv6indexBaseAddr uint32
	ipv4columnSize    uint32
	ipv6columnSize    uint32
}

// Record is base record
type Record struct {
	CountryShort       string  `json:"countryShort"`
	CountryLong        string  `json:"countryLong"`
	Region             string  `json:"region"`
	City               string  `json:"city"`
	Isp                string  `json:"isp"`
	Latitude           float32 `json:"-"`
	Longitude          float32 `json:"-"`
	Domain             string  `json:"-"`
	ZipCode            string  `json:"-"`
	Timezone           string  `json:"-"`
	NetSpeed           string  `json:"-"`
	IddCode            string  `json:"-"`
	AreaCode           string  `json:"-"`
	WeatherStationCode string  `json:"-"`
	WeatherStationName string  `json:"-"`
	Mcc                string  `json:"-"`
	Mnc                string  `json:"-"`
	MobileBrand        string  `json:"-"`
	Elevation          float32 `json:"-"`
	UsageType          string  `json:"-"`
}

var f *fileMock
var meta metaData

var countryPosition = [25]uint8{0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
var regionPosition = [25]uint8{0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}
var cityPosition = [25]uint8{0, 0, 0, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}
var ispPosition = [25]uint8{0, 0, 3, 0, 5, 0, 7, 5, 7, 0, 8, 0, 9, 0, 9, 0, 9, 0, 9, 7, 9, 0, 9, 7, 9}
var latitudePosition = [25]uint8{0, 0, 0, 0, 0, 5, 5, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5}
var longitudePosition = [25]uint8{0, 0, 0, 0, 0, 6, 6, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6}
var domainPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 6, 8, 0, 9, 0, 10, 0, 10, 0, 10, 0, 10, 8, 10, 0, 10, 8, 10}
var zipCodePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 0, 7, 7, 7, 0, 7, 0, 7, 7, 7, 0, 7}
var timezonePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 8, 7, 8, 8, 8, 7, 8, 0, 8, 8, 8, 0, 8}
var netSpeedPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 11, 0, 11, 8, 11, 0, 11, 0, 11, 0, 11}
var iddCodePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 12, 0, 12, 0, 12, 9, 12, 0, 12}
var areaCodePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 13, 0, 13, 0, 13, 10, 13, 0, 13}
var weatherStationCodePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 14, 0, 14, 0, 14, 0, 14}
var weatherStationNamePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 15, 0, 15, 0, 15, 0, 15}
var mccPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 16, 0, 16, 9, 16}
var mncPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 17, 0, 17, 10, 17}
var mobileBrandPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 18, 0, 18, 11, 18}
var elevationPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 19, 0, 19}
var usageTypePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 12, 20}

const apiVersion string = "8.0.3"

var maxIPv6Range = big.NewInt(0)

const (
	countryShort       uint32 = 0x00001
	countryLong               = 0x00002
	region                    = 0x00004
	city                      = 0x00008
	isp                       = 0x00010
	latitude                  = 0x00020
	longitude                 = 0x00040
	domain                    = 0x00080
	zipCode                   = 0x00100
	timezone                  = 0x00200
	netSpeed                  = 0x00400
	iddCode                   = 0x00800
	areaCode                  = 0x01000
	weatherStationCode        = 0x02000
	weatherStationName        = 0x04000
	mcc                       = 0x08000
	mnc                       = 0x10000
	mobileBrand               = 0x20000
	elevation                 = 0x40000
	usageType                 = 0x80000

	all = countryShort | countryLong | region | city | isp | latitude | longitude | domain |
		zipCode | timezone | netSpeed | iddCode | areaCode | weatherStationCode | weatherStationName |
		mcc | mnc | mobileBrand | elevation | usageType
)
const (
	invalidAddress string = "Invalid IP address."
	missingFile           = "Invalid database file."
	notSupported          = "This parameter is unavailable for selected data file. Please upgrade the data file."
)

var metaOK bool

var countryPositionOffset uint32
var regionPositionOffset uint32
var cityPositionOffset uint32
var ispPositionOffset uint32
var domainPositionOffset uint32
var zipCodePositionOffset uint32
var latitudePositionOffset uint32
var longitudePositionOffset uint32
var timezonePositionOffset uint32
var netSpeedPositionOffset uint32
var iddCodePositionOffset uint32
var areaCodePositionOffset uint32
var weatherStationCodePositionOffset uint32
var weatherStationNamePositionOffset uint32
var mccPositionOffset uint32
var mncPositionOffset uint32
var mobileBrandPositionOffset uint32
var elevationPositionOffset uint32
var usageTypePositionOffset uint32

var countryEnabled bool
var regionEnabled bool
var cityEnabled bool
var ispEnabled bool
var domainEnabled bool
var zipCodeEnabled bool
var latitudeEnabled bool
var longitudeEnabled bool
var timezoneEnabled bool
var netSpeedEnabled bool
var iddCodeEnabled bool
var areaCodeEnabled bool
var weatherStationCodeEnabled bool
var weatherStationNameEnabled bool
var mccEnabled bool
var mncEnabled bool
var mobileBrandEnabled bool
var elevationEnabled bool
var usageTypeEnabled bool

// get IP type and calculate IP number; calculates index too if exists
func checkIP(ip string) (IPType uint32, IPNum *big.Int, IPIndex uint32) {
	IPType = 0
	IPNum = big.NewInt(0)
	IPNumTmp := big.NewInt(0)
	IPIndex = 0
	ipAddress := net.ParseIP(ip)

	if ipAddress != nil {
		v4 := ipAddress.To4()

		if v4 != nil {
			IPType = 4
			IPNum.SetBytes(v4)
		} else {
			v6 := ipAddress.To16()

			if v6 != nil {
				IPType = 6
				IPNum.SetBytes(v6)
			}
		}
	}
	if IPType == 4 {
		if meta.ipv4indexBaseAddr > 0 {
			IPNumTmp.Rsh(IPNum, 16)
			IPNumTmp.Lsh(IPNumTmp, 3)
			IPIndex = uint32(IPNumTmp.Add(IPNumTmp, big.NewInt(int64(meta.ipv4indexBaseAddr))).Uint64())
		}
	} else if IPType == 6 {
		if meta.ipv6indexBaseAddr > 0 {
			IPNumTmp.Rsh(IPNum, 112)
			IPNumTmp.Lsh(IPNumTmp, 3)
			IPIndex = uint32(IPNumTmp.Add(IPNumTmp, big.NewInt(int64(meta.ipv6indexBaseAddr))).Uint64())
		}
	}
	return
}

// read byte
func readUint8(pos int64) uint8 {
	var retval uint8
	data := make([]byte, 1)
	_, err := f.ReadAt(data, pos-1)
	if err != nil {
		fmt.Println("File read failed:", err)
	}
	retval = data[0]
	return retval
}

// read unsigned 32-bit integer
func readUint32(pos uint32) uint32 {
	pos2 := int64(pos)
	var retval uint32
	data := make([]byte, 4)
	_, err := f.ReadAt(data, pos2-1)
	if err != nil {
		fmt.Println("File read failed:", err)
	}
	buf := bytes.NewReader(data)
	err = binary.Read(buf, binary.LittleEndian, &retval)
	if err != nil {
		fmt.Println("Binary read failed:", err)
	}
	return retval
}

// read unsigned 128-bit integer
func readUint128(pos uint32) *big.Int {
	pos2 := int64(pos)
	retval := big.NewInt(0)
	data := make([]byte, 16)
	_, err := f.ReadAt(data, pos2-1)
	if err != nil {
		fmt.Println("File read failed:", err)
	}

	// little endian to big endian
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	retval.SetBytes(data)
	return retval
}

// read string
func readStr(pos uint32) string {
	pos2 := int64(pos)
	var retval string
	lenbyte := make([]byte, 1)
	_, err := f.ReadAt(lenbyte, pos2)
	if err != nil {
		fmt.Println("File read failed:", err)
	}
	strlen := lenbyte[0]
	data := make([]byte, strlen)
	_, err = f.ReadAt(data, pos2+1)
	if err != nil {
		fmt.Println("File read failed:", err)
	}
	retval = string(data[:strlen])
	return retval
}

// read float
func readFloat(pos uint32) float32 {
	pos2 := int64(pos)
	var retval float32
	data := make([]byte, 4)
	_, err := f.ReadAt(data, pos2-1)
	if err != nil {
		fmt.Println("File read failed:", err)
	}
	buf := bytes.NewReader(data)
	err = binary.Read(buf, binary.LittleEndian, &retval)
	if err != nil {
		fmt.Println("Binary read failed:", err)
	}
	return retval
}

func initRegion(dbt uint8) {
	if countryPosition[dbt] != 0 {
		countryPositionOffset = uint32(countryPosition[dbt]-1) << 2
		countryEnabled = true
	}
	if regionPosition[dbt] != 0 {
		regionPositionOffset = uint32(regionPosition[dbt]-1) << 2
		regionEnabled = true
	}
	if cityPosition[dbt] != 0 {
		cityPositionOffset = uint32(cityPosition[dbt]-1) << 2
		cityEnabled = true
	}
	if ispPosition[dbt] != 0 {
		ispPositionOffset = uint32(ispPosition[dbt]-1) << 2
		ispEnabled = true
	}
	if domainPosition[dbt] != 0 {
		domainPositionOffset = uint32(domainPosition[dbt]-1) << 2
		domainEnabled = true
	}
	if zipCodePosition[dbt] != 0 {
		zipCodePositionOffset = uint32(zipCodePosition[dbt]-1) << 2
		zipCodeEnabled = true
	}
	if latitudePosition[dbt] != 0 {
		latitudePositionOffset = uint32(latitudePosition[dbt]-1) << 2
		latitudeEnabled = true
	}
	if longitudePosition[dbt] != 0 {
		longitudePositionOffset = uint32(longitudePosition[dbt]-1) << 2
		longitudeEnabled = true
	}
	if timezonePosition[dbt] != 0 {
		timezonePositionOffset = uint32(timezonePosition[dbt]-1) << 2
		timezoneEnabled = true
	}
}

// initialize the component with the database path
func open() error {
	maxIPv6Range.SetString("340282366920938463463374607431768211455", 10)

	var err error
	f, err = newFileMock()
	if err != nil {
		return err
	}

	meta.databaseType = readUint8(1)
	meta.databaseColumn = readUint8(2)
	meta.databaseYear = readUint8(3)
	meta.databaseMonth = readUint8(4)
	meta.databaseDay = readUint8(5)
	meta.ipv4databaseCount = readUint32(6)
	meta.ipv4databaseAddr = readUint32(10)
	meta.ipv6databaseCount = readUint32(14)
	meta.ipv6databaseAddr = readUint32(18)
	meta.ipv4indexBaseAddr = readUint32(22)
	meta.ipv6indexBaseAddr = readUint32(26)
	meta.ipv4columnSize = uint32(meta.databaseColumn << 2)              // 4 bytes each column
	meta.ipv6columnSize = uint32(16 + ((meta.databaseColumn - 1) << 2)) // 4 bytes each column, except IPFrom column which is 16 bytes

	dbt := meta.databaseType

	// since both IPv4 and IPv6 use 4 bytes for the below columns, can just do it once here
	initRegion(dbt)
	if netSpeedPosition[dbt] != 0 {
		netSpeedPositionOffset = uint32(netSpeedPosition[dbt]-1) << 2
		netSpeedEnabled = true
	}
	if iddCodePosition[dbt] != 0 {
		iddCodePositionOffset = uint32(iddCodePosition[dbt]-1) << 2
		iddCodeEnabled = true
	}
	if areaCodePosition[dbt] != 0 {
		areaCodePositionOffset = uint32(areaCodePosition[dbt]-1) << 2
		areaCodeEnabled = true
	}
	if weatherStationCodePosition[dbt] != 0 {
		weatherStationCodePositionOffset = uint32(weatherStationCodePosition[dbt]-1) << 2
		weatherStationCodeEnabled = true
	}
	if weatherStationNamePosition[dbt] != 0 {
		weatherStationNamePositionOffset = uint32(weatherStationNamePosition[dbt]-1) << 2
		weatherStationNameEnabled = true
	}
	if mccPosition[dbt] != 0 {
		mccPositionOffset = uint32(mccPosition[dbt]-1) << 2
		mccEnabled = true
	}
	if mncPosition[dbt] != 0 {
		mncPositionOffset = uint32(mncPosition[dbt]-1) << 2
		mncEnabled = true
	}
	if mobileBrandPosition[dbt] != 0 {
		mobileBrandPositionOffset = uint32(mobileBrandPosition[dbt]-1) << 2
		mobileBrandEnabled = true
	}
	if elevationPosition[dbt] != 0 {
		elevationPositionOffset = uint32(elevationPosition[dbt]-1) << 2
		elevationEnabled = true
	}
	if usageTypePosition[dbt] != 0 {
		usageTypePositionOffset = uint32(usageTypePosition[dbt]-1) << 2
		usageTypeEnabled = true
	}

	metaOK = true
	return nil
}

// APIVersion get api version
func APIVersion() string {
	return apiVersion
}

// populate record with message
func loadMessage(mesg string) Record {
	var x Record

	x.CountryShort = mesg
	x.CountryLong = mesg
	x.Region = mesg
	x.City = mesg
	x.Isp = mesg
	x.Domain = mesg
	x.ZipCode = mesg
	x.Timezone = mesg
	x.NetSpeed = mesg
	x.IddCode = mesg
	x.AreaCode = mesg
	x.WeatherStationCode = mesg
	x.WeatherStationName = mesg
	x.Mcc = mesg
	x.Mnc = mesg
	x.MobileBrand = mesg
	x.UsageType = mesg

	return x
}

// GetAll return get all fields
func GetAll(ipAddress string) Record {
	return query(ipAddress, all)
}

// GetCountryShort return get country code
func GetCountryShort(ipAddress string) Record {
	return query(ipAddress, countryShort)
}

// GetCountryLong return get country name
func GetCountryLong(ipAddress string) Record {
	return query(ipAddress, countryLong)
}

// GetRegion return region
func GetRegion(ipAddress string) Record {
	return query(ipAddress, region)
}

// GetCity return  city
func GetCity(ipAddress string) Record {
	return query(ipAddress, city)
}

// GetIsp return isp
func GetIsp(ipAddress string) Record {
	return query(ipAddress, isp)
}

// GetLatitude return latitude
func GetLatitude(ipAddress string) Record {
	return query(ipAddress, latitude)
}

// GetLongitude return longitude
func GetLongitude(ipAddress string) Record {
	return query(ipAddress, longitude)
}

// GetDomain return domain
func GetDomain(ipAddress string) Record {
	return query(ipAddress, domain)
}

// GetZipCode return zip code
func GetZipCode(ipAddress string) Record {
	return query(ipAddress, zipCode)
}

// GetTimezone return time zone
func GetTimezone(ipAddress string) Record {
	return query(ipAddress, timezone)
}

// GetNetSpeed return net speed
func GetNetSpeed(ipAddress string) Record {
	return query(ipAddress, netSpeed)
}

// GetIddCode return idd code
func GetIddCode(ipAddress string) Record {
	return query(ipAddress, iddCode)
}

// GetAreaCode return area code
func GetAreaCode(ipAddress string) Record {
	return query(ipAddress, areaCode)
}

// GetWeatherstationcode return weather station code
func GetWeatherstationcode(ipAddress string) Record {
	return query(ipAddress, weatherStationCode)
}

// GetWeatherStationName return weather station name
func GetWeatherStationName(ipAddress string) Record {
	return query(ipAddress, weatherStationName)
}

// GetMcc return mobile country code
func GetMcc(ipAddress string) Record {
	return query(ipAddress, mcc)
}

// GetMnc return mobile network code
func GetMnc(ipAddress string) Record {
	return query(ipAddress, mnc)
}

// GetMobilebrand return mobile carrier brand
func GetMobilebrand(ipAddress string) Record {
	return query(ipAddress, mobileBrand)
}

// GetElevation return elevation
func GetElevation(ipAddress string) Record {
	return query(ipAddress, elevation)
}

// GetUsageType return usage type
func GetUsageType(ipAddress string) Record {
	return query(ipAddress, usageType)
}

// main query
func query(ipAddress string, mode uint32) Record {
	x := loadMessage(notSupported) // default message

	// read metadata
	if !metaOK {
		x = loadMessage(missingFile)
		return x
	}

	// check IP type and return IP number & index (if exists)
	ipType, ipNo, ipIndex := checkIP(ipAddress)

	if ipType == 0 {
		x = loadMessage(invalidAddress)
		return x
	}

	var colSize uint32
	var baseAddr uint32
	var low uint32
	var high uint32
	var mid uint32
	var rowOffset uint32
	var rowOffset2 uint32
	var ipFrom *big.Int
	var ipTo *big.Int
	var maxIP *big.Int

	if ipType == 4 {
		baseAddr = meta.ipv4databaseAddr
		high = meta.ipv4databaseCount
		maxIP = big.NewInt(4294967295)
		colSize = meta.ipv4columnSize
	} else {
		baseAddr = meta.ipv6databaseAddr
		high = meta.ipv6databaseCount
		maxIP = maxIPv6Range
		colSize = meta.ipv6columnSize
	}

	// reading index
	if ipIndex > 0 {
		low = readUint32(ipIndex)
		high = readUint32(ipIndex + 4)
	}

	if ipNo.Cmp(maxIP) >= 0 {
		ipNo = ipNo.Sub(ipNo, big.NewInt(1))
	}

	for low <= high {
		mid = (low + high) >> 1
		rowOffset = baseAddr + (mid * colSize)
		rowOffset2 = rowOffset + colSize

		if ipType == 4 {
			ipFrom = big.NewInt(int64(readUint32(rowOffset)))
			ipTo = big.NewInt(int64(readUint32(rowOffset2)))
		} else {
			ipFrom = readUint128(rowOffset)
			ipTo = readUint128(rowOffset2)
		}

		if ipNo.Cmp(ipFrom) >= 0 && ipNo.Cmp(ipTo) < 0 {
			if ipType == 6 {
				rowOffset = rowOffset + 12 // coz below is assuming all columns are 4 bytes, so got 12 left to go to make 16 bytes total
			}
			x = general(mode, x, rowOffset)
			x = latLng(mode, x, rowOffset)
			x = network(mode, x, rowOffset)
			x = mobile(mode, x, rowOffset)
			x = regionRecord(mode, x, rowOffset)

			return x
		}
		if ipNo.Cmp(ipFrom) < 0 {
			high = mid - 1
		} else {
			low = mid + 1
		}

	}
	return x
}

func regionRecord(mode uint32, x Record, rowOffset uint32) Record {
	if mode&countryShort == 1 && countryEnabled {
		x.CountryShort = readStr(readUint32(rowOffset + countryPositionOffset))
	}

	if mode&countryLong != 0 && countryEnabled {
		x.CountryLong = readStr(readUint32(rowOffset+countryPositionOffset) + 3)
	}

	if mode&region != 0 && regionEnabled {
		x.Region = readStr(readUint32(rowOffset + regionPositionOffset))
	}

	if mode&city != 0 && cityEnabled {
		x.City = readStr(readUint32(rowOffset + cityPositionOffset))
	}
	return x
}

func latLng(mode uint32, x Record, rowOffset uint32) Record {
	if mode&latitude != 0 && latitudeEnabled {
		x.Latitude = readFloat(rowOffset + latitudePositionOffset)
	}

	if mode&longitude != 0 && longitudeEnabled {
		x.Longitude = readFloat(rowOffset + longitudePositionOffset)
	}
	return x
}
func network(mode uint32, x Record, rowOffset uint32) Record {

	if mode&domain != 0 && domainEnabled {
		x.Domain = readStr(readUint32(rowOffset + domainPositionOffset))
	}

	if mode&netSpeed != 0 && netSpeedEnabled {
		x.NetSpeed = readStr(readUint32(rowOffset + netSpeedPositionOffset))
	}

	return x
}
func general(mode uint32, x Record, rowOffset uint32) Record {

	if mode&isp != 0 && ispEnabled {
		x.Isp = readStr(readUint32(rowOffset + ispPositionOffset))
	}

	if mode&zipCode != 0 && zipCodeEnabled {
		x.ZipCode = readStr(readUint32(rowOffset + zipCodePositionOffset))
	}

	if mode&timezone != 0 && timezoneEnabled {
		x.Timezone = readStr(readUint32(rowOffset + timezonePositionOffset))
	}

	if mode&iddCode != 0 && iddCodeEnabled {
		x.IddCode = readStr(readUint32(rowOffset + iddCodePositionOffset))
	}

	if mode&areaCode != 0 && areaCodeEnabled {
		x.AreaCode = readStr(readUint32(rowOffset + areaCodePositionOffset))
	}

	if mode&weatherStationCode != 0 && weatherStationCodeEnabled {
		x.WeatherStationCode = readStr(readUint32(rowOffset + weatherStationCodePositionOffset))
	}

	if mode&weatherStationName != 0 && weatherStationNameEnabled {
		x.WeatherStationName = readStr(readUint32(rowOffset + weatherStationNamePositionOffset))
	}

	if mode&elevation != 0 && elevationEnabled {
		f, _ := strconv.ParseFloat(readStr(readUint32(rowOffset+elevationPositionOffset)), 32)
		x.Elevation = float32(f)
	}

	if mode&usageType != 0 && usageTypeEnabled {
		x.UsageType = readStr(readUint32(rowOffset + usageTypePositionOffset))
	}
	return x
}

func mobile(mode uint32, x Record, rowOffset uint32) Record {

	if mode&mcc != 0 && mccEnabled {
		x.Mcc = readStr(readUint32(rowOffset + mccPositionOffset))
	}

	if mode&mnc != 0 && mncEnabled {
		x.Mnc = readStr(readUint32(rowOffset + mncPositionOffset))
	}

	if mode&mobileBrand != 0 && mobileBrandEnabled {
		x.MobileBrand = readStr(readUint32(rowOffset + mobileBrandPositionOffset))
	}
	return x
}

// PrintRecord for debugging purposes
func PrintRecord(x Record) {
	fmt.Printf("countryShort: %s\n", x.CountryShort)
	fmt.Printf("countryLong: %s\n", x.CountryLong)
	fmt.Printf("region: %s\n", x.Region)
	fmt.Printf("city: %s\n", x.City)
	fmt.Printf("isp: %s\n", x.Isp)
	fmt.Printf("latitude: %f\n", x.Latitude)
	fmt.Printf("longitude: %f\n", x.Longitude)
	fmt.Printf("domain: %s\n", x.Domain)
	fmt.Printf("zipCode: %s\n", x.ZipCode)
	fmt.Printf("timezone: %s\n", x.Timezone)
	fmt.Printf("netSpeed: %s\n", x.NetSpeed)
	fmt.Printf("iddCode: %s\n", x.IddCode)
	fmt.Printf("areaCode: %s\n", x.AreaCode)
	fmt.Printf("weatherStationCode: %s\n", x.WeatherStationCode)
	fmt.Printf("weatherStationName: %s\n", x.WeatherStationName)
	fmt.Printf("mcc: %s\n", x.Mcc)
	fmt.Printf("mnc: %s\n", x.Mnc)
	fmt.Printf("mobileBrand: %s\n", x.MobileBrand)
	fmt.Printf("elevation: %f\n", x.Elevation)
	fmt.Printf("usageType: %s\n", x.UsageType)
}
