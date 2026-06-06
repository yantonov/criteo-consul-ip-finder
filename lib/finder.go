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
	log.Printf("Total number of services=%d", numberOfServices)

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

	return extractFoundServices(resultChannel)
}

func extractFoundServices(ch chan string) ([]string, error) {
	var result []string

	for foundService := range ch {
		result = append(result, foundService)
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

	defer func() {
		bar.Add(1)
		<-parallelismLevelChannel
		wg.Done()
	}()

	if verbose {
		log.Println(serviceName)
	}
	serviceInfo, err := consul.GetService(client, serviceName)
	if err != nil {
		if verbose {
			log.Println("Error getting service info from Consul:", err)
		}
	} else {
		for _, instance := range serviceInfo.Instances {
			if verbose {
				log.Printf("service=%s instance with address=%s\n", serviceName, instance.ServiceAddress)
			}
			if instance.ServiceAddress == ip {
				if verbose {
					log.Printf("Found service=%s with address=%s\n", serviceName, instance.ServiceAddress)
				}
				resultChannel <- serviceName
				break
			}
		}
	}
}
