package solana

import (
	"encoding/json"

	"github.com/pokt-network/poktroll/pkg/polylog"

	"github.com/buildwithgrove/path/qos/jsonrpc"
)

// responseUnmarshaller is the entrypoint function for any
// new supported response types.
// E.g. to handle "getBalance" requests, the following need to be defined:
//  1. A new custom responseUnmarshaller
//  2. A new custom struct  to handle the details of the particular response.
type responseUnmarshaller func(logger polylog.Logger, jsonrpcReq jsonrpc.Request, jsonrpcResp jsonrpc.Response) (response, error)

var (
	// All response types needs to implement the response interface.
	// Any new response struct needs to be added to the following list.
	_ response = &responseToGetEpochInfo{}
	_ response = &responseToGetHealth{}
	_ response = &responseGeneric{}

	methodResponseMappings = map[jsonrpc.Method]responseUnmarshaller{
		methodGetHealth:    responseUnmarshallerGetHealth,
		methodGetEpochInfo: responseUnmarshallerGetEpochInfo,
	}
)

func unmarshalResponse(logger polylog.Logger, jsonrpcReq jsonrpc.Request, data []byte) (response, error) {
	var jsonrpcResponse jsonrpc.Response
	err := json.Unmarshal(data, &jsonrpcResponse)
	if err != nil {
		// The response raw payload (e.g. as received from an endpoint) could not be unmarshalled as a JSONRC response.
		// Return a generic response to the user.
		return getGenericJSONRPCErrResponse(logger, jsonrpcReq.ID, data, err), err
	}

	// Validate the JSONRPC response.
	if err := jsonrpcResponse.Validate(jsonrpcReq.ID); err != nil {
		return getGenericJSONRPCErrResponse(logger, jsonrpcReq.ID, data, err), err
	}

	// Note: we intentionally skip checking whether the JSONRPC response indicates an error. This allows the method-specific handler
	// to determine how to respond to the user.
	unmarshaller, found := methodResponseMappings[jsonrpcReq.Method]
	if found {
		return unmarshaller(logger, jsonrpcReq, jsonrpcResponse)
	}

	return responseUnmarshallerGeneric(logger, jsonrpcReq, data)
}
