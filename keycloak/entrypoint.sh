#!/bin/bash

sh /setup_keycloak.sh &

exec /opt/keycloak/bin/standalone.sh -b 0.0.0.0
