FROM golang:1.12-alpine as builder

RUN apk update && apk add sqlite build-base git mercurial
WORKDIR /build

COPY vendor vendor
COPY cmd cmd
COPY pkg pkg
COPY Makefile Makefile
COPY go.mod go.mod
RUN make static

FROM golang:1.10-alpine as probe-builder

RUN apk update && apk add build-base git
ENV ORG github.com/grpc-ecosystem
ENV PROJECT $ORG/grpc_health_probe
WORKDIR /go/src/$PROJECT

COPY --from=builder /build/vendor/$ORG/grpc-health-probe .
COPY --from=builder /build/vendor .
RUN CGO_ENABLED=0 go install -a -tags netgo -ldflags "-w"

FROM scratch
COPY --from=builder /build/bin/bin/initializer /bin/initializer
COPY --from=builder /build/bin/registry-server /bin/registry-server
COPY --from=builder /build/bin/configmap-server /bin/configmap-server
COPY --from=builder /build/bin/appregistry-server /bin/appregistry-server
COPY --from=probe-builder /go/bin/grpc_health_probe /bin/grpc_health_probe
EXPOSE 50051
ENTRYPOINT ["/bin/registry-server"]
