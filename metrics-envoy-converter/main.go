package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/metrics", metricsHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	// Sample data - replace this with the actual data source or a file read operation
	// 	sampleData := `
	// redis::192.168.160.4:6379::rq_total::0
	// redis::192.168.160.4:6379::hostname::redis-replica1
	// redis::192.168.160.4:6379::health_flags::/failed_active_hc/active_hc_timeout
	// redis::192.168.160.4:6379::weight::1
	// `
	// Replace this with the actual URL of your Envoy admin console
	// envoyURL := "http://localhost:8001/clusters"

	//get envoyurl from url form
	envoyURL := r.URL.Query()["envoyurl"]

	// fmt.Println(envoyURL)

	if envoyURL == nil {
		envoyURL = []string{"http://localhost:8001/clusters"}
	}

	resp, err := http.Get(envoyURL[0])
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

	scanner := bufio.NewScanner(strings.NewReader(string(body)))
	var metrics strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "::")

		var formattedMetric string
		switch len(parts) {
		case 4: // Format: [service name]::[address]::[metric/attribute name]::[value]
			serviceName, address, metricName, value := parts[0], parts[1], parts[2], parts[3]
			switch value {
			case "":
				continue

			case "false":
				value = "0"
			case "true":
				value = "1"

			case "healthy":
				metricName = "master"
				value = "1"
			case "/failed_active_hc/active_hc_timeout":
				metricName = "master"
				value = "0"
			case "/pending_dynamic_removal":
				metricName = "master"
				value = "0"

			case "/failed_active_hc/pending_dynamic_removal/active_hc_timeout":
				metricName = "master"
				value = "0"
			}

			switch metricName {
			case "hostname":
				continue
			}

			formattedMetric = formatMetric(serviceName, address, metricName, value)

		case 3: // Format: [service name]::[metric/attribute name]::[value]
			fmt.Println(parts)
			serviceName, metricName, value := parts[0], parts[1], parts[2]
			formattedMetric = formatMetricWithoutAddress(serviceName, metricName, value)
			continue

		case 2: // Format: [service name]::[metric/attribute name]::[value]
			fmt.Println(parts)
			serviceName, value := parts[0], parts[1]
			formattedMetric = formatMetricWithoutAddress(serviceName, "", value)
			continue

		default:
			continue
		}

		metrics.WriteString(formattedMetric + "\n")
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading data: %v", err)
		http.Error(w, "Error reading data", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, metrics.String())
}

func convertToNumericValue(value string) string {
	// Check if the value is a boolean
	if value == "true" {
		return "1"
	} else if value == "false" {
		return "0"
	}
	// If not a boolean, return the value as-is
	return value
}

func formatMetric(serviceName, address, metricName, value string) string {
	// Replace invalid characters in metric names
	metricName = strings.NewReplacer(".", "_", "-", "_").Replace(metricName)

	// Convert boolean values to numeric
	numericValue := convertToNumericValue(value)

	// Format the metric for Prometheus
	return fmt.Sprintf("%s{service=\"%s\", address=\"%s\"} %s", metricName, serviceName, address, numericValue)
}

func formatMetricWithoutAddress(serviceName, metricName, value string) string {
	// Replace invalid characters in metric names
	metricName = strings.NewReplacer(".", "_", "-", "_").Replace(metricName)

	// Format the metric for Prometheus without address
	return fmt.Sprintf("%s{service=\"%s\"} %s", metricName, serviceName, value)
}
