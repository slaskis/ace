package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

type Env struct {
	OnMissing  string   `arg:"--on-missing" default:"error" help:"How to handle when env-file or identity is missing, can be 'ignore', 'warn' or 'error'"`
	EnvFile    string   `arg:"--env-file,-e" default:"./.env.ace"`
	Identities []string `arg:"--identity,-i,separate" help:"Decrypt using the specified IDENTITY. Can be repeated. Defaults to $XDG_CONFIG_HOME/ace/identity"`
	Command    []string `arg:"positional,required"`
}

func (cmd *Env) Run() error {
	if len(cmd.Command) == 0 {
		return fmt.Errorf("missing command to run")
	}
	src, err := os.Open(cmd.EnvFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			switch cmd.OnMissing {
			case "ignore":
				// silence
			case "warn", "warning":
				slog.Warn("env-file not found")
			default:
				return err
			}
		} else {
			return err
		}
	} else {
		defer src.Close()
	}

	identities, err := readIdentities(cmd.Identities, cmd.OnMissing)
	if err != nil {
		return err
	}

	vars, err := readEnvFile(src, identities, false)
	if err != nil {
		if err.Error() == "no identities specified" {
			switch cmd.OnMissing {
			case "ignore":
				// silence
			case "warn", "warning":
				slog.Warn(err.Error())
			default:
				return err
			}
		} else {
			return err
		}
	}

	// run command with vars added
	c := exec.Command(cmd.Command[0], cmd.Command[1:]...)
	c.Env = append(os.Environ(), vars...)
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Stdout = output
	return c.Run()
}
