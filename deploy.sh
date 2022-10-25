#!/bin/bash

set -euo pipefail

# To run: run ./depoy.sh from the root of binpak repo
GOOGLE_PROJECT=[name-of-your-project]
BACKEND_IMAGE=europe-west1-docker.pkg.dev/${GOOGLE_PROJECT}/binpak/binpak
UI_IMAGE=europe-west1-docker.pkg.dev/${GOOGLE_PROJECT}/binpak/binpak-ui
COMBINED_IMAGE=europe-west1-docker.pkg.dev/${GOOGLE_PROJECT}/binpak/binpak-combined

# TODO: Proper versioning
VERSION=latest

# Build and push backend
echo "Building ${BACKEND_IMAGE}:${VERSION}"
docker build -t ${BACKEND_IMAGE}:${VERSION} . --platform linux/amd64
docker push ${BACKEND_IMAGE}:${VERSION}
#
## Build and push UI
echo "Building ${UI_IMAGE}:${VERSION}"
docker build -t ${UI_IMAGE}:${VERSION} ./ui/. --platform linux/amd64
docker push ${UI_IMAGE}:${VERSION}

## Build and push Combined image for mock/demo purposes
echo "Building ${COMBINED_IMAGE}:${VERSION}"
docker build -t europe-west1-docker.pkg.dev/${GOOGLE_PROJECT}/binpak/binpak-combined:latest -f ui/demo-image.Dockerfile ui/ --platform linux/amd64
docker push ${COMBINED_IMAGE}:${VERSION}

# Set kubernetes context
#gcloud container clusters get-credentials binpak --zone europe-west2-c --project ${GOOGLE_PROJECT}
#
## Deploy to kubernetes
#BACKEND_IMAGE_TAG=${BACKEND_IMAGE}:${VERSION} envsubst < manifests/deployment-backend.tpl > \
#  manifests/deployment-backend.yaml
#UI_IMAGE_TAG=${UI_IMAGE}:${VERSION} envsubst < manifests/deployment-ui.tpl > \
#  manifests/deployment-ui.yaml
#
## kubectl create deployment binpak --image=${IMAGE}:${VERSION} --port=8080 --dry-run=client -o yaml > manifests/deployment-backend.tpl
#kubectl apply -f manifests
