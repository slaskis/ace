FROM scratch
ADD --chmod=0755 https://github.com/slaskis/ace/releases/download/v0.1.13/ace-linux-arm64 /bin/ace
COPY ./env /bin/aceenv
ENTRYPOINT ["/bin/ace","env","-e=/run/secrets/env","-i=/run/secrets/identity","--"]
CMD ["/bin/aceenv"]
