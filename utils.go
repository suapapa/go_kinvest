package kinvest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

func mustCreateJsonReader(data any) io.ReadCloser {
	buff := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buff).Encode(data); err != nil {
		panic(err)
	}
	return io.NopCloser(bytes.NewReader(buff.Bytes()))
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

func unmarshalYamlBody(body io.Reader, data any) error {
	if err := yaml.NewDecoder(body).Decode(data); err != nil {
		return fmt.Errorf("failed to unmarshal yaml body: %w", err)
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

func toStr[T any](v T) string {
	switch val := any(v).(type) {
	case bool:
		if val {
			return "Y"
		}
		return "N"
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	// case time.Time:
	// 	return val.Format("20060102")
	default:
		return ""
	}
}

func toInt[T any](v T) int {
	switch val := any(v).(type) {
	case string:
		if val == "" {
			return 0
		}
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0
		}
		return i
	case int:
		return val
	case float64:
		return int(val)
	default:
		return 0
	}
}

func toFloat[T any](v T) float64 {
	switch val := any(v).(type) {
	case string:
		if val == "" {
			return 0
		}
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0
		}
		return f
	case float64:
		return val
	case int:
		return float64(val)
	default:
		return 0
	}
}

func toTime[T any](v T) time.Time {
	switch val := any(v).(type) {
	case int64:
		if val == 0 {
			return time.Time{}
		}
		t := time.Unix(int64(val), 0)
		return t
	case string:
		if val == "" {
			return time.Time{}
		}
		t, err := time.Parse("20060102", val)
		if err != nil {
			return time.Time{}
		}
		return t
	case time.Time:
		return val
	default:
		return time.Time{}
	}
}

func hhmmssToTime(hms string) (time.Time, error) {
	now := time.Now()

	t, err := time.Parse("150405", hms)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}

	return time.Date(
		now.Year(), now.Month(), now.Day(),
		t.Hour(), t.Minute(), t.Second(), 0,
		time.Local,
	), nil
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func getExt(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ""
	}
	return ext[1:]
}

func ptr[T any](v T) *T {
	return &v
}
