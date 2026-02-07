#!/bin/bash

TRINO_IMAGE="trinodb/trino"
JAVA_OPTS="-Dhttp-server.process-forwarded=true"

# Start Trino servers
for i in 1 2; do
    PORT=808$i
    if ! lsof -i:$PORT > /dev/null; then
        docker run --name trino$i -d -p $PORT:8080 \
            -e JAVA_TOOL_OPTIONS="$JAVA_OPTS" $TRINO_IMAGE
    else
        echo "Warning: Port $PORT is already in use. Skipping trino$i."
    fi
done

# Add Trino servers as Gateway backends
add_backend() {
    curl -H "Content-Type: application/json" -X POST \
        localhost:8080/gateway/backend/modify/add \
        -d "{
              \"name\": \"$1\",
              \"proxyTo\": \"http://localhost:808$2\",
              \"active\": true,
              \"routingGroup\": \"adhoc\"
            }"
}

# Adding Trino servers as backends
for i in 1 2; do
    add_backend "trino$i" "$i"
done
