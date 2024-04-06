package main

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io"
	"os"
	"os/exec"
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

type EnvVar struct {
	K, V string
}

type Get struct {
	EnvFile  string   `name:"env-file" short:"e" default:"./.env.ace"`
	Identity string   `name:"identity" short:"i" default:"$XDG_CONFIG_HOME/ace/identity"`
	Keys     []string `name:"keys" index:"*"`
}

func (cmd *Get) Run() error {
	src, err := os.Open(cmd.EnvFile)
	if err != nil {
		return err
	}
	defer src.Close()

	if _, exists := os.LookupEnv("XDG_CONFIG_HOME"); !exists {
		dir, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		os.Setenv("XDG_CONFIG_HOME", dir)
	}

	i, err := os.Open(os.ExpandEnv(cmd.Identity))
	if err != nil {
		return err
	}
	defer i.Close()

	identities, err := age.ParseIdentities(i)
	if err != nil {
		return err
	}

	vars, err := readEnvFile(src, identities)
	if err != nil {
		return err
	}

	for _, kv := range vars {
		if len(cmd.Keys) > 0 {
			var match bool
			for _, k := range cmd.Keys {
				if strings.HasPrefix(kv, k+"=") {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		fmt.Fprintln(os.Stdout, kv)
	}

	return nil
}

type Set struct {
	Recipients []string `name:"recipients" short:"r" default:"./recipients.txt"`
	EnvFile    string   `name:"env-file" short:"e" default:"./.env.ace"`
	EnvPairs   []string `name:"env" index:"*"`
}

func (cmd *Set) Run() error {
	var recipients []age.Recipient
	for _, r := range cmd.Recipients {
		rcp, err := os.Open(r)
		if err != nil {
			return err
		}
		defer rcp.Close()

		rec, err := age.ParseRecipients(rcp)
		if err != nil {
			return err
		}
		recipients = append(recipients, rec...)
	}

	blockKey := make([]byte, chacha20poly1305.KeySize)
	if _, err := rand.Read(blockKey); err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)

	// encrypt the key using age
	err := func() error {
		w, err := age.Encrypt(buf, recipients...)
		if err != nil {
			return err
		}
		defer w.Close()

		_, err = w.Write(blockKey)
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		return err
	}

	dst, err := os.OpenFile(cmd.EnvFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.WriteString(dst, ACE_PREFIX+base32.StdEncoding.EncodeToString(buf.Bytes())+"\n")
	if err != nil {
		return err
	}

	aead, err := chacha20poly1305.NewX(blockKey)
	if err != nil {
		return err
	}

	pairs := cmd.EnvPairs
	if len(pairs) == 0 {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			line := strings.TrimSpace(s.Text())
			if strings.HasPrefix(line, "#") {
				continue
			}
			if !strings.Contains(line, "=") {
				continue
			}
			pairs = append(pairs, line)
		}
	}

	for _, p := range pairs {
		pair := strings.SplitN(p, "=", 2)
		if len(pair) != 2 {
			continue
		}
		nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(pair[1])+aead.Overhead())
		if _, err := rand.Read(nonce); err != nil {
			return err
		}

		secret := base32.StdEncoding.EncodeToString(aead.Seal(nonce, nonce, []byte(pair[1]), nil))
		_, err = io.WriteString(dst, pair[0]+"="+secret+"\n")
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(dst, "\n")
	if err != nil {
		return err
	}

	return nil
}

type Env struct {
	EnvFile  string   `name:"env-file" short:"e" default:"./.env.ace"`
	Identity string   `name:"identity" short:"i" default:"$XDG_CONFIG_HOME/ace/identity"`
	Command  []string `index:"*"`
}

func (cmd *Env) Run() error {
	if len(cmd.Command) == 0 {
		return fmt.Errorf("missing command to run")
	}
	src, err := os.Open(cmd.EnvFile)
	if err != nil {
		return err
	}
	defer src.Close()

	if _, exists := os.LookupEnv("XDG_CONFIG_HOME"); !exists {
		dir, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		os.Setenv("XDG_CONFIG_HOME", dir)
	}

	i, err := os.Open(os.ExpandEnv(cmd.Identity))
	if err != nil {
		return err
	}
	defer i.Close()

	identities, err := age.ParseIdentities(i)
	if err != nil {
		return err
	}

	vars, err := readEnvFile(src, identities)
	if err != nil {
		return err
	}

	// run command with vars added
	c := exec.Command(cmd.Command[0], cmd.Command[1:]...)
	c.Env = append(os.Environ(), vars...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func readEnvFile(src io.Reader, identities []age.Identity) ([]string, error) {
	var vars []EnvVar

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
			if err != nil {
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

		vars = append(vars, EnvVar{
			K: pair[0],
			V: string(plaintext),
		})
	}

	mostRecentAt := map[string]int{}
	for i := len(vars) - 1; i >= 0; i-- {
		kv := vars[i]
		if _, ok := mostRecentAt[kv.K]; ok {
			// skip previous
			continue
		}
		mostRecentAt[kv.K] = i
	}

	var newVars []string
	for i, kv := range vars {
		idx := mostRecentAt[kv.K]
		if idx == i {
			newVars = append(newVars, kv.K+"="+kv.V)
		}
	}

	return newVars, nil
}

func main() {
	main := &Main{}
	argp := argp.NewCmd(main, "age")
	argp.AddCmd(&Set{}, "set", "Append encrypted env vars to file")
	argp.AddCmd(&Get{}, "get", "Decrypt env with available identities")
	argp.AddCmd(&Env{}, "env", "Expand to env and pass to command")
	argp.Parse()
}
