package kinvest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
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
	accountParts := strings.Split(account, "-")
	if len(accountParts) != 2 {
		return nil, nil, fmt.Errorf("invalid account format: %s", account)
	}
	first, second := accountParts[0], accountParts[1]
	if len(first) != 8 || len(second) != 2 {
		return nil, nil, fmt.Errorf("invalid account format: %s", account)
	}

	secondInt, err := strconv.Atoi(second)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid account format: %s", account)
	}

	return &first, &secondInt, nil
}

func getLocalIPAndMAC() (string, string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", "", fmt.Errorf("failed to get network interfaces: %w", err)
	}

	for _, iface := range interfaces {
		// Skip loopback interfaces and interfaces without a MAC address
		if iface.Flags&net.FlagLoopback != 0 || iface.HardwareAddr == nil {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", "", fmt.Errorf("failed to get addresses for interface %s: %w", iface.Name, err)
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				return ipNet.IP.String(), iface.HardwareAddr.String(), nil
			}
		}
	}

	return "", "", fmt.Errorf("no valid network interface found")
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

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
