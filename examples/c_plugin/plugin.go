package main

import (
	"../../sdk"

	"log"
	"strconv"
	"time"
)


// ExamplePluginHandler is a plugin-specific handler required by the
// SDK. It defines the plugin's read and write functionality.
type ExamplePluginHandler struct {}

func (h *ExamplePluginHandler) Read(in sdk.Device) (sdk.ReadResource, error) {
	id, err := strconv.Atoi(in.Data()["id"]); if err != nil {
		log.Fatalf("Invalid device ID - should be an integer in configuration.")
	}
	value := CRead(id, in.Model())
	return sdk.ReadResource{
		Device: in.Uid(),
		Reading: []sdk.Reading{{time.Now().String(), in.Type(), value}},
	}, nil
}

func (h *ExamplePluginHandler) Write(in sdk.Device, data *sdk.WriteData) error {
	return nil
}


// ExampleDeviceHandler is a plugin-specific handler required by the
// SDK. It defines functions which are needed to parse/make sense of
// some of the plugin-specific configurations.
type ExampleDeviceHandler struct {}

func (h *ExampleDeviceHandler) GetProtocolIdentifiers(data map[string]string) string {
	return data["id"]
}

// The main function - this is where we will configure, create, and run
// the plugin.
func main() {
	config := sdk.PluginConfig{}
	config.FromFile("plugin.yml")

	p, err := sdk.NewPlugin(
		config,
		&ExamplePluginHandler{},
		&ExampleDeviceHandler{},
	)
	if err != nil {
		log.Fatal(err)
	}

	p.Run()
}