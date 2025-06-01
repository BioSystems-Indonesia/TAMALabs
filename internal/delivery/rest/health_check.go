package rest

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type HealthCheckHandler struct {
	cfg *config.Schema
}

func NewHealthCheckHandler(cfg *config.Schema) *HealthCheckHandler {
	return &HealthCheckHandler{cfg: cfg}
}

// LocalIP get the host machine local IP address
func LocalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if isPrivateIP(ip) {
				return ip, nil
			}
		}
	}

	return nil, errors.New("no IP")
}

func isPrivateIP(ip net.IP) bool {
	var privateIPBlocks []*net.IPNet
	for _, cidr := range []string{
		// don't check loopback ips
		//"127.0.0.0/8",    // IPv4 loopback
		//"::1/128",        // IPv6 loopback
		//"fe80::/10",      // IPv6 link-local
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

func (h HealthCheckHandler) Ping(c echo.Context) error {
	var serverIP string
	localIP, err := LocalIP()
	if err != nil {
		serverIP = err.Error()
	} else {
		serverIP = localIP.String()
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":   "OK",
		"version":  h.cfg.Version,
		"revision": h.cfg.Revision,
		"logLevel": h.cfg.LogLevel,
		"serverIP": serverIP,
		"port":     h.cfg.Port,
	})
}

func (h HealthCheckHandler) CheckAuth(c echo.Context) error {
	admin := entity.GetEchoContextUser(c)

	return c.JSON(http.StatusOK, map[string]string{
		"status":     "OK",
		"id":         fmt.Sprintf("%d", admin.ID),
		"fullname":   admin.Fullname,
		"email":      admin.Email,
		"is_active":  fmt.Sprintf("%t", admin.IsActive),
		"created_at": admin.CreatedAt,
		"updated_at": admin.UpdatedAt,
	})
}
