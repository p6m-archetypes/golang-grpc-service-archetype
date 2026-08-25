FROM golang:1.23 AS builder
RUN apt-get update && apt-get install -y --no-install-recommends protobuf-compiler && rm -rf /var/lib/apt/lists/*
# protoc plugin versions pinned to the protobuf/grpc versions in go.mod
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
WORKDIR /app
# A fresh render carries no go.sum (the dependency set is decided at render time), so the
# lockfile is produced inside the build: `go.*` copies whichever of go.mod/go.sum exist, and
# `go mod tidy` completes go.sum before the compile. Run `make setup` and commit go.sum to pin.
COPY go.* ./
RUN go mod download
COPY . .
# Generate the gRPC bindings so a clean render serves its API out of the box
RUN mkdir -p gen && protoc -I proto \
    --go_out=gen --go_opt=paths=source_relative \
    --go-grpc_out=gen --go-grpc_opt=paths=source_relative \
    proto/{{ project_name }}/{{ project_name }}.proto
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o server ./cmd/server

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/server /server
EXPOSE {{ service_port }}
EXPOSE {{ management_port }}
ENTRYPOINT ["/server"]
