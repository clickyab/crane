package ip2location

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"net"
	"strconv"
)

type ip2locationmeta struct {
	databasetype      uint8
	databasecolumn    uint8
	databaseday       uint8
	databasemonth     uint8
	databaseyear      uint8
	ipv4databasecount uint32
	ipv4databaseaddr  uint32
	ipv6databasecount uint32
	ipv6databaseaddr  uint32
	ipv4indexbaseaddr uint32
	ipv6indexbaseaddr uint32
	ipv4columnsize    uint32
	ipv6columnsize    uint32
}

type IP2Locationrecord struct {
	Country_short      string  `json:"country_short"`
	Country_long       string  `json:"country_long"`
	Region             string  `json:"region"`
	City               string  `json:"city"`
	Isp                string  `json:"-"`
	Latitude           float32 `json:"-"`
	Longitude          float32 `json:"-"`
	Domain             string  `json:"-"`
	Zipcode            string  `json:"-"`
	Timezone           string  `json:"-"`
	Netspeed           string  `json:"-"`
	Iddcode            string  `json:"-"`
	Areacode           string  `json:"-"`
	Weatherstationcode string  `json:"-"`
	Weatherstationname string  `json:"-"`
	Mcc                string  `json:"-"`
	Mnc                string  `json:"-"`
	Mobilebrand        string  `json:"-"`
	Elevation          float32 `json:"-"`
	Usagetype          string  `json:"-"`
}

var f *fileMock
var meta ip2locationmeta

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

const api_version string = "8.0.3"

var maxIPv4Range = big.NewInt(4294967295)
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

	all = countryShort | countryLong | region | city | isp | latitude | longitude | domain | zipCode | timezone | netSpeed | iddCode | areaCode | weatherStationCode | weatherStationName | mcc | mnc | mobileBrand | elevation | usageType
)
const (
	invalid_address string = "Invalid IP address."
	missing_file           = "Invalid database file."
	not_supported          = "This parameter is unavailable for selected data file. Please upgrade the data file."
)

var metaOK bool

var country_position_offset uint32
var region_position_offset uint32
var city_position_offset uint32
var isp_position_offset uint32
var domain_position_offset uint32
var zipcode_position_offset uint32
var latitude_position_offset uint32
var longitude_position_offset uint32
var timezone_position_offset uint32
var netspeed_position_offset uint32
var iddcode_position_offset uint32
var areacode_position_offset uint32
var weatherstationcode_position_offset uint32
var weatherstationname_position_offset uint32
var mcc_position_offset uint32
var mnc_position_offset uint32
var mobilebrand_position_offset uint32
var elevation_position_offset uint32
var usagetype_position_offset uint32

var country_enabled bool
var region_enabled bool
var city_enabled bool
var isp_enabled bool
var domain_enabled bool
var zipcode_enabled bool
var latitude_enabled bool
var longitude_enabled bool
var timezone_enabled bool
var netspeed_enabled bool
var iddcode_enabled bool
var areacode_enabled bool
var weatherstationcode_enabled bool
var weatherstationname_enabled bool
var mcc_enabled bool
var mnc_enabled bool
var mobilebrand_enabled bool
var elevation_enabled bool
var usagetype_enabled bool

