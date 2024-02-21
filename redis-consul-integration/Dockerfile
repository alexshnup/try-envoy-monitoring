FROM python:3.8-slim

WORKDIR /app

# Copy the Python script into the container
COPY update_consul.py .

# Install required Python packages
RUN pip install redis python-consul

# Command to run the script
CMD ["python", "update_consul.py"]
