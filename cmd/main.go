package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pokt-network/poktroll/pkg/polylog"
	"github.com/pokt-network/poktroll/pkg/polylog/polyzero"

	"github.com/pokt-foundation/portal-middleware/config"
	"github.com/pokt-foundation/portal-middleware/gateway"
	"github.com/pokt-foundation/portal-middleware/relayer"
	"github.com/pokt-foundation/portal-middleware/relayer/morse"
	"github.com/pokt-foundation/portal-middleware/relayer/shannon"
	"github.com/pokt-foundation/portal-middleware/request"
	"github.com/pokt-foundation/portal-middleware/router"
)

const configPath = ".config.yaml"

func getProtocol(config config.GatewayConfig, logger polylog.Logger) (relayer.Protocol, error) {

	// Config YAML validation enforces that exactly one protocol config is set,
	// so first check if the protocol config is set for Shannon.
	if shannonConfig := config.GetShannonConfig(); shannonConfig != nil {
		logger.Info().Msg("Starting PATH gateway with Shannon protocol")

		fullNode, err := shannon.NewFullNode(shannonConfig.FullNodeConfig, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create shannon full node: %v", err)
		}

		protocol, err := shannon.NewProtocol(context.Background(), fullNode)
		if err != nil {
			return nil, fmt.Errorf("failed to create shannon protocol: %v", err)
		}

		return protocol, nil
	}

	// If the protocol config is not set for Shannon, then it must be set for Morse.
	if morseConfig := config.GetMorseConfig(); morseConfig != nil {
		logger.Info().Msg("Starting PATH gateway with Morse protocol")

		fullNode, err := morse.NewFullNode(morseConfig.FullNodeConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create morse full node: %v", err)
		}

		protocol, err := morse.NewProtocol(context.Background(), fullNode, morseConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create morse protocol: %v", err)
		}

		return protocol, nil
	}

	// this should never happen but guard against it
	return nil, fmt.Errorf("no protocol config set")
}

func main() {
	logger := polyzero.NewLogger()

	config, err := config.LoadGatewayConfigFromYAML(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	protocol, err := getProtocol(config, logger)
	if err != nil {
		log.Fatalf("failed to create protocol: %v", err)
	}

	requestParser, err := request.NewParser(config, logger)
	if err != nil {
		log.Fatalf("failed to create request parser: %v", err)
	}

	relayer := &relayer.Relayer{Protocol: protocol}

	gateway := &gateway.Gateway{
		HTTPRequestParser: requestParser,
		Relayer:           relayer,
	}

	apiRouter := router.NewRouter(gateway, config.GetRouterConfig(), logger)
	if err != nil {
		log.Fatalf("failed to create API router: %v", err)
	}

	if err := apiRouter.Start(context.Background()); err != nil {
		log.Fatalf("failed to start API router: %v", err)
	}
}
