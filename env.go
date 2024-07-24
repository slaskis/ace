package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"filippo.io/age"
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
				log.Println("env-file not found")
			default:
				return err
			}
		} else {
			return err
		}
	} else {
		defer src.Close()
	}

	if _, exists := os.LookupEnv("XDG_CONFIG_HOME"); !exists {
		dir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("unable to read user config dir: %w", err)
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
				if errors.Is(err, os.ErrNotExist) {
					switch cmd.OnMissing {
					case "ignore":
						return nil
					case "warn", "warning":
						log.Println("identity not found")
						return nil
					default:
						return err
					}
				} else {
					return err
				}
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
		if err.Error() == "no identities specified" {
			switch cmd.OnMissing {
			case "ignore":
				// silence
			case "warn", "warning":
				log.Println(err.Error())
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
