package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"fmt"

	"github.com/hashicorp/consul/api"
)

func Reap(address string, datacenter string) {
	client, err := api.NewClient(&api.Config{Address: address})
	if err != nil {
		log.Fatalf("Failed: Creating client %v", err)
	}
	node, err := client.Agent().NodeName()
	if err != nil {
		log.Fatalf("Getting Node name %v", err)
	}

	for {
		var criticalServices []*api.HealthCheck
		criticalServices, _, err := client.Health().State("critical", &api.QueryOptions{
			Datacenter:        datacenter,
			AllowStale:        false,
			RequireConsistent: false,
			WaitTime:          5 * time.Second,
		})
		if err != nil {
			log.Fatalf("Failed: Getting critical services %v", err)
		}

		for _, critical := range criticalServices {
			if critical.Node == node {
				log.Println("Deregistering " + critical.ServiceID)

				err := client.Agent().ServiceDeregister(critical.ServiceID)
				if err != nil {
					log.Println("Cound not deregister service on agent" + critical.ServiceID)
				}
			} else {
				log.Printf("%v %v", node, critical.Node)
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func main() {
	log.Println("Starting...")
	time.Sleep(20 * time.Second)
	host, _ := httpGet("http://rancher-metadata.rancher.internal/2015-12-19/self/host/agent_ip")
	log.Printf("Got Rancher Host IP %v", host)
	port := "8500"
	address := fmt.Sprintf("%v:%v", host, port)
	log.Println("Connecting to Consul on " + address)
	Reap(address, "")
}

func httpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed: Get %v %v", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	str := string(body)
	return str, err
}
