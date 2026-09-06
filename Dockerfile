# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/vanguard ./cmd/vanguard

FROM scratch

COPY --from=builder /out/vanguard /vanguard
USER 65532:65532
EXPOSE 8081
ENTRYPOINT ["/vanguard"]
