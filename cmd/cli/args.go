package cli

import (
	"errors"
	"flag"
)

type CmdParameters struct {
	Datacenter  string
	Environment string
	Ip          string
}

func ParseParameters() (*CmdParameters, error) {
	ipParam := flag.String("ip", "", "ip address")
	dcParam := flag.String("dc", "", "data center code")
	environmentParam := flag.String("env", "", "environment")
	flag.Parse()
	if *ipParam == "" || *dcParam == "" || *environmentParam == "" {
		flag.Usage()
		return nil, errors.New("ip and datacenter and environment are required")
	}
	return &CmdParameters{
		Ip:          *ipParam,
		Datacenter:  *dcParam,
		Environment: *environmentParam,
	}, nil
}
