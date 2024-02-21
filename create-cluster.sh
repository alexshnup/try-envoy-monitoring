#!/bin/bash

echo "HOSTNAME=$HOSTNAME"

set -e

# Start Redis server in the background
redis-server --cluster-enabled yes --cluster-config-file /data/nodes.conf --cluster-node-timeout 5000 --appendonly yes &
# apt update; apt install curl inetutils-ping net-tools telnet vim bash -y
# netstat -atpn

# Function to check if Redis is ready
check_redis() {
    redis-cli -h $1 -p 6379 ping
}

# Wait for the current Redis node to be ready
until check_redis "127.0.0.1"; do
    echo "check_redis 127.0.0.1" 
    # netstat -atpn
    echo "Waiting for Redis at localhost to be ready..." 
    sleep 3
done

count=3
# Wait for all Redis nodes to be ready
for ip in $@; do
    if [ "$count" = 0 ]; then
        break
    fi
    count=$((count-1))
    echo "check_redis $ip"
    until check_redis $ip; do
        # netstat -atpn
        echo "Waiting for Redis at $ip to be ready..."
        sleep 3
    done
done

count=3
# Create a string of hosts with the port appended
hosts_with_ports=""
for ip in $@; do
    if [ "$count" = 0 ]; then
        break
    fi
    count=$((count-1))
    hosts_with_ports="$hosts_with_ports $ip:6379"
done

# Create the cluster
# Note: Assumes the script is run only on one node
if [ "$HOSTNAME" = "redis-node-1" ]; then
    echo "Creating cluster with nodes: $hosts_with_ports"
    yes 'yes' | redis-cli --cluster create --cluster-replicas 0 $hosts_with_ports
fi
