package main

import (
	"fmt"
	"os"
	"os/exec"

	"filippo.io/age"
)

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
