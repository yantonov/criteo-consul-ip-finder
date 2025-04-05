package lib

import (
	"consul-ip-finder/cmd/ui"
	"consul-ip-finder/lib/consul"
	"fmt"
	"log"
	"sync"
)

func FindService(
	ip string,
	dc string, env string,
	parallelismLevel int,
	verbose bool,
	bar ui.ProgressBar) ([]string, error) {

	client := consul.Create(dc, env, verbose)
	services, err := consul.GetListOfServices(client)
	if err != nil {
		return nil, fmt.Errorf("error getting services from Consul: %v", err)
	}

	numberOfServices := len(services)
	println(fmt.Sprintf("Total number of services=%d", numberOfServices))

	resultChannel := make(chan string, numberOfServices)

	parallelismLevelChannel := make(chan int, parallelismLevel)

	bar.Init(numberOfServices)

	wg := sync.WaitGroup{}
	for _, service := range services {
		wg.Add(1)
		parallelismLevelChannel <- 1
		go inspectService(client, service, ip, resultChannel, parallelismLevelChannel, verbose, bar, &wg)
	}

	wg.Wait()
	close(resultChannel)

	return extractFoundServices(numberOfServices, resultChannel)
}

func extractFoundServices(numberOfServices int,
	ch chan string) ([]string, error) {
	var result []string

	for _ = range numberOfServices {
		foundService := <-ch
		if foundService != "" {
			result = append(result, foundService)
		}
	}
	return result, nil
}

func inspectService(
	client consul.Client,
	serviceName string,
	ip string,
	resultChannel chan string,
	parallelismLevelChannel chan int, verbose bool, bar ui.ProgressBar,
	wg *sync.WaitGroup) {

	defer wg.Done()

	if verbose {
		println(serviceName)
	}
	serviceInfo, err := consul.GetService(client, serviceName)
	found := false
	if err != nil {
		if verbose {
			log.Println("Error getting service info from Consul:", err)
		}
	} else {
		for _, instance := range serviceInfo.Instances {
			if verbose {
				fmt.Printf("service=%s instance with address=%s\n", serviceName, instance.ServiceAddress)
			}
			if instance.ServiceAddress == ip {
				if verbose {
					fmt.Printf("Found service=%s with address=%s\n", serviceName, instance.ServiceAddress)
				}
				resultChannel <- serviceName
				found = true
				break
			}
		}
	}
	if !found {
		resultChannel <- ""
	}
	bar.Add(1)
	<-parallelismLevelChannel
}
