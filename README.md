## Project Overview

This project integrates Envoy Proxy with a Go-Control-Plane to dynamically manage configurations for a Redis cluster. The setup includes TCP health checks for Redis nodes and dynamic updates of Envoy's configuration in response to changes in the Redis cluster. The stack also features Grafana for visualization, VictoriaMetrics for monitoring, and a custom control plane for Envoy.

## Stack Components

-   **Envoy**: Front proxy that routes traffic to the Redis cluster. Configured to dynamically update its configuration based on the Redis nodes' status.
-   **Redis Cluster**: Includes one master (`redis-master`) and two replicas (`redis-replica1`, `redis-replica2`). Redis Sentinels (`redis-sentinel1`, `redis-sentinel2`, `redis-sentinel3`) are used for high availability.
-   **Grafana**: Visualization tool for monitoring metrics. Integrated with VictoriaMetrics.
-   **VictoriaMetrics**: Time-series database used for storing and querying metrics. `vmagent` is used for scraping and sending metrics to VictoriaMetrics.
-   **Go-Control-Plane (`redis-go-control-plane`)**: Custom implementation of the control plane for Envoy, handling dynamic configuration updates.

## Getting Started

To get started, clone the repository and navigate to the project directory. Use `docker-compose up --build` to build and run the containers. This will set up the entire stack, including Envoy, Redis, Grafana, VictoriaMetrics, and the Go-Control-Plane.

## Configuration

-   **Envoy**: Configured via `./envoy.yaml` to interact with the Go-Control-Plane.
-   **Redis**: Master-slave configuration with Sentinel for failover. Custom scripts can be added for cluster setup.
-   **Grafana**: Provisioning configurations are provided in `./grafana/provisioning`. The default admin password is set to `secret`.
-   **VictoriaMetrics**: Time-series database, configured via `victoriametrics.yml`.
-   **Go-Control-Plane**: Custom control plane for Envoy, configured through environment variables `REDIS_HOST_PORT_LIST` and `ENVOY_CLUSTERS_LIST`.

## Health Check Implementation

TCP health checks are set up for Redis nodes to monitor their status. The Go-Control-Plane sends `PING` commands to the Redis nodes and expects a `PONG` response to confirm their status.

## Usage

Once the setup is running, Envoy will dynamically adjust its routing based on the health and status of the Redis nodes as provided by the Go-Control-Plane. Grafana can be used to visualize metrics collected by VictoriaMetrics.




# Envoy Dynamic Configuration with Go-Control-Plane

This project demonstrates the integration of Envoy Proxy with the Go-Control-Plane to dynamically manage configurations, particularly for a Redis cluster. The focus is on implementing TCP health checks to monitor the status of Redis nodes and update Envoy's configuration in response to changes in the Redis cluster.

## Overview

Envoy is used as a front proxy in this setup, routing traffic to a Redis cluster. The Go-Control-Plane serves as an xDS server to dynamically update Envoy's configuration, especially for handling changes in the Redis master node.

## Features

-   **Dynamic Redis Cluster Management**: Automatically updates Envoy's configuration when the Redis master node changes.
-   **TCP Health Checks**: Implements health checks for Redis nodes to monitor their status.
-   **Efficient Configuration Updates**: Utilizes Envoy's xDS APIs for seamless and efficient configuration updates without restarting Envoy.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

-   Docker
-   Golang (for running the Go-Control-Plane)
-   An Envoy Proxy setup

### Installation
    
1.  Build and run the containers (including Envoy, Redis, and the Go-Control-Plane):
    
    ```bash
    docker-compose up --build
    ```
    

### Configuration

-   **Envoy Configuration**: Configured to communicate with the Go-Control-Plane for dynamic updates.
-   **Go-Control-Plane**: Implements the xDS protocol to manage Envoy's configuration, particularly for updating Redis cluster information and health checks.

### Health Check Implementation

The health checks for Redis are implemented using TCP checks. The Go-Control-Plane sends `PING` commands to the Redis nodes and expects a `PONG` response to confirm their status.

### Usage

Once the setup is running, Envoy will dynamically adjust its routing based on the health and status of the Redis nodes as provided by the Go-Control-Plane.
