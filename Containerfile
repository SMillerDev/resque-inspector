# Create a production stage to run the application binary
FROM scratch

WORKDIR /app
COPY resque-inspector ./

EXPOSE 5678
CMD ["/app/resque-inspector", "serve"]