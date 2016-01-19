package main

import (
	"github.com/hashicorp/consul/api"
	"log"
	"os"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	var (
		q  *api.QueryOptions
		nc nagiosConfig
	)
	nc.CfgPath = "/tmp/nagios.cfg"
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}
	catalog := client.Catalog()
	health := client.Health()
	nodes, _, err := catalog.Nodes(q)
	if err != nil {
		log.Fatal(err)
	}
	for _, node := range nodes {
		nh := nagiosHost{
			Name: node.Node,
		}
		nc.addHost(&nh)
		healthChecks, _, err := health.Node(node.Node, q)
		if err != nil {
			log.Fatal(err)
		}
		for _, healthCheck := range healthChecks {
			ns := nagiosService{
				Description: healthCheck.Name,
				Host:        &nh,
			}
			nc.addService(&ns)
		}
	}
	err = nc.render()
	if err != nil {
		log.Fatal(err)
	}
	return 0

}
