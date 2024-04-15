package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/slaskis/ace/internal/test"
	"github.com/tdewolff/argp"
)

func TestAce(t *testing.T) {
	t.Run("single recipient", func(t *testing.T) {
		os.Remove("testdata/.env1.ace")
		{
			cmd := &Set{EnvFile: "testdata/.env1.ace", Recipients: argp.Append{I: &([]string{"testdata/recipients1.txt"})}, EnvPairs: []string{"A=1", "B=2", "C=1 2 3 "}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}

		{
			input = strings.NewReader("X=1\nY=2\nZ=3")
			cmd := &Set{EnvFile: "testdata/.env1.ace", Recipients: argp.Append{I: &([]string{"testdata/recipients1.txt"})}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
		{
			buf := &bytes.Buffer{}
			output = buf
			cmd := &Get{EnvFile: "testdata/.env1.ace", Identities: argp.Append{I: &([]string{"testdata/identity1"})}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		}
	})

	t.Run("multiple recipients", func(t *testing.T) {
		os.Remove("testdata/.env2.ace")
		{
			cmd := &Set{EnvFile: "testdata/.env2.ace", Recipients: argp.Append{I: &([]string{"testdata/recipients12.txt"})}, EnvPairs: []string{"A=1", "B=2", "C=1 2 3 "}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
		{
			cmd := &Set{EnvFile: "testdata/.env2.ace", Recipients: argp.Append{I: &([]string{"testdata/recipients1.txt"})}, EnvPairs: []string{"A=2", "D=3"}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
		{
			cmd := &Set{EnvFile: "testdata/.env2.ace", Recipients: argp.Append{I: &([]string{"testdata/recipients13.txt"})}, EnvPairs: []string{"E=5"}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
		t.Run("identity1", func(t *testing.T) {
			buf := &bytes.Buffer{}
			output = buf
			cmd := &Get{EnvFile: "testdata/.env2.ace", Identities: argp.Append{I: &([]string{"testdata/identity1"})}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})
		t.Run("identity2", func(t *testing.T) {
			buf := &bytes.Buffer{}
			output = buf
			cmd := &Get{EnvFile: "testdata/.env2.ace", Identities: argp.Append{I: &([]string{"testdata/identity2"})}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})
		t.Run("identity3", func(t *testing.T) {
			buf := &bytes.Buffer{}
			output = buf
			cmd := &Get{EnvFile: "testdata/.env2.ace", Identities: argp.Append{I: &([]string{"testdata/identity3"})}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})
	})

	t.Run("env", func(t *testing.T) {
		os.Remove("testdata/.env3.ace")
		{
			cmd := &Set{EnvFile: "testdata/.env3.ace", Recipients: argp.Append{I: &([]string{"testdata/recipients1.txt"})}, EnvPairs: []string{"A=1", "B=2", "C=1 2 3 "}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Run("identity1", func(t *testing.T) {
			buf := &bytes.Buffer{}
			output = buf
			cmd := &Env{EnvFile: "testdata/.env3.ace", Identities: argp.Append{I: &([]string{"testdata/identity1"})}, Command: []string{"sh", "-c", "echo $A $B $C"}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})
		t.Run("identity2", func(t *testing.T) {
			buf := &bytes.Buffer{}
			output = buf
			cmd := &Env{EnvFile: "testdata/.env3.ace", Identities: argp.Append{I: &([]string{"testdata/identity2"})}, Command: []string{"sh", "-c", "echo $A $B $C"}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})
	})

	t.Run("multiple recipients repeated flags", func(t *testing.T) {
		os.Remove("testdata/.env4.ace")
		{
			cmd := &Set{EnvFile: "testdata/.env4.ace", Recipients: argp.Append{I: &([]string{"testdata/recipients1.txt", "testdata/recipients2.txt"})}, EnvPairs: []string{"A=1", "B=2", "C=1 2 3 "}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
		{
			cmd := &Set{EnvFile: "testdata/.env4.ace", Recipients: argp.Append{I: &([]string{"testdata/recipients1.txt"})}, EnvPairs: []string{"A=2", "D=3"}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
		{
			cmd := &Set{EnvFile: "testdata/.env4.ace", Recipients: argp.Append{I: &([]string{"testdata/recipients2.txt"})}, EnvPairs: []string{"C=333 "}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
		t.Run("identity1", func(t *testing.T) {
			buf := &bytes.Buffer{}
			output = buf
			cmd := &Get{EnvFile: "testdata/.env4.ace", Identities: argp.Append{I: &([]string{"testdata/identity1"})}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})
		t.Run("identity2", func(t *testing.T) {
			buf := &bytes.Buffer{}
			output = buf
			cmd := &Get{EnvFile: "testdata/.env4.ace", Identities: argp.Append{I: &([]string{"testdata/identity2"})}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})
		t.Run("identity1,identity2", func(t *testing.T) {
			buf := &bytes.Buffer{}
			output = buf
			cmd := &Get{EnvFile: "testdata/.env4.ace", Identities: argp.Append{I: &([]string{"testdata/identity1", "testdata/identity2"})}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})
		t.Run("identity2,identity1", func(t *testing.T) {
			buf := &bytes.Buffer{}
			output = buf
			cmd := &Get{EnvFile: "testdata/.env4.ace", Identities: argp.Append{I: &([]string{"testdata/identity2", "testdata/identity1"})}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})
	})
}
