package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"io"
	"os"
	"strings"

	"filippo.io/age"
	"github.com/tdewolff/argp"
	"golang.org/x/crypto/chacha20poly1305"
)

type Set struct {
	Recipients argp.Append `name:"recipients" short:"r" desc:"Defaults to ./recipients.txt"`
	EnvFile    string      `name:"env-file" short:"e" default:"./.env.ace"`
	EnvPairs   []string    `name:"env" index:"*"`
}

func (cmd *Set) Run() error {
	recs := *cmd.Recipients.I.(*[]string)
	if len(recs) == 0 {
		recs = []string{"./recipients.txt"}
	}

	var recipients []age.Recipient
	for _, r := range recs {
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
		s := bufio.NewScanner(input)
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
