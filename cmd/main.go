package main

import (
	"consul-ip-finder/cmd/cli"
	"consul-ip-finder/cmd/ui"
	"consul-ip-finder/lib"
	"log"
)

func main() {

	parameters, cmdParamErr := cli.ParseParameters()
	if cmdParamErr != nil {
		log.Fatal(cmdParamErr)
	}

	var progressBar ui.ProgressBar = nil

	if parameters.Verbose {
		progressBar = &ui.EmptyProgressBar{}
	} else {
		progressBar = &ui.TerminalWidgetProgressBar{}
	}
	services, err := lib.FindService(
		parameters.Ip,
		parameters.Datacenter,
		parameters.Environment,
		parameters.ParallelismLevel,
		parameters.Verbose,
		progressBar)
	if err != nil {
		log.Fatal(err)
	}
	if len(services) == 0 {
		println("No services found for ip=" + parameters.Ip)
	} else {
		println("Found services:")
		for _, service := range services {
			println(service)
		}
	}
}
