package main

import (
	"fmt"
	"os"
	"strings"

	"filippo.io/age"
	"github.com/tdewolff/argp"
)

type Get struct {
	EnvFile    string      `name:"env-file" short:"e" default:"./.env.ace"`
	Identities argp.Append `name:"identity" short:"i" desc:"Defaults to $XDG_CONFIG_HOME/ace/identity"`
	Keys       []string    `name:"keys" index:"*"`
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

	idents := *cmd.Identities.I.(*[]string)
	if len(idents) == 0 {
		idents = []string{"$XDG_CONFIG_HOME/ace/identity"}
	}

	var identities []age.Identity
	for _, id := range idents {
		err := func() error {
			i, err := os.Open(os.ExpandEnv(id))
			if err != nil {
				return err
			}
			defer i.Close()

			idents, err := age.ParseIdentities(i)
			if err != nil {
				return err
			}
			identities = append(identities, idents...)
			return nil
		}()
		if err != nil {
			return err
		}
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
		fmt.Fprintln(output, kv)
	}

	return nil
}
