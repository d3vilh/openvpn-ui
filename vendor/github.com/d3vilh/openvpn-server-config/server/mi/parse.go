package mi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//ParsePid gets pid from string
func ParsePid(input string) (int64, error) {
	a := strings.Split(trim(input), "\n")
	if len(a) != 1 {
		return int64(0), fmt.Errorf("Wrong number of lines, expected %d, got %d", 1, len(a))
	}
	if !isSuccess(a[0]) {
		return int64(0), fmt.Errorf("Bad response: %s", a[0])
	}
	return strconv.ParseInt(stripPrefix(a[0], "SUCCESS: pid="), 10, 64)
}

//ParseVersion gets version information from string
func ParseVersion(input string) (*Version, error) {
	v := Version{}
	a := strings.Split(trim(input), "\n")
	if len(a) != 3 {
		return nil, fmt.Errorf("Wrong number of lines, expected %d, got %d", 3, len(a))
	}
	v.OpenVPN = stripPrefix(a[0], "OpenVPN Version: ")
	v.Management = stripPrefix(a[1], "Management Version: ")

	return &v, nil
}

//ParseStats gets stats from string
func ParseStats(input string) (*LoadStats, error) {
	ls := LoadStats{}
	a := strings.Split(trim(input), "\n")

	if len(a) != 1 {
		return nil, fmt.Errorf("Wrong number of lines, expected %d, got %d", 1, len(a))
	}
	line := a[0]
	if !isSuccess(line) {
		return nil, fmt.Errorf("Bad response: %s", line)
	}

	dString := stripPrefix(line, "SUCCESS: ")
	dElements := strings.Split(dString, ",")
	var err error
	ls.NClients, err = getLStatsValue(dElements[0])
	if err != nil {
		return nil, err
	}
	ls.BytesIn, err = getLStatsValue(dElements[1])
	if err != nil {
		return nil, err
	}
	ls.BytesOut, err = getLStatsValue(dElements[2])
	if err != nil {
		return nil, err
	}
	return &ls, nil
}

//ParseStatus gets status information from string
func ParseStatus(input string) (*Status, error) {
	s := Status{}
	s.ClientList = make([]*OVClient, 0, 0)
	s.RoutingTable = make([]*RoutingPath, 0, 0)
	a := strings.Split(trim(input), "\n")
	for _, line := range a {
		fields := strings.Split(trim(line), ",")
		c := fields[0]
		switch {
		case c == "TITLE":
			s.Title = fields[1]
		case c == "TIME":
			s.Time = fields[1]
			s.TimeT = fields[2]
		case c == "ROUTING_TABLE":
			item := &RoutingPath{
				VirtualAddress: fields[1],
				CommonName:     fields[2],
				RealAddress:    fields[3],
				LastRef:        fields[4],
				LastRefT:       fields[5],
			}
			s.RoutingTable = append(s.RoutingTable, item)
		case c == "CLIENT_LIST":
			bytesR, _ := strconv.ParseUint(fields[5], 10, 64)
			bytesS, _ := strconv.ParseUint(fields[6], 10, 64)
			item := &OVClient{
				CommonName:      fields[1],
				RealAddress:     fields[2],
				VirtualAddress:  fields[3],
				VirtualIPv6:     fields[4],
				BytesReceived:   bytesR,
				BytesSent:       bytesS,
				ConnectedSince:  fields[7],
				ConnectedSinceT: fields[8],
				Username:        fields[9],
				ClientID:        fields[10],
				PeerID:          fields[11],
				DataCipher:      fields[12],
			}
			s.ClientList = append(s.ClientList, item)
		}
	}
	return &s, nil
}

//ParseSignal checks for error in response string
func ParseSignal(input string) error {
	a := strings.Split(trim(input), "\n")
	if len(a) != 1 {
		return fmt.Errorf("Wrong number of lines, expected %d, got %d", 1, len(a))
	}
	if !isSuccess(a[0]) {
		return fmt.Errorf("Bad response: %s", a[0])
	}
	return nil
}

//ParseKillSession gets kill command result from string
func ParseKillSession(input string) (string, error) {
	a := strings.Split(trim(input), "\n")

	if len(a) != 1 {
		return "", fmt.Errorf("Wrong number of lines, expected %d, got %d", 1, len(a))
	}
	line := a[0]
	if !isSuccess(line) {
		return "", errors.New(line)
	}

	return stripPrefix(line, "SUCCESS: "), nil
}

func getLStatsValue(s string) (int64, error) {
	a := strings.Split(s, "=")
	if len(a) != 2 {
		return int64(-1), errors.New("Parsing error")
	}
	return strconv.ParseInt(a[1], 10, 64)
}

func trim(s string) string {
	return strings.Trim(strings.Trim(s, "\r\n"), "\n")
}

func stripPrefix(s, prefix string) string {
	return trim(strings.Replace(s, prefix, "", 1))
}

func isSuccess(s string) bool {
	if strings.HasPrefix(s, "SUCCESS: ") {
		return true
	}
	return false
}
