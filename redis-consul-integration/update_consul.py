import redis
import consul
import time

# Configure Redis Sentinel
# sentinel = redis.StrictRedis(
#     host='localhost',
#     port=26379,
#     db=0,
#     decode_responses=True
# )
# sentinel = redis.sentinel.Sentinel([(sentinel, 26379)])

# Configure Redis Sentinel
sentinel = redis.sentinel.Sentinel([
    ('redis-sentinel1', 26379),
    ('redis-sentinel2', 26379),
    ('redis-sentinel3', 26379)
], decode_responses=True)

# Configure Consul Client
consul_client = consul.Consul(host='consul', port=8500)

def update_consul_master():
    try:
        # Get current master info from Redis Sentinel
        master = sentinel.discover_master('mymaster')
        master_host, master_port = master

        # Register the new master in Consul
        consul_client.agent.service.register(
            "redis-master",
            service_id="redis-master",
            address=master_host,
            port=master_port,
            tags=["master"],
            check=consul.Check.tcp(master_host, master_port, "10s")
        )
        print(f"Updated master in Consul: {master}")
    except Exception as e:
        print(f"Error updating Consul: {e}")

# Main loop
while True:
    update_consul_master()
    time.sleep(1)  # Sleep for 1 seconds