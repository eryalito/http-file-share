package listener

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"unicode"
)

// getLocalIPs returns a list of non-loopback IPv4 addresses of the machine.
func getLocalIPs() ([]string, error) {
	var ips []string
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	for _, i := range interfaces {
		if i.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if i.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := i.Addrs()
		if err != nil {
			log.Printf("Warning: failed to get addresses for interface %s: %v", i.Name, err)
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ipv4 := ip.To4()
			if ipv4 == nil {
				continue // not an IPv4 address
			}
			ips = append(ips, ipv4.String())
		}
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("no non-loopback IPv4 addresses found")
	}
	return ips, nil
}

func sanitizeFilename(name string) string {
	// Remove accents and replace with closest ASCII
	safe := make([]rune, 0, len(name))
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			safe = append(safe, r)
		} else if r == '.' || r == '-' || r == '_' {
			safe = append(safe, r)
		} else if unicode.IsSpace(r) {
			safe = append(safe, '_')
		}
		// skip other runes (quotes, commas, etc.)
	}
	// Remove any leading/trailing dots or dashes
	re := regexp.MustCompile(`^[\.\-_]+|[\.\-_]+$`)
	return re.ReplaceAllString(string(safe), "")
}
