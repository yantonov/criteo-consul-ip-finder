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

	ch := make(chan string, len(services))

	wg := sync.WaitGroup{}
	for _, service := range services {
		wg.Add(1)
		go inspectService(client, service, ip, ch, &wg)
	}
	wg.Wait()
	close(ch)

	var result []string

	for foundService := range ch {
		if foundService != "" {
			result = append(result, foundService)
		}
	}
	return result, nil
}

func inspectService(client consul.Client, serviceName string, ip string, ch chan string, wg *sync.WaitGroup) {
	serviceInfo, err := consul.GetService(client, serviceName)
	found := false
	if err == nil {
		for _, instance := range serviceInfo.Instances {
			fmt.Printf("service=%s instance with address=%s\n", serviceName, instance.ServiceAddress)
			if instance.ServiceAddress == ip {
				fmt.Printf("Found service=%s with address=%s\n", serviceName, instance.ServiceAddress)
				ch <- serviceName
				found = true
				break
			}
		}
	}
	if !found {
		ch <- ""
	}
	wg.Done()
}
