FROM gcr.io/distroless/static

ARG BINARY

COPY $BINARY /usr/bin/$BINARY

ENTRYPOINT ["/usr/bin/$BINARY"]
