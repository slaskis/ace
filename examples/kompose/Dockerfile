FROM alpine
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ADD --chmod=0755 --link https://github.com/slaskis/ace/releases/download/v0.1.14/ace-$TARGETOS-$TARGETARCH /bin/ace
ENTRYPOINT ["/bin/ace","env","-e=/run/secrets/env","-i=/run/secrets/identity","--"]
CMD ["/usr/bin/env"]
