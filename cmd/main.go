package main

import (
	"consul-ip-finder/lib"
	"errors"
	"flag"
	"log"
)

type CmdParameters struct {
	datacenter  string
	environment string
	ip          string
}

func ParseParameters() (*CmdParameters, error) {
	ipParam := flag.String("ip", "", "ip address")
	dcParam := flag.String("datacenter", "", "data center code")
	environmentParam := flag.String("environment", "", "environment")
	flag.Parse()
	if *ipParam == "" || *dcParam == "" || *environmentParam == "" {
		flag.Usage()
		return nil, errors.New("ip or datacenter or environment required")
	}
	return &CmdParameters{
		ip:          *ipParam,
		datacenter:  *dcParam,
		environment: *environmentParam,
	}, nil
}

func main() {

	parameters, cmdParamErr := ParseParameters()
	if cmdParamErr != nil {
		log.Fatal(cmdParamErr)
	}

	services, err := lib.FindService(parameters.ip, parameters.datacenter, parameters.environment)
	if err != nil {
		log.Fatal(err)
	}
	if len(services) == 0 {
		log.Fatal("No services found for ip=" + parameters.ip)
	}
	println("Found services:")
	for _, service := range services {
		println(service)
	}
}
