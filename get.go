package main

import (
	"fmt"
	"os"
	"strings"

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

	identities, err := readIdentities(*cmd.Identities.I.(*[]string), "error")
	if err != nil {
		return err
	}

	vars, err := readEnvFile(src, identities, true)
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
