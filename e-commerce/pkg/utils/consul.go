package utils

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"os"
	"strconv"
)

func ConsulClient() *api.Client {
	consulAddr := os.Getenv("CONSUL_ADDR")
	if consulAddr == "" {
		consulAddr = "consul:8500"
	}

	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulAddr
	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	return client
}

func RegistrationService() *api.AgentServiceRegistration {
	serviceID := os.Getenv("SERVICE_ID")
	serviceName := os.Getenv("SERVICE_NAME")
	servicePortStr := os.Getenv("PORT")
	port, err := strconv.Atoi(servicePortStr)
	if err != nil {
		log.Fatalf("Failed to convert SERVICE_PORT to int: %v", err)
		return nil
	}
	serviceAddress := os.Getenv("SERVICE_ADDRESS")
	if serviceAddress == "" {
		serviceAddress = serviceName
	}

	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Port:    port,
		Address: serviceAddress,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", serviceAddress, port),
			Interval: "10s",
			Timeout:  "3s",
		},
	}
	return registration
}

func RegisterToConsul(client *api.Client, registration *api.AgentServiceRegistration) {
	err := client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("Failed to register to Consul: %v", err.Error())
		return
	}
}
