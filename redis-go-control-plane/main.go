// Copyright 2020 Envoyproxy Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	example "envoy/redis-go-control-plane/example"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/envoyproxy/go-control-plane/pkg/test/v3"
)

var (
	l                 example.Logger
	port              uint
	nodeID            string
	redisclusters     map[string]string
	currentMasterNode string
	masterChanged     = make(chan string)
	ctx               context.Context
	versionNumber     uint64 = 1
)

type MetricsData struct {
	serviceName, address, metricName, value string
}

var metricsData = make(map[string]MetricsData)

func init() {
	l = example.Logger{}

	flag.BoolVar(&l.Debug, "debug", false, "Enable xDS server debug logging")

	// The port that this xDS server listens on
	flag.UintVar(&port, "port", 18000, "xDS management server port")

	// Tell Envoy to use this Node ID
	flag.StringVar(&nodeID, "nodeID", "test-id", "Node ID")
}

func main() {
	flag.Parse()

	go func() {
		for {
			// fmt.Println("Check Redis master node")
			updateEnvoyConfigWithRedisMaster()
			time.Sleep(time.Second * 5) // Polling interval
		}
	}()

	initializeEnvironment()

	// Create a cache
	cache := cache.NewSnapshotCache(false, cache.IDHash{}, l)

	// Create the snapshot that we'll serve to Envoy
	snapshot := example.GenerateSnapshot()
	if err := snapshot.Consistent(); err != nil {
		l.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
		os.Exit(1)
	}
	l.Debugf("will serve snapshot %+v", snapshot)

	// Add the snapshot to the cache
	if err := cache.SetSnapshot(context.Background(), nodeID, snapshot); err != nil {
		l.Errorf("snapshot error %q for %+v", err, snapshot)
		os.Exit(1)
	}

	// Run the xDS server
	go runXdsServer(cache)

	// Export Prometheus metrics
	go MetricsServer()

	// Listen for changes in the Redis master and update Envoy config
	for masterAddr := range masterChanged {
		versionNumber++
		// Call function to update Envoy configuration
		fmt.Printf("++++++Updating Envoy configuration with the new master node: %s\n", masterAddr)
		updateEnvoyClusterConfig(ctx, cache, nodeID, masterAddr, fmt.Sprintf("version-%d", versionNumber))
	}
}

func initializeEnvironment() {
	// Parse Redis node and cluster configuration from environment variables
	redisNodeList := strings.Split(os.Getenv("REDIS_HOST_PORT_LIST"), ",")
	clusterList := strings.Split(os.Getenv("ENVOY_CLUSTERS_LIST"), ",")

	redisclusters = make(map[string]string)
	for i, addr := range redisNodeList {
		redisclusters[addr] = clusterList[i]
	}

	fmt.Printf("redisclusters %+v\n", redisclusters)
}

func updateEnvoyConfigWithRedisMaster() {
	var masterNode string

	countMaster := 0
	for addr, _ := range redisclusters {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			// handle error (could not connect to Redis node)
			metricsData[addr] = MetricsData{"redis", addr, "master", "0"}
			fmt.Printf("Could not connect to Redis node %s\n", addr)
			continue
		}

		// Set a timeout for the connection
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))

		_, err = conn.Write([]byte("INFO replication\r\n"))
		if err != nil {
			// handle error (could not send command to Redis node)
			//delete from metricsData[addr]
			metricsData[addr] = MetricsData{"redis", addr, "master", "0"}
			fmt.Printf("Could not send command to Redis node %s\n", addr)
			conn.Close()
			continue
		}

		isMaster, err := checkIfMaster(conn)
		if err != nil {
			// handle error (could not read or parse response from Redis node)
			metricsData[addr] = MetricsData{"redis", addr, "master", "0"}
			conn.Close()
			continue
		}

		if isMaster {
			countMaster++
			masterNode = addr
			// conn.Close()
			// break
		} else {
			metricsData[addr] = MetricsData{"redis", addr, "master", "0"}
		}

		conn.Close()
	}

	if masterNode != "" && masterNode != currentMasterNode && countMaster == 1 {
		currentMasterNode = masterNode
		fmt.Printf("_Updating Envoy configuration with the new master node: %s\n", masterNode)

		metricsData[masterNode] = MetricsData{"redis", masterNode, "master", "1"}

		// Update Envoy configuration with the new master node
		// Send the master node address to the channel
		masterChanged <- redisclusters[masterNode]
	}
}

func checkIfMaster(conn net.Conn) (bool, error) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "role:master") {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}

func runXdsServer(cache cache.SnapshotCache) {
	ctx = context.Background()
	cb := &test.Callbacks{Debug: l.Debug}
	srv := server.NewServer(ctx, cache, cb)
	example.RunServer(srv, port) // Ensure RunServer is defined
}

func updateEnvoyClusterConfig(ctx context.Context, c cache.SnapshotCache, nodeID, masterNode, version string) error {
	// Create a new cluster with the Redis master as the endpoint
	clusterName := "redis_proxy_cluster"
	// clusterName := "redis2"

	redisCluster := example.MakeRedisCluster(clusterName, masterNode)

	// Create a snapshot with the new cluster configuration
	snap, _ := cache.NewSnapshot(version, map[resource.Type][]types.Resource{
		resource.ClusterType: {redisCluster},

		// Using static cluster
		// resource.EndpointType: {example.MakeRedisEndpoint(clusterName, masterNode)},

		// Using dinamic cluster
		resource.ListenerType: {example.MakeTCPListener(example.ListenerName)},
	})

	// Add the snapshot to the cache
	return c.SetSnapshot(ctx, nodeID, snap)
}

func MetricsServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	http.HandleFunc("/metrics", metricsHandler)
	log.Println("Starting metrics server on :8002")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	var metrics strings.Builder

	// output metricsData
	for addr, data := range metricsData {
		formattedMetric := formatMetric(data.serviceName, addr, data.metricName, data.value)
		metrics.WriteString(formattedMetric + "\n")
	}
	// fmt.Printf("metricsData %+v\n", metricsData)

	fmt.Fprint(w, metrics.String())
}

func formatMetric(serviceName, address, metricName, value string) string {
	// Replace invalid characters in metric names
	metricName = strings.NewReplacer(".", "_", "-", "_").Replace(metricName)

	// Format the metric for Prometheus
	return fmt.Sprintf("%s{service=\"%s\", address=\"%s\"} %s", metricName, serviceName, address, value)
}
