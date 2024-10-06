package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/tdewolff/argp"
)

type Env struct {
	OnMissing  string      `name:"on-missing" default:"error" desc:"How to handle when env-file or identity is missing, can be 'ignore', 'warn' or 'error'"`
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

	identities, err := readIdentities(*cmd.Identities.I.(*[]string), cmd.OnMissing)
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
