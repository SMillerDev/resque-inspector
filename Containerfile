FROM curlimages/curl:latest AS curl

# Create a production stage to run the application binary
FROM scratch

WORKDIR /app
COPY --from=curl /usr/bin/curl /usr/bin/curl
COPY resque-inspector ./

EXPOSE 5678
HEALTHCHECK --interval=30s --timeout=3s --retries=3 CMD ["/usr/bin/curl", "--fail", "--silent", "http://localhost:5678/health"]
CMD ["/app/resque-inspector", "serve"]