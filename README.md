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
