FROM scratch

# Set these as empty so resque-inspector uses the default
ENV REDIS_HOST=''
ENV REDIS_PORT=''
ENV REDIS_DSN=''

WORKDIR /app
COPY resque-inspector ./

EXPOSE 5678
CMD ["/app/resque-inspector", "serve"]