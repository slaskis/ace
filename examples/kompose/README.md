# Docker Compose example

This shows how ace can be used with docker compose secrets to mount secrets and identity and expand them using `ace env` as an entrypoint.

It also shows that it can be used with `kompose convert` to do the same thing in kubernetes.

## How to run

```sh
# to test locally
docker compose up

# to test in k8s
# (currently requires a version of kompose newer than v1.32.0)
kompose convert -o kompose.yaml
kubectl apply -f kompose.yaml
```
