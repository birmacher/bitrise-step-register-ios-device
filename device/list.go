package device

import (
	"fmt"

	"github.com/bitrise-steplib/steps-ios-auto-provision-appstoreconnect/appstoreconnect"
	"github.com/bitrise-steplib/steps-ios-auto-provision-appstoreconnect/autoprovision"
)

func ascListDevice(client *appstoreconnect.Client, device Device) ([]appstoreconnect.Device, error) {
	var ascDevices []appstoreconnect.Device

	if client == nil {
		return ascDevices, fmt.Errorf("Failed to estabilish connection: App Store Connect client not provided")
	}

	ascDevices, err := autoprovision.ListDevices(client, device.UDID, appstoreconnect.DevicePlatform(device.ASCPlatform()))
	if err != nil {
		rerr, ok := err.(*appstoreconnect.ErrorResponse)
		if ok && rerr.Response != nil {
			errorStr := fmt.Sprintf("Failed to register device %s (%s)", device.Name, device.UDID)
			for _, error := range rerr.Errors {
				errorStr += "\n" + error.Title + ": " + error.Detail
			}
			return ascDevices, fmt.Errorf("%s", errorStr)
		}
	}

	return ascDevices, nil
}
