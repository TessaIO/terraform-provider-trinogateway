#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


VERSION=17
BASE_URL="https://repo1.maven.org/maven2/io/trino/gateway/gateway-ha"
JAR_FILE="gateway-ha-$VERSION-jar-with-dependencies.jar"
GATEWAY_JAR="gateway-ha.jar"
CONFIG_YAML="config.yaml"

# Copy necessary files
copy_files() {
    if [[ ! -f "$GATEWAY_JAR" ]]; then
        echo "Fetching $GATEWAY_JAR version $VERSION"
        curl -O "$BASE_URL/$VERSION/$JAR_FILE"
        mv "$JAR_FILE" "$GATEWAY_JAR"
    fi
}

# Start PostgreSQL database if not running
start_postgres_db() {
    if ! docker ps --format '{{.Names}}' | grep -q '^local-postgres$'; then
        echo "Starting PostgreSQL database container"
        PGPASSWORD=mysecretpassword
        docker run -v "$PWD/$POSTGRES_SQL:/tmp/$POSTGRES_SQL" \
            --name local-postgres -p 5432:5432 -e POSTGRES_PASSWORD=$PGPASSWORD -d postgres
        sleep 5
        docker exec local-postgres psql -U postgres -h localhost -c 'CREATE DATABASE gateway'
    fi
}

# Main execution flow
copy_files
start_postgres_db

# Start Trino Gateway server
echo "Starting Trino Gateway server..."
java --version
java -Xmx1g -jar ./$GATEWAY_JAR ./$CONFIG_YAML
