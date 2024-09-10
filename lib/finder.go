package lib

import (
	"consul-ip-finder/lib/consul"
	"fmt"
	"sync"
)

func FindService(ip string, dc string, env string) ([]string, error) {
	client := consul.Create(dc, env)
	services, err := consul.GetListOfServices(client)
	if err != nil {
		return nil, fmt.Errorf("error getting services from Consul: %v", err)
	}

	println(fmt.Sprintf("Total number of service=%d", len(services)))

	resultChannel := make(chan string, len(services))

	parallelismLevel := 20
	parallelismLevelChannel := make(chan int, parallelismLevel)

	wg := sync.WaitGroup{}
	for _, service := range services {
		wg.Add(1)
		parallelismLevelChannel <- 1
		go inspectService(client, service, ip, resultChannel, parallelismLevelChannel, &wg)
	}
	wg.Wait()
	close(resultChannel)

	return extractFoundServices(resultChannel)
}

func extractFoundServices(ch chan string) ([]string, error) {
	var result []string

	for foundService := range ch {
		if foundService != "" {
			result = append(result, foundService)
		}
	}
	return result, nil
}

func inspectService(client consul.Client, serviceName string, ip string, resultChannel chan string, parallelismLevelChannel chan int, wg *sync.WaitGroup) {
	<-parallelismLevelChannel
	serviceInfo, err := consul.GetService(client, serviceName)
	found := false
	if err == nil {
		for _, instance := range serviceInfo.Instances {
			fmt.Printf("service=%s instance with address=%s\n", serviceName, instance.ServiceAddress)
			if instance.ServiceAddress == ip {
				fmt.Printf("Found service=%s with address=%s\n", serviceName, instance.ServiceAddress)
				resultChannel <- serviceName
				found = true
				break
			}
		}
	}
	if !found {
		resultChannel <- ""
	}
	wg.Done()
}
