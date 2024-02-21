docker-compose down 
# docker volume prune -a -f
yes | docker volume rm envoy_victoria_data
yes | docker volume rm envoy_redis-node-1-data
yes | docker volume rm envoy_redis-node-2-data
yes | docker volume rm envoy_redis-node-3-data
yes | docker volume rm envoy_grafana_data
docker-compose up -d --build
# docker exec -it  envoy-redis-node-1-1 redis-cli --cluster create  redis-node-1:6379 redis-node-2:6379 redis-node-3:6379 --cluster-replicas 0
