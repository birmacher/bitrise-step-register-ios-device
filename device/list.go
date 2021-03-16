package device

import (
	"fmt"

	"github.com/bitrise-steplib/steps-ios-auto-provision-appstoreconnect/appstoreconnect"
	"github.com/bitrise-steplib/steps-ios-auto-provision-appstoreconnect/autoprovision"
)

func ascListDevices(client *appstoreconnect.Client) ([]appstoreconnect.Device, error) {
	var ascDevices []appstoreconnect.Device

	if client == nil {
		return ascDevices, fmt.Errorf("Failed to estabilish connection: App Store Connect client not provided")
	}

	var err error
	ascDevices, err = autoprovision.ListDevices(client, "", appstoreconnect.IOSDevice)
	if err != nil {
		return ascDevices, fmt.Errorf("Failed to list devices from App Store Connect\n%s", err)
	}

	return ascDevices, nil
}
