package sdk

import (
	"testing"

	"github.com/vapor-ware/synse-server-grpc/go"
)

func TestValidateReadRequest(t *testing.T) {
	// everything is there
	request := &synse.ReadRequest{Device: "device", Board: "board", Rack: "rack"}
	err := validateReadRequest(request)
	if err != nil {
		t.Error("Got error when validating read request, but none expected.")
	}

	// missing device
	request = &synse.ReadRequest{Board: "board", Rack: "rack"}
	err = validateReadRequest(request)
	if err == nil {
		t.Error("Got no error when validating read request, but was expecting one.")
	}

	// missing board
	request = &synse.ReadRequest{Device: "device", Rack: "rack"}
	err = validateReadRequest(request)
	if err == nil {
		t.Error("Got no error when validating read request, but was expecting one.")
	}

	// missing rack
	request = &synse.ReadRequest{Device: "device", Board: "board"}
	err = validateReadRequest(request)
	if err == nil {
		t.Error("Got no error when validating read request, but was expecting one.")
	}
}

func TestValidateWriteRequest(t *testing.T) {
	// everything is there
	request := &synse.WriteRequest{Device: "device", Board: "board", Rack: "rack"}
	err := validateWriteRequest(request)
	if err != nil {
		t.Error("Got error when validating write request, but none expected.")
	}

	// missing device
	request = &synse.WriteRequest{Board: "board", Rack: "rack"}
	err = validateWriteRequest(request)
	if err == nil {
		t.Error("Got no error when validating write request, but was expecting one.")
	}

	// missing board
	request = &synse.WriteRequest{Device: "device", Rack: "rack"}
	err = validateWriteRequest(request)
	if err == nil {
		t.Error("Got no error when validating write request, but was expecting one.")
	}

	// missing rack
	request = &synse.WriteRequest{Device: "device", Board: "board"}
	err = validateWriteRequest(request)
	if err == nil {
		t.Error("Got no error when validating write request, but was expecting one.")
	}
}
