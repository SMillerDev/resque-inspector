# builder
FROM golang:1.24-bookworm AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -trimpath -o resque-inspector resque-inspector

# Create a production stage to run the application binary
FROM scratch AS production

WORKDIR /app
COPY --from=builder /build/resque-inspector ./

EXPOSE 5678
CMD ["/app/resque-inspector", "serve"]