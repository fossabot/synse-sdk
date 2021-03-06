package devices

import (
	"github.com/vapor-ware/synse-sdk/sdk"
)

// Volt1103 is the handler for the example voltage device with model "volt1103".
var Volt1103 = sdk.DeviceHandler{
	Name: "voltage",

	Read: func(device *sdk.Device) ([]*sdk.Reading, error) {
		reading, err := device.GetOutput("voltage").MakeReading(1)
		if err != nil {
			return nil, err
		}
		return []*sdk.Reading{reading}, nil
	},
}
