FROM alpine
RUN apk add -U --no-cache ca-certificates

ADD pod_watcher /pod_watcher/pod_watcher
WORKDIR /pod_watcher
ENTRYPOINT ./pod_watcher
