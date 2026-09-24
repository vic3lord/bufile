# syntax=docker/dockerfile:1

ARG GO_VERSION=1.27
ARG KUBECTL_VERSION=v1.31.0

FROM golang:${GO_VERSION} AS build
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /out/bufile .

FROM build AS test
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go test -v ./...

FROM curlimages/curl:latest AS kubectl
ARG TARGETOS
ARG TARGETARCH
ARG KUBECTL_VERSION
RUN curl -fsSL -o /tmp/kubectl \
      "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/${TARGETOS}/${TARGETARCH}/kubectl" \
    && chmod 0750 /tmp/kubectl

FROM gcr.io/distroless/static AS final
COPY --from=kubectl /tmp/kubectl /usr/bin/kubectl
COPY --from=build /out/bufile /bufile
ENTRYPOINT ["/bufile"]
