FROM golang:1.25@sha256:a22b2e6c5e753345b9759fba9e5c1731ebe28af506745e98f406cc85d50c828e AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o rcm-node-installer ./cmd/node-installer
RUN /app/rcm-node-installer -h

# Using busybox instead of scratch so that the nsenter utility is present, as used in restarter logic
FROM busybox:1.37@sha256:d80cd694d3e9467884fcb94b8ca1e20437d8a501096cdf367a5a1918a34fc2fd
COPY --from=builder /app/rcm-node-installer /rcm-node-installer

ENTRYPOINT ["/rcm-node-installer"]
