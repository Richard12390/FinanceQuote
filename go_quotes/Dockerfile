# Build
FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -o /out/quote-publisher ./cmd/quotepub

# Runtime
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /out/quote-publisher .
COPY configs ./configs
ENV QUOTE_CONFIG=./configs/pipelines
ENTRYPOINT ["./quote-publisher","-config","./configs"]
