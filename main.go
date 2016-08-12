package main

import (
	"log"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/rogeralsing/goconsole"
	"os"
)

func Reap(address string, datacenter string) {

	client, err := api.NewClient(&api.Config{
		Address: address,
	})
	if err != nil {
		log.Println("Cound not create Consul client")
		return
	}
	for ; ; {
		log.Println("Getting critical services")
		criticalServices, _, err := client.Health().State("critical", &api.QueryOptions{
			Datacenter:        datacenter,
			AllowStale:        false,
			RequireConsistent: false,
			WaitTime:          5 * time.Second,
		})

		if err != nil {
			log.Println("Cound not get service health")
			return
		}

		log.Println("Got critical services")
		for _, critical := range criticalServices {
			log.Println("Deregistering " + critical.ServiceID)
			_, err := client.Catalog().Deregister(&api.CatalogDeregistration{
				Datacenter: datacenter,
				Node:       critical.Node,
				ServiceID:  critical.ServiceID,
			}, &api.WriteOptions{})
			if err != nil {
				log.Println("Cound not deregister " + critical.ServiceID)
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func main() {
	address := os.Getenv("ConsulAddress")
	if address == "" {
		address = "127.0.0.1:8500"
	}
	go Reap("", "")
	console.ReadLine()
}
