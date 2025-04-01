# builder
FROM ghcr.io/goreleaser/goreleaser:latest AS builder

WORKDIR /build

COPY . .

RUN goreleaser build --single-target --auto-snapshot --clean

# Create a production stage to run the application binary
FROM scratch AS production

WORKDIR /app
COPY --from=builder /build/dist/resque-inspector*/resque-inspector ./

EXPOSE 5678
CMD ["/app/resque-inspector", "serve"]