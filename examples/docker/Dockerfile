FROM alpine

ADD https://github.com/slaskis/ace/releases/download/v0.1.8/ace-linux-arm64 /bin/ace
RUN chmod +x /bin/ace

ENTRYPOINT ["ace","env","-e=/run/secrets/env","-i=/run/secrets/identity","--"]
CMD env
