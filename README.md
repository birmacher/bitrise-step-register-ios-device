# Register iOS Device [![Bitrise Build Status]()]() [![Bitrise Step Version](https://img.shields.io/badge/version-0.0.1-blue)](https://www.bitrise.io/integrations/steps/register-ios-device) [![GitHub License](https://img.shields.io/badge/license-MIT-lightgrey.svg)](https://raw.githubusercontent.com/bitrise-steplib/steps-go-list/master/LICENSE) [![Bitrise Community](https://img.shields.io/badge/community-Bitrise%20Discuss-lightgrey)](https://discuss.bitrise.io/)

You can use this step to register your iOS device to the Apple Developer portal

## Examples

### List packages in the working directory excluding vendor/*

```yml
---
format_version: '8'
default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
project_type: other
workflows:
  register_device:
    steps:
    - git-clone: {}
    - register-ios-device:
        inputs:
        - api_key_path: $BITRISE_API_KEY_PATH # Path to your p8 file
        - api_issuer: $BITRISE_API_ISSUER     # iTunes Connect Issuer Key
        - device_name: "QA iPhone 12 Pro Max"
        - device_udid: "00000000-000000A000000000"
        - device_platform: "ios"
```

## Configuration

### Inputs

| Parameter | Description | Required | Default |
| --- | --- | --- | --- |
| api_key_path | Path to local or remote file that holds the API Key for iTunes Connect API (p8 file) | 👍 | "" |
| api_issuer | iTunes Connect API Issuer Key | 👍 | "" |
| build_api_token | Bitrise.io Build API token | - | $BITRISE_BUILD_API_TOKEN |
| build_url | Build URL on bitrise.io | - | $BITRISE_BUILD_URL |
| device_name | The name of the device that you want to register | 👍 | "" |
| device_udid | The UDID of the device that you want to register | 👍 | "" |
| device_platform | The platform of the device that you want to register | 👍 | ios |

### Outputs

This step does not generate any outputs

## Contributing

We welcome [pull requests](https://github.com/birmacher/steps-register-ios-device/pulls) and [issues](https://github.com/birmacher/steps-register-ios-device/issues) against this repository. 

For pull requests, work on your changes in a forked repository

### Running locally

Copy the tmp.bitrise.secrets.yml file and set the correct environment variables in it
```sh
cp tmp.bitrise.secrets.yml .bitrise.secrets.yml
```

Use the bitrise cli to [run your tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/)

### Creating your own steps

Follow [this guide](https://devcenter.bitrise.io/contributors/create-your-own-step/) if you would like to create your own step

