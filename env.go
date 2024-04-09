package main

import (
	"fmt"
	"os"
	"os/exec"

	"filippo.io/age"
	"github.com/tdewolff/argp"
)

type Env struct {
	EnvFile    string      `name:"env-file" short:"e" default:"./.env.ace"`
	Identities argp.Append `name:"identity" short:"i" desc:"Defaults to $XDG_CONFIG_HOME/ace/identity"`
	Command    []string    `index:"*"`
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

	// run command with vars added
	c := exec.Command(cmd.Command[0], cmd.Command[1:]...)
	c.Env = append(os.Environ(), vars...)
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Stdout = output
	return c.Run()
}
