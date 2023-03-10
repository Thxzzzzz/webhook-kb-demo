FROM golang:1.19 as builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /workspace
COPY ./ /workspace
RUN go mod download
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o sidecar main.go


FROM alpine/curl:3.14
WORKDIR /
COPY --from=builder workspace/sidecar ./sidecar

ENTRYPOINT ["/sidecar"]