# build stage
FROM golang:1.17 as builder

ENV GO111MODULE=on

ENV GOFLAGS="-mod=readonly"

WORKDIR /workspace/

# Copy the Go Modules manifests
COPY ./go.mod /workspace/
COPY ./go.sum /workspace/

RUN go mod download

COPY ./ /workspace/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on make binary

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/bin/main .
USER nonroot:nonroot

EXPOSE 9898

ENTRYPOINT ["/main"]