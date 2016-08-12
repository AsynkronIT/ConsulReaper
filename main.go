package main

import (
	"log"
	"time"

	"github.com/hashicorp/consul/api"
)

func Reap(datacenter string) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Println("Cound not create Consul client")
		return
	}

	criticalServices, _, err := client.Health().State("critical", &api.QueryOptions{
		Datacenter:        datacenter,
		AllowStale:        true,
		RequireConsistent: false,
		WaitTime:          5 * time.Second,
	})

	if err != nil {
		log.Println("Cound not get service health")
		panic(err)
		return
	}
	for _, critical := range criticalServices {
		_, err := client.Catalog().Deregister(&api.CatalogDeregistration{
			Datacenter: datacenter,
			Node:       critical.Node,
			ServiceID:  critical.ServiceID,
		}, &api.WriteOptions{})
		if err != nil {
			log.Println("Cound not deregister " + critical.ServiceID)
		}
	}
}

func main() {
	Reap("")
}
