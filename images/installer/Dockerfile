FROM golang:1.25@sha256:3b02b6795aa38e69f8e39b93a32acc02a423de1a8d508bbcb16b858cd6b4ac07 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o rcm-node-installer ./cmd/node-installer
RUN /app/rcm-node-installer -h

# Using busybox instead of scratch so that the nsenter utility is present, as used in restarter logic
FROM busybox:1.37@sha256:e3652a00a2fabd16ce889f0aa32c38eec347b997e73bd09e69c962ec7f8732ee
COPY --from=builder /app/rcm-node-installer /rcm-node-installer

ENTRYPOINT ["/rcm-node-installer"]
