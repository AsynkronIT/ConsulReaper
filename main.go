package main

import (
	"log"
	"time"

	"os"

	"github.com/hashicorp/consul/api"
	"fmt"
)

func Reap(address string, datacenter string) {

	client, err := api.NewClient(&api.Config{
		Address: address,
	})
	if err != nil {
		log.Println("Cound not create Consul client")
		log.Println(err.Error())
		return
	}
	for {
		log.Println("Getting critical services")
		criticalServices, _, err := client.Health().State("critical", &api.QueryOptions{
			Datacenter:        datacenter,
			AllowStale:        false,
			RequireConsistent: false,
			WaitTime:          5 * time.Second,
		})

		if err != nil {
			log.Println("Cound not get service health")
			log.Println(err.Error())
			return
		}

		log.Println("Got critical services")
		for _, critical := range criticalServices {
			log.Println("Deregistering " + critical.ServiceID)

			err = client.Agent().ServiceDeregister(critical.ServiceID)
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
	address := fmt.Sprintf("%v:%v",host,port)
	time.Sleep(5 * time.Second)
	log.Println("Connecting to Consul on " + address)
	Reap(address, "")
}
