package engine

import (
	"net"
	"net/http"
	"strings"

	"github.com/samuelngs/hyper/router"
	"github.com/ua-parser/uap-go/uaparser"
)

// Client struct
type Client struct {
	ip, host, protocol string
	device             router.Device
	req                *http.Request
	uaparser           *uaparser.Parser
}

// IP returns client ip address
func (v *Client) IP() string {
	if v.ip == "" {
		ip := strings.TrimSpace(v.req.Header.Get("X-Real-Ip"))
		if len(ip) > 0 {
			v.ip = ip
			return ip
		}
		ip = v.req.Header.Get("X-Forwarded-For")
		if index := strings.IndexByte(ip, ','); index >= 0 {
			ip = ip[0:index]
		}
		ip = strings.TrimSpace(ip)
		if len(ip) > 0 {
			v.ip = ip
			return ip
		}
		if ip, _, err := net.SplitHostPort(strings.TrimSpace(v.req.RemoteAddr)); err == nil {
			v.ip = ip
			return ip
		}
	}
	return v.ip
}

// Host returns request hostname
func (v *Client) Host() string {
	if v.host == "" {
		v.host = v.req.Host
	}
	return v.host
}

// Protocol returns protocol name
func (v *Client) Protocol() string {
	if v.protocol == "" {
		if v.req.TLS != nil {
			v.protocol = "https"
		} else {
			v.protocol = "http"
		}
	}
	return v.protocol
}

// Device returns device information
func (v *Client) Device() router.Device {
	if v.device == nil {
		v.device = &Device{
			req:      v.req,
			uaparser: v.uaparser,
		}
	}
	return v.device
}