// get IP type and calculate IP number; calculates index too if exists
func checkIP(ip string) (IPType uint32, IPNum *big.Int, IPIndex uint32) {
	IPType = 0
	IPNum = big.NewInt(0)
	IPNumTmp := big.NewInt(0)
	IPIndex = 0
	IPAddress := net.ParseIP(ip)

	if IPAddress != nil {
		v4 := IPAddress.To4()

		if v4 != nil {
			IPType = 4
			IPNum.SetBytes(v4)
		} else {
			v6 := IPAddress.To16()

			if v6 != nil {
				IPType = 6
				IPNum.SetBytes(v6)
			}
		}
	}
	if IPType == 4 {
		if meta.ipv4indexbaseaddr > 0 {
			IPNumTmp.Rsh(IPNum, 16)
			IPNumTmp.Lsh(IPNumTmp, 3)
			IPIndex = uint32(IPNumTmp.Add(IPNumTmp, big.NewInt(int64(meta.ipv4indexbaseaddr))).Uint64())
		}
	} else if IPType == 6 {
		if meta.ipv6indexbaseaddr > 0 {
			IPNumTmp.Rsh(IPNum, 112)
			IPNumTmp.Lsh(IPNumTmp, 3)
			IPIndex = uint32(IPNumTmp.Add(IPNumTmp, big.NewInt(int64(meta.ipv6indexbaseaddr))).Uint64())
		}
	}
	return
}

