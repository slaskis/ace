# ACE

*A*ppend-only en*C*rypted *E*nvironment variables.

- Append-only allows the file to be updated without needing to be able to decrypt any previous environment variables.
- Encrypted variables, while keys are public so it's easy to see any changes

## API

- `ace set .env.enc [KEY=VALUE]`
- `ace set .env.enc < .env`
- `ace get .env.enc`
- `ace env .env.enc -- command`

Ex.

```
# ace/v1:<encryption header for the written block, including recipients, key information>
MY_KEY=[Base58 encoded string encrypted as defined by block header]

# ace/v1:<encryption header for the written block, including recipients, key information>
# a normal comment which is not encrypted
# this block will override the previously set MY_KEY
MY_KEY=[Base58 encoded string encrypted as defined by block header]

# ace/v1:<encryption header for the written block, including recipients, key information>
# this block will unset the MY_KEY env var
MY_KEY=
```
