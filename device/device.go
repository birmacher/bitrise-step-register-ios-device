package device

import (
	"strings"

	"github.com/bitrise-steplib/steps-ios-auto-provision-appstoreconnect/appstoreconnect"
)

type Device struct {
	Name     string
	UDID     string
	Platform string
}

func (d Device) ASCPlatform() appstoreconnect.BundleIDPlatform {
	switch strings.ToLower(d.Platform) {
	case "ios":
		return appstoreconnect.IOS
	case "macos":
		return appstoreconnect.MacOS
	case "universal":
		return appstoreconnect.Universal
	}

	return appstoreconnect.BundleIDPlatform("UNKNOWN")
}
