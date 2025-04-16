package kinvest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
)

func mustCreateJsonReader(data map[string]any) *bytes.Reader {
	buff := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buff).Encode(data); err != nil {
		panic(err)
	}
	return bytes.NewReader(buff.Bytes())
}

func mustUnmarshalJsonBody(body io.Reader) map[string]any {
	var ret map[string]any
	if err := json.NewDecoder(body).Decode(&ret); err != nil {
		panic(err)
	}
	return ret
}

func unmarshalJsonBody(body io.Reader, data any) error {
	if err := json.NewDecoder(body).Decode(data); err != nil {
		return fmt.Errorf("failed to unmarshal json body: %w", err)
	}
	return nil
}

func parseAccount(account string) (*string, *int, error) {
	var first string
	var second string
	fmt.Scanf(account, "%s-%s", &first, &second)
	if len(first) != 8 || len(second) != 2 {
		return nil, nil, fmt.Errorf("invalid account format: %s", account)
	}

	secondInt, err := strconv.Atoi(second)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid account format: %s", account)
	}

	return &first, &secondInt, nil
}

// NetIF represents a network interface
type NetIF struct {
	Name    string
	MacAddr string
}

// GetNetIFs returns a list of network interfaces
// that are not loopback and have a valid MAC address
// and do not start with 'v'
func GetNetIFs() ([]*NetIF, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	netifs := make([]*NetIF, 0, len(ifas))
	for _, ifa := range ifas {
		name := ifa.Name
		mac := ifa.HardwareAddr.String()
		if mac == "" || len(mac) != 17 {
			continue
		}
		if len(name) == 0 || name == "lo" || name == "lo0" {
			continue
		}
		if len(name) > 0 && name[0] == 'v' {
			continue
		}
		netifs = append(netifs, &NetIF{
			Name:    name,
			MacAddr: mac,
		})
	}

	return netifs, nil
}

func strToInt(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func strToFloat(s string) float64 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}
