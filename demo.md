# Demo Steps

## Build `opm` binary

```sh
$ make build
GOFLAGS="-mod=vendor" go build  -o bin/appregistry-server ./cmd/appregistry-server
GOFLAGS="-mod=vendor" go build  -o bin/configmap-server ./cmd/configmap-server
GOFLAGS="-mod=vendor" go build  -o bin/initializer ./cmd/initializer
GOFLAGS="-mod=vendor" go build  -o bin/opm ./cmd/opm
GOFLAGS="-mod=vendor" go build  -o bin/registry-server ./cmd/registry-server
$ tree bin
bin
├── appregistry-server
├── configmap-server
├── initializer
├── opm
└── registry-server

0 directories, 5 files
$ chmod +x bin/opm
$ alias opm=$(pwd)/bin/opm
```

## Create initial index with bundle image

```sh
$ opm registry add --database bundles.db --bundle-images "quay.io/galletti94/prometheus-test@sha256:3ed785959f57fe70ccd895f8ed1aeeba769e45694de51a372b7bda8e4026495e"
...
$ opm registry serve --database bundles.db &
[1] 60909
WARN[0000] unable to set termination log path            error="open /dev/termination-log: operation not permitted"
INFO[0000] serving registry                              database=bundles.db port=50051
$ grpcurl -plaintext  localhost:50051 api.Registry/ListPackages
{
  "name": "prometheus"
}
```

## Update an existing registry database

```sh
$ opm registry add --database bundles.db --bundle-images "quay.io/olmtest/kiali:1.2.4"
...
$ opm registry serve --database bundles.db &
[1] 61417
WARN[0000] unable to set termination log path            error="open /dev/termination-log: operation not permitted"
INFO[0000] serving registry                              database=bundles.db port=50051
$ grpcurl -plaintext  localhost:50051 api.Registry/ListPackages
{
  "name": "kiali"
}
{
  "name": "prometheus"
}
```
