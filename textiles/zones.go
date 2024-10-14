package textiles

import (
	"cmp"
	"encoding/binary"
	"github.com/memmaker/go/geometry"
	"io"
	"os"
	"slices"
)

type MapZone map[geometry.Point]bool

func SaveZones(filename string, zones map[string]MapZone) error {
	file, openErr := os.Create(filename)
	if openErr != nil {
		return openErr
	}
	defer file.Close()
	handleErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	zoneCount := int32(len(zones))
	handleErr(binary.Write(file, binary.LittleEndian, zoneCount))
	for _, zoneName := range sortedNames(zones) {
		zone := zones[zoneName]
		writeString(file, zoneName)
		handleErr(binary.Write(file, binary.LittleEndian, int32(len(zone))))
		locations := sortedValues(zone)
		for _, pos := range locations {
			handleErr(binary.Write(file, binary.LittleEndian, int32(pos.X)))
			handleErr(binary.Write(file, binary.LittleEndian, int32(pos.Y)))
		}
	}

	return nil
}

func writeString(file io.Writer, str string) {
	handleErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	size := int32(len(str))
	handleErr(binary.Write(file, binary.LittleEndian, size))
	handleErr(binary.Write(file, binary.LittleEndian, []byte(str)))
}

func readString(file io.Reader) string {
	handleErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	var size int32
	handleErr(binary.Read(file, binary.LittleEndian, &size))
	buf := make([]byte, size)
	handleErr(binary.Read(file, binary.LittleEndian, buf))
	return string(buf)
}

func ReadZones(filename string) map[string]MapZone {
	file, openErr := os.Open(filename)
	if openErr != nil {
		return nil
	}
	defer file.Close()

	handleErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	var zones = make(map[string]MapZone)

	var zoneCount int32

	handleErr(binary.Read(file, binary.LittleEndian, &zoneCount))

	for i := 0; i < int(zoneCount); i++ {
		var zoneSize int32
		zoneName := readString(file)
		handleErr(binary.Read(file, binary.LittleEndian, &zoneSize))
		zone := make(MapZone)
		for j := 0; j < int(zoneSize); j++ {
			var x, y int32
			handleErr(binary.Read(file, binary.LittleEndian, &x))
			handleErr(binary.Read(file, binary.LittleEndian, &y))
			zone[geometry.Point{X: int(x), Y: int(y)}] = true
		}
		zones[zoneName] = zone
	}
	return zones
}

func sortedValues(zone MapZone) []geometry.Point {
	values := make([]geometry.Point, 0, len(zone))
	for pos := range zone {
		values = append(values, pos)
	}
	slices.SortStableFunc(values, func(i, j geometry.Point) int {
		if i.Y == j.Y {
			return cmp.Compare(i.X, j.X)
		}
		return cmp.Compare(i.Y, j.Y)
	})
	return values
}

func sortedNames(zones map[string]MapZone) []string {
	names := make([]string, 0, len(zones))
	for name := range zones {
		names = append(names, name)
	}
	slices.SortStableFunc(names, func(i, j string) int {
		return cmp.Compare(i, j)
	})
	return names
}
