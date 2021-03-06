package sdk

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/errors"
)

var (
	// deviceConfigLocations is a map to track the locations for the unified
	// DeviceConfig. The key is the name of the Location.
	deviceConfigLocations map[string]*LocationConfig
)

// verifyConfigs verifies that all device configurations that the plugin has
// found are correct.
//
// Config verification is different than config validation. In general,
// config validation consists of checks to ensure that a field is supported,
// required fields are set, and that fields are set correctly. Validation
// can be thought of as happening at an individual config struct level.
//
// Config verification is done at a higher level. It uses all known configs
// to verify the "global config state". This means, for example, that if a
// DeviceConfig references a Location, that the Location is defined somewhere.
//
// Config verification is necessary because the SDK allows multiple config
// files to be specified, for certain configs. This means that we can not verify
// that all the information in a given config is correct until we have the
// whole picture of what exists.
func verifyConfigs(unifiedDeviceConfig *DeviceConfig) *errors.MultiError {
	log.Debug("[sdk] verifying unified device config")

	var multiErr = errors.NewMultiError("config verification")

	// Reset the tracking state for every run of config verification
	deviceConfigLocations = map[string]*LocationConfig{}

	// Verify that there are no conflicting device configurations. We want to
	// do this first. This has the side-effect of building the deviceConfigLocations
	// map, which we will use later to verify that all DeviceInstances reference a
	// known location.
	verifyDeviceConfigLocations(unifiedDeviceConfig, multiErr)

	// Verify that all device instances reference a valid location.
	verifyDeviceConfigInstances(unifiedDeviceConfig, multiErr)

	// Verify that device kinds/instances reference valid output types.
	verifyDeviceConfigOutputs(unifiedDeviceConfig, multiErr)

	log.Debugf("[sdk] config verification found %d error(s)", len(multiErr.Errors))
	return multiErr
}

// verifyDeviceConfigLocations verifies that there are no Locations specified
// in the unified DeviceConfig that have conflicting data.
func verifyDeviceConfigLocations(deviceConfig *DeviceConfig, multiErr *errors.MultiError) {
	log.Debug("[sdk] verifying device config location data")
	for _, location := range deviceConfig.Locations {
		loc, hasName := deviceConfigLocations[location.Name]

		// If we do not already have a Location cached with the given name, there
		// can be no conflict, so we just add it to the cache.
		if !hasName {
			deviceConfigLocations[location.Name] = location
			continue
		}

		// If we already have the location cached, make sure that this Location
		// is the same as the existing one. If not, we have a conflict.
		if !loc.Equals(location) {
			log.WithField("name", loc.Name).Error("[sdk] duplicate location name")
			multiErr.Add(
				errors.NewVerificationConflictError(
					"device",
					fmt.Sprintf("differing Location config with the same name: %s", loc.Name),
				),
			)
		}
	}
}

// verifyDeviceConfigInstances verifies that the device instances all reference valid
// locations. verifyDeviceConfigLocations needs to be called before this verification
// function, as the deviceConfigLocations map is populated there.
func verifyDeviceConfigInstances(deviceConfig *DeviceConfig, multiErr *errors.MultiError) {
	log.Debug("[sdk] verifying device config instance data")
	for _, device := range deviceConfig.Devices {
		for _, instance := range device.Instances {
			if instance.Location == "" {
				log.Error("[sdk] instance config does not specify location")
				multiErr.Add(
					errors.NewVerificationInvalidError(
						"device",
						"device instance needs a location specified, but is empty",
					),
				)
				continue
			}

			_, hasLocation := deviceConfigLocations[instance.Location]
			if !hasLocation {
				log.WithField("name", instance.Location).Error("[sdk] unknown location specified")
				multiErr.Add(
					errors.NewVerificationInvalidError(
						"device",
						fmt.Sprintf("unknown device instance location specified: %s", instance.Location),
					),
				)
				continue
			}

		}
	}
}

// verifyDeviceConfigOutputs verifies that the DeviceOutputs of DeviceKinds and
// DeviceInstances reference valid output types.
func verifyDeviceConfigOutputs(deviceConfig *DeviceConfig, multiErr *errors.MultiError) {
	log.Debug("[sdk] verifying device config output type data")
	for _, device := range deviceConfig.Devices {
		// Check the device-level outputs
		for _, output := range device.Outputs {
			_, hasOutput := ctx.outputTypes[output.Type]
			if !hasOutput {
				log.WithField("name", output.Type).Error("[sdk] unknown output type specified")
				multiErr.Add(
					errors.NewVerificationInvalidError(
						"device",
						fmt.Sprintf("unknown output type specified: %s", output.Type),
					),
				)
				continue
			}
		}

		// Check the instance-level outputs
		for _, instance := range device.Instances {
			for _, output := range instance.Outputs {
				_, hasOutput := ctx.outputTypes[output.Type]
				if !hasOutput {
					log.WithField("name", output.Type).Error("[sdk] unknown output type specified")
					multiErr.Add(
						errors.NewVerificationInvalidError(
							"device",
							fmt.Sprintf("unknown output type specified: %s", output.Type),
						),
					)
					continue
				}
			}
		}
	}
}
