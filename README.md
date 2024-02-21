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
