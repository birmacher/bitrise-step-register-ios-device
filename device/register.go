package device

import (
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-ios-auto-provision-appstoreconnect/appstoreconnect"
)

func registerDevice(client *appstoreconnect.Client, device Device) error {
	if client == nil {
		return fmt.Errorf("Failed to estabilish connection: App Store Connect client not provided")
	}

	// Register device
	// The API seems to recognize existing devices even with different casing and '-' separator removed.
	// The Developer Portal UI does not let adding devices with unexpected casing or separators removed.
	// Did not fully validate the ability to add devices with changed casing (or '-' removed) via the API, so passing the UDID through unchanged.
	req := appstoreconnect.DeviceCreateRequest{
		Data: appstoreconnect.DeviceCreateRequestData{
			Attributes: appstoreconnect.DeviceCreateRequestDataAttributes{
				Name:     device.Name,
				UDID:     device.UDID,
				Platform: device.ASCPlatform(),
			},
			Type: "devices",
		},
	}

	_, err := client.Provisioning.RegisterNewDevice(req)
	if err != nil {
		rerr, ok := err.(*appstoreconnect.ErrorResponse)
		if ok && rerr.Response != nil {
			errorStr := fmt.Sprintf("Failed to register device %s (%s)", device.Name, device.UDID)
			for _, error := range rerr.Errors {
				errorStr += "\n" + error.Title + ": " + error.Detail
			}
			return fmt.Errorf("%s", errorStr)
		}
	}

	return nil
}

func RegisterDevices(client *appstoreconnect.Client, devices []Device) error {
	if client == nil {
		return fmt.Errorf("Failed to estabilish connection: App Store Connect client not provided")
	}

	for _, device := range devices {
		log.Printf("")
		log.Infof("Registering device %s (%s)", device.Name, device.UDID)

		ascDevices, err := ascListDevice(client, device)
		if err != nil {
			return err
		}

		if len(ascDevices) > 0 {
			log.Warnf("Device is already registered on App Store Connect, skipping")
			continue
		}

		if err := registerDevice(client, device); err != nil {
			return err
		}
		log.Donef("Device %s (%s) successfully registered", device.Name, device.UDID)
	}

	return nil
}
