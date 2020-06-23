rm pod_watcher
GOOS=linux GOARCH=amd64 go build -o pod_watcher
GOOS=linux GOARCH=amd64 go build
docker build -t chaocai/pod-watcher:v0.0.1 .
docker push chaocai/pod-watcher:v0.0.1