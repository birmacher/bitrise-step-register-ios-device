package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-deploy-to-itunesconnect-deliver/appleauth"
	"github.com/bitrise-steplib/steps-deploy-to-itunesconnect-deliver/devportalservice"
	"github.com/bitrise-steplib/steps-ios-auto-provision-appstoreconnect/appstoreconnect"
)

const noDeveloperAccountConnectedWarning = `Connected Apple Developer Portal Account not found.
Most likely because there is no Apple Developer Portal Account connected to the build.
Read more: https://devcenter.bitrise.io/getting-started/configuring-bitrise-steps-that-require-apple-developer-account-data/`

func handleSessionDataError(err error) {
	if err == nil {
		return
	}

	if networkErr, ok := err.(devportalservice.NetworkError); ok && networkErr.Status == http.StatusUnauthorized {
		log.Warnf("Building a Pull Request for a Public App. Secret environments are not available in this build to protect them.\nThis will prevent us to fetch Bitrise Apple Developer Portal connection.")
		return
	}

	log.Errorf("Failed to activate Bitrise Apple Developer Portal connection:\n%v", err)
	log.Warnf("Failing back to step inputs.\nRead more about this issue: https://devcenter.bitrise.io/getting-started/configuring-bitrise-steps-that-require-apple-developer-account-data/")
}

func setupStepConfigs() (Config, error) {
	var stepConf Config
	if err := stepconf.Parse(&stepConf); err != nil {
		return Config{}, fmt.Errorf("Failed to read step configs:\n%s", err)
	}
	stepconf.Print(stepConf)
	log.Printf("")

	return stepConf, nil
}

func setupAppStoreConnectAPIClient(config Config) (*appstoreconnect.Client, error) {
	// Creating AppstoreConnectAPI client
	log.Infof("Setup App Store Connect API connection")

	// Setup API connections
	authInputs := appleauth.Inputs{
		APIIssuer:  config.APIIssuer,
		APIKeyPath: string(config.APIKeyPath),
	}
	if err := authInputs.Validate(); err != nil {
		return nil, fmt.Errorf("Failed to validate App Store Connect API inputs:\n%v", err)
	}

	// Authentication sources
	// First try to authenticate with connected account fetched from Bitrise
	// if it fails try with configs fetched from step configs
	authSources := []appleauth.Source{
		&appleauth.ConnectionAPIKeySource{},
		&appleauth.InputAPIKeySource{},
	}

	// Setup connection with the connected account stored on bitrise.io
	var devportalConnectionProvider *devportalservice.BitriseClient
	var appleDeveloperPortalConnection *devportalservice.AppleDeveloperConnection
	if config.BuildURL != "" && config.BuildAPIToken != "" {
		devportalConnectionProvider = devportalservice.NewBitriseClient(http.DefaultClient, config.BuildURL, string(config.BuildAPIToken))

		if devportalConnectionProvider != nil {
			appleDeveloperPortalConnection, err := devportalConnectionProvider.GetAppleDeveloperConnection()
			if err != nil {
				handleSessionDataError(err)
			}

			if appleDeveloperPortalConnection != nil && (appleDeveloperPortalConnection.APIKeyConnection == nil) {
				log.Warnf("%s", noDeveloperAccountConnectedWarning)
			}
		}
	} else {
		log.Warnf("Failed to fetch connected Apple Developer Portal Account from bitrise.io.\nStep is not running on bitrise.io: BITRISE_BUILD_URL and BITRISE_BUILD_API_TOKEN envs are not set")
	}

	// Setup configs with newly acquired bitrise account, or fall back to step inputs
	authConfig, err := appleauth.Select(appleDeveloperPortalConnection, authSources, authInputs)
	if err != nil {
		return nil, fmt.Errorf("Failed to configure App Store Connect API authentication:\n%v", err)
	}

	// Setup connection
	client := appstoreconnect.NewClient(http.DefaultClient, authConfig.APIKey.KeyID, authConfig.APIKey.IssuerID, []byte(authConfig.APIKey.PrivateKey))
	client.EnableDebugLogs = false

	log.Donef("App Store Connect API connection setup successfully")
	return client, nil
}

func logErrorAndExitIfAny(err error) {
	if err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}
}

func main() {
	config, err := setupStepConfigs()
	logErrorAndExitIfAny(err)

	_, err = setupAppStoreConnectAPIClient(config)
	logErrorAndExitIfAny(err)

	os.Exit(0)
}
