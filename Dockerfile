FROM alpine:latest
COPY ./.secret/* /root/.secret/
COPY focus /usr/local/
EXPOSE 7001 7002
ENTRYPOINT ["/usr/local/focus", "--env", "prod"]