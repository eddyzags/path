package evm

import (
	"encoding/json"

	"github.com/pokt-network/poktroll/pkg/polylog"

	qosobservations "github.com/buildwithgrove/path/observation/qos"
	"github.com/buildwithgrove/path/qos/jsonrpc"
)

// responseToChainID provides the functionality required from a response by a requestContext instance.
var _ response = responseToChainID{}

// TODO_TECHDEBT(@adshmh): consider refactoring all unmarshallers to remove any duplicated logic.

// responseUnmarshallerChainID deserializes the provided byte slice into a responseToChainID struct,
// adding any encountered errors to the returned struct for constructing a response payload.
func responseUnmarshallerChainID(
	logger polylog.Logger,
	jsonrpcReq jsonrpc.Request,
	jsonrpcResp jsonrpc.Response,
) (response, error) {
	// The endpoint returned an error: no need to do further processing of the response.
	if jsonrpcResp.IsError() {
		return responseToChainID{
			logger:          logger,
			jsonRPCResponse: jsonrpcResp,

			// DEV_NOTE: A valid JSONRPC error response is considered a valid response.
			valid: true,
		}, nil
	}

	resultBz, err := jsonrpcResp.GetResultAsBytes()
	if err != nil {
		validationError := qosobservations.EVMResponseValidationError_EVM_RESPONSE_VALIDATION_ERROR_UNMARSHAL
		return responseToChainID{
			logger:          logger,
			jsonRPCResponse: jsonrpcResp,
			validationError: &validationError,
		}, err
	}

	var result string
	validationError := qosobservations.EVMResponseValidationError_EVM_RESPONSE_VALIDATION_ERROR_UNSPECIFIED

	err = json.Unmarshal(resultBz, &result)
	if err != nil {
		validationError = qosobservations.EVMResponseValidationError_EVM_RESPONSE_VALIDATION_ERROR_UNMARSHAL
	}

	return &responseToChainID{
		logger:          logger,
		jsonRPCResponse: jsonrpcResp,
		result:          result,

		// if unmarshaling succeeded, the response is considered valid.
		valid: (err == nil),

		validationError: &validationError,
	}, err
}

// responseToChainID captures the fields expected in a
// response to an `eth_chainId` request.
type responseToChainID struct {
	logger polylog.Logger

	// jsonRPCResponse stores the JSONRPC response parsed from an endpoint's response bytes.
	jsonRPCResponse jsonrpc.Response

	// result captures the `result` field of a JSONRPC response to an `eth_chainId` request.
	result string

	// valid is set to true if the parsed response is deemed valid.
	// As of PR #152, a response is valid if either of the following holds:
	//	- It is a valid JSONRPC error response
	//	- It is a valid JSONRPC response with any string value in `result` field.
	valid bool

	// Why the response has failed validation.
	// Used when generating observations.
	validationError *qosobservations.EVMResponseValidationError
}

// GetObservation returns an observation of the endpoint's response to an `eth_chainId` request.
// Implements the response interface.
// Reference: https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_chainid
func (r responseToChainID) GetObservation() qosobservations.EVMEndpointObservation {
	return qosobservations.EVMEndpointObservation{
		ResponseObservation: &qosobservations.EVMEndpointObservation_ChainIdResponse{
			ChainIdResponse: &qosobservations.EVMChainIDResponse{
				ChainIdResponse:         r.result,
				Valid:                   r.valid,
				ResponseValidationError: r.validationError,
			},
		},
	}
}

// TODO_MVP(@adshmh): handle the following scenarios:
//  1. An endpoint returned a malformed, i.e. Not in JSONRPC format, response.
//     The user-facing response should include the request's ID.
//  2. An endpoint returns a JSONRPC response indicating a user error:
//     This should be returned to the user as-is.
//  3. An endpoint returns a valid JSONRPC response to a valid user request:
//     This should be returned to the user as-is.
//
// GetResponsePayload returns the raw byte slice payload to be returned as the response to the JSONRPC request.
// It implements the response interface.
func (r responseToChainID) GetResponsePayload() []byte {
	bz, err := json.Marshal(r.jsonRPCResponse)
	if err != nil {
		// This should never happen: log an entry but return the response anyway.
		r.logger.Warn().Err(err).Msg("responseToChainID: Marshaling JSONRPC response failed.")
	}
	return bz
}
