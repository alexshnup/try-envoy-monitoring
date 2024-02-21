package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type EnvoyClusters struct {
	ClusterStatuses []ClusterStatus `json:"cluster_statuses"`
}

type ClusterStatus struct {
	Name              string       `json:"name"`
	HostStatuses      []HostStatus `json:"host_statuses"`
	ObservabilityName string       `json:"observability_name"`
}

type HostStatus struct {
	Address      Address                `json:"address"`
	Stats        []Stat                 `json:"stats"`
	HealthStatus map[string]interface{} `json:"health_status"`
	Weight       int                    `json:"weight"`
	Hostname     string                 `json:"hostname"`
}

type Address struct {
	SocketAddress SocketAddress `json:"socket_address"`
}

type SocketAddress struct {
	Address   string `json:"address"`
	PortValue int    `json:"port_value"`
}

type Stat struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

func main() {
	http.HandleFunc("/metrics", metricsHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	// Replace this with the actual URL of your Envoy admin console
	envoyURL := "http://localhost:8001/clusters?format=json"

	resp, err := http.Get(envoyURL)
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	var envoyClusters EnvoyClusters
	if err := json.Unmarshal(body, &envoyClusters); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Error unmarshalling JSON", http.StatusInternalServerError)
		return
	}

	var metrics strings.Builder
	for _, cluster := range envoyClusters.ClusterStatuses {
		for _, host := range cluster.HostStatuses {
			for _, stat := range host.Stats {
				metricName := fmt.Sprintf("envoy_%s_%s", cluster.Name, stat.Name)
				metricValue := 1 // Assuming a dummy value. Replace with actual value if available in your JSON.
				metrics.WriteString(fmt.Sprintf("%s{cluster=\"%s\", host=\"%s\"} %d\n", metricName, cluster.Name, host.Hostname, metricValue))
			}
		}
	}

	fmt.Fprint(w, metrics.String())
}
