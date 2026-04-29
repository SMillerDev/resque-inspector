# Create a production stage to run the application binary
FROM scratch

WORKDIR /app
COPY resque-inspector ./

EXPOSE 5678
HEALTHCHECK CMD ["/app/resque-inspector", "health"]
CMD ["/app/resque-inspector", "serve"]