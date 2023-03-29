//go:build windows2012R2
// +build windows2012R2

package main

import (
	"encoding/json"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"code.cloudfoundry.org/diego-ssh/server"
	"code.cloudfoundry.org/lager/v3"
)

type PortMapping struct {
	Internal int `json:"internal"`
	External int `json:"external"`
}

func createServer(
	logger lager.Logger,
	address string,
	sshDaemon server.ConnectionHandler,
) (*server.Server, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	jsonPortMappings := os.Getenv("CF_INSTANCE_PORTS")
	var portMappings []PortMapping
	json.Unmarshal([]byte(jsonPortMappings), &portMappings)
	for _, mapping := range portMappings {
		if strconv.Itoa(mapping.Internal) == port {
			port = strconv.Itoa(mapping.External)
		}
	}
	address = strings.Join([]string{host, port}, ":")
	return server.NewServer(logger, address, sshDaemon, 5*time.Minute), err
}