// read byte
func readuint8(pos int64) uint8 {
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
func readuint32(pos uint32) uint32 {
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
func readuint128(pos uint32) *big.Int {
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
func readstr(pos uint32) string {
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
func readfloat(pos uint32) float32 {
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

// initialize the component with the database path
func Open() error {
	maxIPv6Range.SetString("340282366920938463463374607431768211455", 10)

	var err error
	f, err = newFileMock()
	if err != nil {
		return err
	}

	meta.databasetype = readuint8(1)
	meta.databasecolumn = readuint8(2)
	meta.databaseyear = readuint8(3)
	meta.databasemonth = readuint8(4)
	meta.databaseday = readuint8(5)
	meta.ipv4databasecount = readuint32(6)
	meta.ipv4databaseaddr = readuint32(10)
	meta.ipv6databasecount = readuint32(14)
	meta.ipv6databaseaddr = readuint32(18)
	meta.ipv4indexbaseaddr = readuint32(22)
	meta.ipv6indexbaseaddr = readuint32(26)
	meta.ipv4columnsize = uint32(meta.databasecolumn << 2)              // 4 bytes each column
	meta.ipv6columnsize = uint32(16 + ((meta.databasecolumn - 1) << 2)) // 4 bytes each column, except IPFrom column which is 16 bytes

	dbt := meta.databasetype

	// since both IPv4 and IPv6 use 4 bytes for the below columns, can just do it once here
	if countryPosition[dbt] != 0 {
		country_position_offset = uint32(countryPosition[dbt]-1) << 2
		country_enabled = true
	}
	if regionPosition[dbt] != 0 {
		region_position_offset = uint32(regionPosition[dbt]-1) << 2
		region_enabled = true
	}
	if cityPosition[dbt] != 0 {
		city_position_offset = uint32(cityPosition[dbt]-1) << 2
		city_enabled = true
	}
	if ispPosition[dbt] != 0 {
		isp_position_offset = uint32(ispPosition[dbt]-1) << 2
		isp_enabled = true
	}
	if domainPosition[dbt] != 0 {
		domain_position_offset = uint32(domainPosition[dbt]-1) << 2
		domain_enabled = true
	}
	if zipCodePosition[dbt] != 0 {
		zipcode_position_offset = uint32(zipCodePosition[dbt]-1) << 2
		zipcode_enabled = true
	}
	if latitudePosition[dbt] != 0 {
		latitude_position_offset = uint32(latitudePosition[dbt]-1) << 2
		latitude_enabled = true
	}
	if longitudePosition[dbt] != 0 {
		longitude_position_offset = uint32(longitudePosition[dbt]-1) << 2
		longitude_enabled = true
	}
	if timezonePosition[dbt] != 0 {
		timezone_position_offset = uint32(timezonePosition[dbt]-1) << 2
		timezone_enabled = true
	}
	if netSpeedPosition[dbt] != 0 {
		netspeed_position_offset = uint32(netSpeedPosition[dbt]-1) << 2
		netspeed_enabled = true
	}
	if iddCodePosition[dbt] != 0 {
		iddcode_position_offset = uint32(iddCodePosition[dbt]-1) << 2
		iddcode_enabled = true
	}
	if areaCodePosition[dbt] != 0 {
		areacode_position_offset = uint32(areaCodePosition[dbt]-1) << 2
		areacode_enabled = true
	}
	if weatherStationCodePosition[dbt] != 0 {
		weatherstationcode_position_offset = uint32(weatherStationCodePosition[dbt]-1) << 2
		weatherstationcode_enabled = true
	}
	if weatherStationNamePosition[dbt] != 0 {
		weatherstationname_position_offset = uint32(weatherStationNamePosition[dbt]-1) << 2
		weatherstationname_enabled = true
	}
	if mccPosition[dbt] != 0 {
		mcc_position_offset = uint32(mccPosition[dbt]-1) << 2
		mcc_enabled = true
	}
	if mncPosition[dbt] != 0 {
		mnc_position_offset = uint32(mncPosition[dbt]-1) << 2
		mnc_enabled = true
	}
	if mobileBrandPosition[dbt] != 0 {
		mobilebrand_position_offset = uint32(mobileBrandPosition[dbt]-1) << 2
		mobilebrand_enabled = true
	}
	if elevationPosition[dbt] != 0 {
		elevation_position_offset = uint32(elevationPosition[dbt]-1) << 2
		elevation_enabled = true
	}
	if usageTypePosition[dbt] != 0 {
		usagetype_position_offset = uint32(usageTypePosition[dbt]-1) << 2
		usagetype_enabled = true
	}

	metaOK = true
	return nil
}

// close database file handle
func Close() {
}

// get api version
func Api_version() string {
	return api_version
}

// populate record with message
func loadmessage(mesg string) IP2Locationrecord {
	var x IP2Locationrecord

	x.Country_short = mesg
	x.Country_long = mesg
	x.Region = mesg
	x.City = mesg
	x.Isp = mesg
	x.Domain = mesg
	x.Zipcode = mesg
	x.Timezone = mesg
	x.Netspeed = mesg
	x.Iddcode = mesg
	x.Areacode = mesg
	x.Weatherstationcode = mesg
	x.Weatherstationname = mesg
	x.Mcc = mesg
	x.Mnc = mesg
	x.Mobilebrand = mesg
	x.Usagetype = mesg

	return x
}

// get all fields
func Get_all(ipaddress string) IP2Locationrecord {
	return query(ipaddress, all)
}

// get country code
func Get_country_short(ipaddress string) IP2Locationrecord {
	return query(ipaddress, countryShort)
}

// get country name
func Get_country_long(ipaddress string) IP2Locationrecord {
	return query(ipaddress, countryLong)
}

// get region
func Get_region(ipaddress string) IP2Locationrecord {
	return query(ipaddress, region)
}

// get city
func Get_city(ipaddress string) IP2Locationrecord {
	return query(ipaddress, city)
}

// get isp
func Get_isp(ipaddress string) IP2Locationrecord {
	return query(ipaddress, isp)
}

// get latitude
func Get_latitude(ipaddress string) IP2Locationrecord {
	return query(ipaddress, latitude)
}

// get longitude
func Get_longitude(ipaddress string) IP2Locationrecord {
	return query(ipaddress, longitude)
}

// get domain
func Get_domain(ipaddress string) IP2Locationrecord {
	return query(ipaddress, domain)
}

// get zip code
func Get_zipcode(ipaddress string) IP2Locationrecord {
	return query(ipaddress, zipCode)
}

// get time zone
func Get_timezone(ipaddress string) IP2Locationrecord {
	return query(ipaddress, timezone)
}

// get net speed
func Get_netspeed(ipaddress string) IP2Locationrecord {
	return query(ipaddress, netSpeed)
}

// get idd code
func Get_iddcode(ipaddress string) IP2Locationrecord {
	return query(ipaddress, iddCode)
}

// get area code
func Get_areacode(ipaddress string) IP2Locationrecord {
	return query(ipaddress, areaCode)
}

// get weather station code
func Get_weatherstationcode(ipaddress string) IP2Locationrecord {
	return query(ipaddress, weatherStationCode)
}

// get weather station name
func Get_weatherstationname(ipaddress string) IP2Locationrecord {
	return query(ipaddress, weatherStationName)
}

// get mobile country code
func Get_mcc(ipaddress string) IP2Locationrecord {
	return query(ipaddress, mcc)
}

// get mobile network code
func Get_mnc(ipaddress string) IP2Locationrecord {
	return query(ipaddress, mnc)
}

// get mobile carrier brand
func Get_mobilebrand(ipaddress string) IP2Locationrecord {
	return query(ipaddress, mobileBrand)
}

// get elevation
func Get_elevation(ipaddress string) IP2Locationrecord {
	return query(ipaddress, elevation)
}

// get usage type
func Get_usagetype(ipaddress string) IP2Locationrecord {
	return query(ipaddress, usageType)
}

// main query
func query(ipaddress string, mode uint32) IP2Locationrecord {
	x := loadmessage(not_supported) // default message

	// read metadata
	if !metaOK {
		x = loadmessage(missing_file)
		return x
	}

	// check IP type and return IP number & index (if exists)
	iptype, ipno, ipindex := checkIP(ipaddress)

	if iptype == 0 {
		x = loadmessage(invalid_address)
		return x
	}

	var colsize uint32
	var baseaddr uint32
	var low uint32
	var high uint32
	var mid uint32
	var rowoffset uint32
	var rowoffset2 uint32
	ipfrom := big.NewInt(0)
	ipto := big.NewInt(0)
	maxip := big.NewInt(0)

	if iptype == 4 {
		baseaddr = meta.ipv4databaseaddr
		high = meta.ipv4databasecount
		maxip = maxIPv4Range
		colsize = meta.ipv4columnsize
	} else {
		baseaddr = meta.ipv6databaseaddr
		high = meta.ipv6databasecount
		maxip = maxIPv6Range
		colsize = meta.ipv6columnsize
	}

	// reading index
	if ipindex > 0 {
		low = readuint32(ipindex)
		high = readuint32(ipindex + 4)
	}

	if ipno.Cmp(maxip) >= 0 {
		ipno = ipno.Sub(ipno, big.NewInt(1))
	}

	for low <= high {
		mid = ((low + high) >> 1)
		rowoffset = baseaddr + (mid * colsize)
		rowoffset2 = rowoffset + colsize

		if iptype == 4 {
			ipfrom = big.NewInt(int64(readuint32(rowoffset)))
			ipto = big.NewInt(int64(readuint32(rowoffset2)))
		} else {
			ipfrom = readuint128(rowoffset)
			ipto = readuint128(rowoffset2)
		}

		if ipno.Cmp(ipfrom) >= 0 && ipno.Cmp(ipto) < 0 {
			if iptype == 6 {
				rowoffset = rowoffset + 12 // coz below is assuming all columns are 4 bytes, so got 12 left to go to make 16 bytes total
			}

			if mode&countryShort == 1 && country_enabled {
				x.Country_short = readstr(readuint32(rowoffset + country_position_offset))
			}

			if mode&countryLong != 0 && country_enabled {
				x.Country_long = readstr(readuint32(rowoffset+country_position_offset) + 3)
			}

			if mode&region != 0 && region_enabled {
				x.Region = readstr(readuint32(rowoffset + region_position_offset))
			}

			if mode&city != 0 && city_enabled {
				x.City = readstr(readuint32(rowoffset + city_position_offset))
			}

			if mode&isp != 0 && isp_enabled {
				x.Isp = readstr(readuint32(rowoffset + isp_position_offset))
			}

			if mode&latitude != 0 && latitude_enabled {
				x.Latitude = readfloat(rowoffset + latitude_position_offset)
			}

			if mode&longitude != 0 && longitude_enabled {
				x.Longitude = readfloat(rowoffset + longitude_position_offset)
			}

			if mode&domain != 0 && domain_enabled {
				x.Domain = readstr(readuint32(rowoffset + domain_position_offset))
			}

			if mode&zipCode != 0 && zipcode_enabled {
				x.Zipcode = readstr(readuint32(rowoffset + zipcode_position_offset))
			}

			if mode&timezone != 0 && timezone_enabled {
				x.Timezone = readstr(readuint32(rowoffset + timezone_position_offset))
			}

			if mode&netSpeed != 0 && netspeed_enabled {
				x.Netspeed = readstr(readuint32(rowoffset + netspeed_position_offset))
			}

			if mode&iddCode != 0 && iddcode_enabled {
				x.Iddcode = readstr(readuint32(rowoffset + iddcode_position_offset))
			}

			if mode&areaCode != 0 && areacode_enabled {
				x.Areacode = readstr(readuint32(rowoffset + areacode_position_offset))
			}

			if mode&weatherStationCode != 0 && weatherstationcode_enabled {
				x.Weatherstationcode = readstr(readuint32(rowoffset + weatherstationcode_position_offset))
			}

			if mode&weatherStationName != 0 && weatherstationname_enabled {
				x.Weatherstationname = readstr(readuint32(rowoffset + weatherstationname_position_offset))
			}

			if mode&mcc != 0 && mcc_enabled {
				x.Mcc = readstr(readuint32(rowoffset + mcc_position_offset))
			}

			if mode&mnc != 0 && mnc_enabled {
				x.Mnc = readstr(readuint32(rowoffset + mnc_position_offset))
			}

			if mode&mobileBrand != 0 && mobilebrand_enabled {
				x.Mobilebrand = readstr(readuint32(rowoffset + mobilebrand_position_offset))
			}

			if mode&elevation != 0 && elevation_enabled {
				f, _ := strconv.ParseFloat(readstr(readuint32(rowoffset+elevation_position_offset)), 32)
				x.Elevation = float32(f)
			}

			if mode&usageType != 0 && usagetype_enabled {
				x.Usagetype = readstr(readuint32(rowoffset + usagetype_position_offset))
			}

			return x
		} else {
			if ipno.Cmp(ipfrom) < 0 {
				high = mid - 1
			} else {
				low = mid + 1
			}
		}
	}
	return x
}

// for debugging purposes
func Printrecord(x IP2Locationrecord) {
	fmt.Printf("country_short: %s\n", x.Country_short)
	fmt.Printf("country_long: %s\n", x.Country_long)
	fmt.Printf("region: %s\n", x.Region)
	fmt.Printf("city: %s\n", x.City)
	fmt.Printf("isp: %s\n", x.Isp)
	fmt.Printf("latitude: %f\n", x.Latitude)
	fmt.Printf("longitude: %f\n", x.Longitude)
	fmt.Printf("domain: %s\n", x.Domain)
	fmt.Printf("zipCode: %s\n", x.Zipcode)
	fmt.Printf("timezone: %s\n", x.Timezone)
	fmt.Printf("netSpeed: %s\n", x.Netspeed)
	fmt.Printf("iddCode: %s\n", x.Iddcode)
	fmt.Printf("areaCode: %s\n", x.Areacode)
	fmt.Printf("weatherStationCode: %s\n", x.Weatherstationcode)
	fmt.Printf("weatherStationName: %s\n", x.Weatherstationname)
	fmt.Printf("mcc: %s\n", x.Mcc)
	fmt.Printf("mnc: %s\n", x.Mnc)
	fmt.Printf("mobileBrand: %s\n", x.Mobilebrand)
	fmt.Printf("elevation: %f\n", x.Elevation)
	fmt.Printf("usageType: %s\n", x.Usagetype)
}
