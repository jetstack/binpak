#!/bin/bash

# To run: run ./depoy.sh from the root of binpak repo

version=0.0.1

# Build and push
docker build -t europe-west1-docker.pkg.dev/jetstack-wil/binpak/binpak:${version} .
docker push europe-west1-docker.pkg.dev/jetstack-wil/binpak/binpak:${version}

# Set kubernetes context
gcloud container clusters get-credentials binpak --zone europe-west2-c --project jetstack-wil

# Deploy to kubernetes
# TODO