package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"filippo.io/age"
)

type Get struct {
	EnvFile  string   `name:"env-file" short:"e" default:"./.env.ace"`
	Identity string   `name:"identity" short:"i" default:"$XDG_CONFIG_HOME/ace/identity"`
	Keys     []string `name:"keys" index:"*"`
	Output   io.Writer
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
		fmt.Fprintln(cmd.Output, kv)
	}

	return nil
}
