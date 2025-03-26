#!/bin/bash

# Create the volume for the database file
fly volumes create micromdm_data --size 1

# Deploy the app
fly deploy

# Scale to one instance
fly scale count 1

# Create SSL certificate
fly certs create micromdm-app.fly.dev

# Show the app status
fly status 