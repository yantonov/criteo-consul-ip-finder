package consul

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
)

type Client struct {
	URL string
}

type ServiceInstanceResponse struct {
	ServiceAddress string `json:"ServiceAddress"`
}

type ServiceResponse struct {
	Instances []ServiceInstanceResponse
}

func Create(dc string, env string) Client {
	return Client{
		URL: getConsulURL(dc, env),
	}
}

func getConsulURL(dc string, env string) string {
	return fmt.Sprintf("https://consul-relay.%s.%s.crto.in", dc, env)
}

func GetListOfServices(client Client) ([]string, error) {
	url := fmt.Sprintf("%s/v1/catalog/services", client.URL)

	println(url)

	consulClient := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, getErr := consulClient.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	var response interface{}
	err = json.Unmarshal(body, &response)
	serviceMap := response.(map[string]interface{})
	if err != nil {
		return nil, err
	}

	serviceNames := make([]string, len(serviceMap))
	i := 0
	for k := range serviceMap {
		serviceNames[i] = k
		i++
	}
	sort.Strings(serviceNames)
	return serviceNames, nil
}

func GetService(client Client, serviceName string) (*ServiceResponse, error) {
	url := fmt.Sprintf("%s/v1/catalog/service/%s", client.URL, serviceName)

	consulClient := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, getErr := consulClient.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	response := ServiceResponse{}
	err = json.Unmarshal(body, &response.Instances)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
