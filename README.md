# binpak
A new way of visualising Kubernetes workloads. 

## Building

Build the container with:

```shell
docker build --tag binpak:latest .
```

## Running

Run the container locally using your kubeconfig with:

```shell
gcloud container clusters get-credentials ${CLUSTER_NAME} --zone ${ZONE}
docker run \
  --mount type=bind,source=$HOME/.kube/config,target=/.kube/config \
  --env KUBECONFIG=/.kube/config \
  --expose 8080 -p 8080:8080/tcp \
  binpak:latest
```

Then test with:

```shell
curl localhost:8080/info
```

## UI

Build the UI container with:

```shell
docker build --tag binpak-ui:latest ./ui/.
```

Run this with:

```shell
docker run --expose 3000 -p 3000:3000/tcp binpak-ui:latest
```

## Deploy

```shell
gcloud auth configure-docker europe-west1-docker.pkg.dev
./deploy.sh
```
