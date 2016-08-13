package main

import (
	"log"
	"time"

	"os"

	"github.com/hashicorp/consul/api"
	"fmt"
)

func do(body func() error) {
	for {
		err := body()
		if err != nil {
			log.Println(err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		return
	}
}

func Reap(address string, datacenter string) {
	var client *api.Client
	do(func() error {
		c, err := api.NewClient(&api.Config{Address: address})
		client = c
		return err
	})

	for {
		log.Println("Getting critical services")
		var criticalServices []*api.HealthCheck
		do(func() error {
			c, _, err := client.Health().State("critical", &api.QueryOptions{
				Datacenter:        datacenter,
				AllowStale:        false,
				RequireConsistent: false,
				WaitTime:          5 * time.Second,
			})
			criticalServices = c
			return err
		})

		log.Println("Got critical services")
		for _, critical := range criticalServices {
			log.Println("Deregistering " + critical.ServiceID)

			err := client.Agent().ServiceDeregister(critical.ServiceID)
			if err != nil {
				log.Println("Cound not deregister service on agent" + critical.ServiceID)
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func main() {
	log.Println("Starting...")
	host := os.Getenv("ConsulHost")
	if host == "" {
		host = "127.0.0.1"
	}
	port := "8500"
	address := fmt.Sprintf("%v:%v", host, port)
	log.Println("Connecting to Consul on " + address)
	Reap(address, "")
}
