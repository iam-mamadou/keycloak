#!/bin/bash

while [[ "$(curl --silent --output /dev/null --write-out ''%{http_code}'' http://localhost:8080)" != "200" ]]; do 
    echo "waiting for keycloak to be ready..."
    sleep 5
done

echo "Creating Keycloak admin user"
./opt/keycloak/bin/add-user-keycloak.sh -u admin -p admin -r master
./opt/keycloak/bin/jboss-cli.sh --connect command=:reload