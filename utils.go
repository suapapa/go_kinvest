package kinvest

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
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
