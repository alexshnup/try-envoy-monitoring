import redis

# Connect to Redis
# Replace 'localhost' and '6379' with your Redis server's address and port if different
redis_client = redis.Redis(host='localhost', port=1999, db=0)

# Set a key-value pair
redis_client.set('mykey', 'myvalue')

# Get the value of the key
value = redis_client.get('mykey')

# Print the value retrieved
if value is not None:
    print(f"The value of 'mykey' is: {value.decode('utf-8')}")
else:
    print("Key not found")
