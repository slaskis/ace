package main

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"encoding/base32"
	"fmt"
	"io"
	"os"
	"strings"

	"filippo.io/age"
	"github.com/tdewolff/argp"
	"golang.org/x/crypto/chacha20poly1305"
)

type Main struct{}

func (cmd *Main) Run() error {
	return argp.ShowUsage
}

const ACE_PREFIX = "# ace/v1:"

func readEnvFile(src io.Reader, identities []age.Identity) ([]string, error) {
	var keys []string
	vals := map[string]string{}

	s := bufio.NewScanner(src)
	var aead cipher.AEAD
	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		// split on ACE_PREFIX
		if strings.HasPrefix(line, ACE_PREFIX) {
			// base32decode and armor decode age header
			header, err := base32.StdEncoding.DecodeString(strings.TrimPrefix(line, ACE_PREFIX))
			if err != nil {
				return nil, err
			}

			var r io.Reader
			r = bytes.NewReader(header)

			// decrypt the block key using identities
			r, err = age.Decrypt(r, identities...)
			if err != nil && err.Error() == "no identity matched any of the recipients" {
				// try next env block
				aead = nil
				continue
			} else if err != nil {
				return nil, err
			}
			blockKey, err := io.ReadAll(r)
			if err != nil {
				return nil, err
			}
			aead, err = chacha20poly1305.NewX(blockKey)
			if err != nil {
				return nil, err
			}
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		// decrypt each secret using block key
		pair := strings.SplitN(line, "=", 2)
		if len(pair) != 2 {
			continue
		}

		if aead == nil {
			continue
		}

		secret, err := base32.StdEncoding.DecodeString(pair[1])
		if err != nil {
			return nil, err
		}

		if len(secret) < aead.NonceSize() {
			return nil, fmt.Errorf("ciphertext too short")
		}
		nonce, ciphertext := secret[:aead.NonceSize()], secret[aead.NonceSize():]

		// Decrypt the message and check it wasn't tampered with.
		plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return nil, err
		}

		if _, exists := vals[pair[0]]; !exists {
			keys = append(keys, pair[0])
		}
		vals[pair[0]] = string(plaintext)
	}

	var newVars []string
	for _, k := range keys {
		newVars = append(newVars, k+"="+vals[k])
	}

	return newVars, nil
}

// configurable for tests
var input io.Reader = os.Stdin
var output io.Writer = os.Stdout

// this is set using `-ldflags "-X main.version=1.2.3"`
var version string

func main() {
	var r, f, i []string
	cmd := argp.NewCmd(&Main{}, "ace")
	cmd.AddCmd(&Set{Recipients: argp.Append{I: &r}, RecipientFiles: argp.Append{I: &f}}, "set", "Append encrypted env vars to file")
	cmd.AddCmd(&Get{Identities: argp.Append{I: &i}}, "get", "Decrypt env with available identities")
	cmd.AddCmd(&Env{Identities: argp.Append{I: &i}}, "env", "Expand to env and pass to command")
	cmd.AddCmd(&Version{version: version}, "version", "Command version")
	cmd.Parse()
}
