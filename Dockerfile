FROM gcr.io/distroless/static

ARG BINARY
COPY $BINARY /usr/bin/app

ENTRYPOINT ["/usr/bin/app"]
