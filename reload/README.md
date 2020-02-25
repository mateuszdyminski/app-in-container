# Reload

All steps to demo how to reload configuration automatically in local Docker and K8s deployments.

## Demo - Local docker deployment

Build docker

```zsh
make docker-build
```

```zsh
docker run -p 8080:8080 -v $PWD/config:/config --rm --name reload -d mateuszdyminski/app-container-reload:latest
```

Go to: [http://localhost:8080/config](http://localhost:8080/config)

Change configuration in `config` directory.

Reload configuration via http request:

```zsh
curl -X POST http://localhost:8080/-/reload
```

Check new configuration: [http://localhost:8080/config](http://localhost:8080/config)

Change again configuration in `config` directory.

Send signal to docker to reload configuration

```zsh
docker kill --signal=SIGHUP $(docker ps -aqf "name=reload")
```

## Demo - Kubernetes deployment

Run minikube:

```zsh
minikube start
```

Deploy configuration:

```zsh
kubectl apply -f k8s/config.yaml
```

Deploy application:

```zsh
kubectl apply -f k8s/deployment.yaml
```

Go to: [http://192.168.64.2:32090/config](http://192.168.64.2:32090/config) to verify if everything works

Change configuration in file `k8s/config.yaml`.

Re-deploy configuration:

```zsh
kubectl apply -f k8s/config.yaml
```

Go to: [http://192.168.64.2:32090/config](http://192.168.64.2:32090/config) and verify new configuration

To clean up:

```zsh
kubectl delete -f k8s
```
