package cli

import (
	"errors"
	"flag"
	"fmt"
)

type CmdParameters struct {
	Datacenter       string
	Environment      string
	Ip               string
	ParallelismLevel int
	Verbose          bool
}

const MaxNumberOfConcurrentConnections = 50
const DefaultNumberOfConcurrentConnections = 20

func ParseParameters() (*CmdParameters, error) {
	ipParam := flag.String("ip", "", "ip address")
	dcParam := flag.String("dc", "", "data center code")
	environmentParam := flag.String("env", "", "environment")
	threadsParam := flag.Int("threads",
		DefaultNumberOfConcurrentConnections,
		fmt.Sprintf("number of simultaneous connections, default=%d", DefaultNumberOfConcurrentConnections))
	verboseParam := flag.Bool("v", false, "verbose mode")
	flag.Parse()
	if *ipParam == "" || *dcParam == "" || *environmentParam == "" {
		flag.Usage()
		return nil, errors.New("ip and datacenter and environment are required")
	}
	if *threadsParam > MaxNumberOfConcurrentConnections {
		return nil, errors.New(fmt.Sprintf("Maximum number of concurrent connections = %d", MaxNumberOfConcurrentConnections))
	}
	if *threadsParam < 0 {
		return nil, errors.New("threads=number of simultaneous connections must be a positive number")
	}
	return &CmdParameters{
		Ip:               *ipParam,
		Datacenter:       *dcParam,
		Environment:      *environmentParam,
		ParallelismLevel: *threadsParam,
		Verbose:          *verboseParam,
	}, nil
}
