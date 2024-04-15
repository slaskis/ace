# Docker example

This examples how to use ace with docker volumes to mount secrets and identity and expand them using `ace env` as an entrypoint.

It also shows that it works with a scratch image, which may be useful in statically compiled projects.

## How to run

```sh
docker build -t slaskis/ace-example:scratch .
docker run \
	--mount type=bind,source="$(pwd)"/.env.ace,target=/run/secrets/env,readonly \
	--mount type=bind,source="$(pwd)"/identity,target=/run/secrets/identity,readonly \
	slaskis/ace-example:scratch
```
