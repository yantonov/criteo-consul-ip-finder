package main

import (
	"consul-ip-finder/cmd/cli"
	"consul-ip-finder/lib"
	"log"
)

func main() {

	parameters, cmdParamErr := cli.ParseParameters()
	if cmdParamErr != nil {
		log.Fatal(cmdParamErr)
	}

	services, err := lib.FindService(parameters.Ip, parameters.Datacenter, parameters.Environment)
	if err != nil {
		log.Fatal(err)
	}
	if len(services) == 0 {
		log.Fatal("No services found for ip=" + parameters.Ip)
	}
	println("Found services:")
	for _, service := range services {
		println(service)
	}
}
