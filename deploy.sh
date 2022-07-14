#!/bin/bash

# To run: run ./depoy.sh from the root of binpak repo
IMAGE=europe-west1-docker.pkg.dev/jetstack-wil/binpak/binpak
VERSION=0.0.4

# # Build and push
docker build -t ${IMAGE}:${VERSION} . --platform linux/amd64
docker push ${IMAGE}:${VERSION}

# Set kubernetes context
gcloud container clusters get-credentials binpak --zone europe-west2-c --project jetstack-wil

# Deploy to kubernetes
# kubectl create deployment binpak --image=${IMAGE}:${VERSION} --port=8080 --dry-run=client -o yaml > manifests/binpak-deployment.yaml
kubectl apply -f manifests -n binpak
