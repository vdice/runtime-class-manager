FROM golang:1.25@sha256:cc737435e2742bd6da3b7d575623968683609a3d2e0695f9d85bee84071c08e6 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o rcm-node-installer ./cmd/node-installer
RUN /app/rcm-node-installer -h

# Using busybox instead of scratch so that the nsenter utility is present, as used in restarter logic
FROM busybox:1.37@sha256:b3255e7dfbcd10cb367af0d409747d511aeb66dfac98cf30e97e87e4207dd76f
COPY --from=builder /app/rcm-node-installer /rcm-node-installer

ENTRYPOINT ["/rcm-node-installer"]
