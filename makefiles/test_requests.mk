#######################
#### Test Requests ####
#######################

# This Makefile provides examples of the various ways to make requests to PATH.

# NOTE: All of these requests assume a Shannon Gateway, as the service ID is 'anvil'.

.PHONY: debug_anvil_supplier_info_msg
debug_anvil_supplier_info_msg: ## Displays debugging guidance for Anvil supplier issues
	@echo "#######################################################################################################################################"
	@echo "INFO: If the request did not succeed, look into debugging the Anvil supplier by reviewing:"
	@echo "  https://www.notion.so/buildwithgrove/PATH-Anvil-RelayMiner-Supplier-in-E2E-Test-infrastructure-17da36edfff680da98f2ff01705be00b?pvs=4"
	@echo "########################################################################################################################################"

####################################
#### PATH + Envoy Test Requests ####
####################################

# For all of the below requests:
# - The full PATH stack including Envoy Proxy must be running
# - The 'anvil' service must be configured in the '.config.yaml' file.

# The following are the various ways to make requests to PATH with Envoy running:
# - **Auth**: static API key or no auth (JWT requires a non-expired JWT token, which cannot be statically set)
# - **Service ID**: passed as the subdomain or in the 'target-service-id' header
# - **Endpoint ID**: passed in the URL path or in the 'endpoint-id' header

.PHONY: test_request__endpoint_url_path_mode__no_auth
test_request__endpoint_url_path_mode__no_auth: debug_anvil_supplier_info_msg ## Test request with no auth, endpoint ID passed in the URL path, and the service ID passed as the subdomain
	curl http://anvil.localhost:3070/v1/endpoint_3_no_auth \
		-X POST \
		-d '{"jsonrpc": "2.0", "id": 1, "method": "eth_blockNumber" }'

.PHONY: test_request__endpoint_header_mode__no_auth
test_request__endpoint_header_mode__no_auth: debug_anvil_supplier_info_msg ## Test request with no auth, endpoint ID passed in the endpoint-id header, and the service ID passed as the subdomain
	curl http://anvil.localhost:3070/v1 \
		-X POST \
		-H "endpoint-id: endpoint_3_no_auth" \
		-d '{"jsonrpc": "2.0", "id": 1, "method": "eth_blockNumber" }'

.PHONY: test_request__endpoint_url_path_mode__no_auth__service_id_header
test_request__endpoint_url_path_mode__no_auth__service_id_header: debug_anvil_supplier_info_msg ## Test request with no auth, endpoint ID passed in the URL path, and the service ID passed in the target-service-id header
	curl http://localhost:3070/v1/endpoint_3_no_auth \
		-X POST \
		-H "target-service-id: anvil" \
		-d '{"jsonrpc": "2.0", "id": 1, "method": "eth_blockNumber" }'

.PHONY: test_request__endpoint_header_mode__static_key
test_request__endpoint_header_mode__static_key: debug_anvil_supplier_info_msg ## Test request with static key auth, endpoint ID passed in the endpoint-id header and the service ID passed as the subdomain
	curl http://anvil.localhost:3070/v1 \
		-X POST \
		-H "endpoint-id: endpoint_1_static_key" \
		-H "authorization: api_key_1" \
		-d '{"jsonrpc": "2.0", "id": 1, "method": "eth_blockNumber" }'

.PHONY: test_request__endpoint_url_path_mode__static_key_service_id_header
test_request__endpoint_url_path_mode__static_key_service_id_header: debug_anvil_supplier_info_msg ## Test request with static key auth, endpoint ID passed in the URL path, and the service ID passed in the target-service-id header
	curl http://localhost:3070/v1/endpoint_1_static_key \
		-X POST \
		-H "authorization: api_key_1" \
		-H "target-service-id: anvil" \
		-d '{"jsonrpc": "2.0", "id": 1, "method": "eth_blockNumber" }'

.PHONY: test_request__endpoint_header_mode__static_key_service_id_header
test_request__endpoint_header_mode__static_key_service_id_header: debug_anvil_supplier_info_msg ## Test request with all possible values passed as headers: service ID, endpoint ID and authorization
	curl http://localhost:3070/v1 \
		-X POST \
		-H "endpoint-id: endpoint_1_static_key" \
		-H "authorization: api_key_1" \
		-H "target-service-id: anvil" \
		-d '{"jsonrpc": "2.0", "id": 1, "method": "eth_blockNumber" }'

############################
#### PATH Test Requests ####
############################

.PHONY: test_request_path_only
test_request_path_only: debug_anvil_supplier_info_msg ## Test request against the PATH Gateway running on port 3069 without Envoy Proxy
	curl http://localhost:3069/v1/ \
		-X POST \
		-H "Content-Type: application/json" \
		-H "target-service-id: anvil" \
		-d '{"jsonrpc": "2.0", "id": 1, "method": "eth_blockNumber" }'