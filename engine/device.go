package engine

import (
	"net/http"

	"github.com/samuelngs/hyper/router"
	"github.com/ua-parser/uap-go/uaparser"
)

// Device struct
type Device struct {
	useragent router.Useragent
	hardware  router.Hardware
	os        router.OS
	req       *http.Request
	parsed    *uaparser.Client
	uaparser  *uaparser.Parser
}

// Useragent struct
type Useragent struct {
	family, major, minor, patch, raw string
}

// OS struct
type OS struct {
	family, major, minor, patch string
}

// Hardware struct
type Hardware struct {
	family, brand, model string
}

// Useragent to return ua information
func (v *Device) Useragent() router.Useragent {
	if v.useragent == nil {
		ru := v.req.UserAgent()
		if v.parsed == nil {
			v.parsed = v.uaparser.Parse(ru)
		}
		v.useragent = &Useragent{
			family: v.parsed.UserAgent.Family,
			major:  v.parsed.UserAgent.Major,
			minor:  v.parsed.UserAgent.Minor,
			patch:  v.parsed.UserAgent.Patch,
			raw:    ru,
		}
	}
	return v.useragent
}

// Hardware to return client hardware information
func (v *Device) Hardware() router.Hardware {
	if v.hardware == nil {
		ru := v.req.UserAgent()
		if v.parsed == nil {
			v.parsed = v.uaparser.Parse(ru)
		}
		v.hardware = &Hardware{
			family: v.parsed.Device.Family,
			brand:  v.parsed.Device.Brand,
			model:  v.parsed.Device.Model,
		}
	}
	return v.hardware
}

// OS returns operating system information
func (v *Device) OS() router.OS {
	if v.os == nil {
		ru := v.req.UserAgent()
		if v.parsed == nil {
			v.parsed = v.uaparser.Parse(ru)
		}
		v.os = &OS{
			family: v.parsed.Os.Family,
			major:  v.parsed.Os.Major,
			minor:  v.parsed.Os.Minor,
			patch:  v.parsed.Os.Patch,
		}
	}
	return v.os
}

// Family returns useragent family name
func (v *Useragent) Family() string {
	return v.family
}

// Major returns useragent major version
func (v *Useragent) Major() string {
	return v.major
}

// Minor returns useragent minor version
func (v *Useragent) Minor() string {
	return v.minor
}

// Patch returns useragent patch version
func (v *Useragent) Patch() string {
	return v.patch
}

// String returns useragent string
func (v *Useragent) String() string {
	return v.raw
}

// Family returns operating system family name
func (v *OS) Family() string {
	return v.family
}

// Major returns operating system major version
func (v *OS) Major() string {
	return v.major
}

// Minor returns operating system minor version
func (v *OS) Minor() string {
	return v.minor
}

// Patch returns operating system patch version
func (v *OS) Patch() string {
	return v.patch
}

// Family returns hardware family name
func (v *Hardware) Family() string {
	return v.family
}

// Brand returns hardware brand name
func (v *Hardware) Brand() string {
	return v.brand
}

// Model returns hardware model name
func (v *Hardware) Model() string {
	return v.model
}
