FROM alpine
ADD --chmod=0755 https://github.com/slaskis/ace/releases/download/v0.1.11/ace-linux-arm64 /bin/ace
ENTRYPOINT ["/bin/ace","env","-e=/run/secrets/env","-i=/run/secrets/identity","--"]
CMD ["/usr/bin/env"]
