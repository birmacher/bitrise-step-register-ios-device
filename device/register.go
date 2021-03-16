package device

import (
	"fmt"
	"net/http"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-deploy-to-itunesconnect-deliver/devportalservice"
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
				Platform: appstoreconnect.IOS,
				UDID:     device.UDID,
			},
			Type: "devices",
		},
	}

	log.Infof("Registering device %s (%s) to App Store Connect", device.Name, device.UDID)
	_, err := client.Provisioning.RegisterNewDevice(req)
	if err != nil {
		rerr, ok := err.(*appstoreconnect.ErrorResponse)
		if ok && rerr.Response != nil && rerr.Response.StatusCode == http.StatusConflict {
			log.Warnf("Failed to register device: %s (%s), skipping", device.Name, device.UDID)
			for _, error := range rerr.Errors {
				log.Warnf("%s - %s", error.Title, error.Detail)
			}
		}

		return err
	}

	return nil
}

func registerDevices(client *appstoreconnect.Client, devices []Device) error {
	if client == nil {
		return fmt.Errorf("Failed to estabilish connection: App Store Connect client not provided")
	}

	ascDevices, err := ascListDevices(client)
	if err != nil {
		return err
	}

	for _, device := range devices {
		for _, ascDevice := range ascDevices {
			if devportalservice.IsEqualUDID(ascDevice.Attributes.UDID, device.UDID) {
				log.Infof("Device with %s (%s) UDID is already registered on App Store Connect, skipping", device.Name, device.UDID)
				continue
			}

			if err := registerDevice(client, device); err != nil {
				return err
			}
		}
	}

	return nil
}
