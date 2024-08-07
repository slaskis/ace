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
	t.Run("set with missing default recipient file path", func(t *testing.T) {
		cmd := &Set{EnvFile: "testdata/.env.invalid.ace", Recipients: argp.Append{I: &([]string{})}, RecipientFiles: argp.Append{I: &([]string{})}, EnvPairs: []string{"A=1", "B=2"}}
		err := cmd.Run()
		if err == nil {
			t.Fatal("expected an error due to missing recipients file, but none occurred")
		}
	})
	t.Run("get with invalid identity file path", func(t *testing.T) {
		cmd := &Get{EnvFile: "testdata/.env.invalid.ace", Identities: argp.Append{I: &([]string{"testdata/nonexistent_identity.txt"})}}
		err := cmd.Run()
		if err == nil {
			t.Fatal("expected an error due to missing identity file, but none occurred")
		}
	})
	t.Run("single recipient", func(t *testing.T) {
		os.Remove("testdata/.env1.ace")
		{
			cmd := &Set{EnvFile: "testdata/.env1.ace", Recipients: argp.Append{I: &([]string{})}, RecipientFiles: argp.Append{I: &([]string{"testdata/recipients1.txt"})}, EnvPairs: []string{"A=1", "B=2", "C=1 2 3 ", "D ignored"}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}

		{
			input = strings.NewReader("X=1\nY=2\nZ=3\n# comment\ninvalid line")
			cmd := &Set{EnvFile: "testdata/.env1.ace", Recipients: argp.Append{I: &([]string{})}, RecipientFiles: argp.Append{I: &([]string{"testdata/recipients1.txt"})}}
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
			cmd := &Set{EnvFile: "testdata/.env2.ace", Recipients: argp.Append{I: &([]string{})}, RecipientFiles: argp.Append{I: &([]string{"testdata/recipients12.txt"})}, EnvPairs: []string{"A=1", "B=2", "C=1 2 3 "}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
		{
			cmd := &Set{EnvFile: "testdata/.env2.ace", Recipients: argp.Append{I: &([]string{})}, RecipientFiles: argp.Append{I: &([]string{"testdata/recipients1.txt"})}, EnvPairs: []string{"A=2", "D=3"}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
		{
			cmd := &Set{EnvFile: "testdata/.env2.ace", Recipients: argp.Append{I: &([]string{})}, RecipientFiles: argp.Append{I: &([]string{"testdata/recipients13.txt"})}, EnvPairs: []string{"E=5"}}
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
			cmd := &Set{EnvFile: "testdata/.env3.ace", Recipients: argp.Append{I: &([]string{})}, RecipientFiles: argp.Append{I: &([]string{"testdata/recipients1.txt"})}, EnvPairs: []string{"A=1", "B=2", "C=1 2 3 "}}
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
		t.Run("env-file on-missing=error", func(t *testing.T) {
			buf := &bytes.Buffer{}
			output = buf
			cmd := &Env{EnvFile: "testdata/.env.not-found.ace", Identities: argp.Append{I: &([]string{"testdata/identity2"})}, Command: []string{"sh", "-c", "echo $A $B $C"}}
			err := cmd.Run()
			if err == nil {
				t.Fatal("expected not such file or directory")
			}
			test.Snapshot(t, buf.Bytes())
		})

		t.Run("env-file on-missing=warn", func(t *testing.T) {
			buf := &bytes.Buffer{}
			output = buf
			t.Setenv("A", "woop")
			cmd := &Env{OnMissing: "warn", EnvFile: "testdata/.env.not-found.ace", Identities: argp.Append{I: &([]string{"testdata/identity2"})}, Command: []string{"sh", "-c", "echo $A $B $C"}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})

		t.Run("env-file on-missing=ignore", func(t *testing.T) {
			buf := &bytes.Buffer{}
			output = buf
			t.Setenv("A", "woop")
			cmd := &Env{OnMissing: "ignore", EnvFile: "testdata/.env.not-found.ace", Identities: argp.Append{I: &([]string{"testdata/identity2"})}, Command: []string{"sh", "-c", "echo $A $B $C"}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})

		t.Run("identity on-missing=error", func(t *testing.T) {
			buf := &bytes.Buffer{}
			output = buf
			cmd := &Env{EnvFile: "testdata/.env3.ace", Identities: argp.Append{I: &([]string{"testdata/identitynot-found"})}, Command: []string{"sh", "-c", "echo $A $B $C"}}
			err := cmd.Run()
			if err == nil {
				t.Fatal("expected not such file or directory")
			}
			test.Snapshot(t, buf.Bytes())
		})

		t.Run("identity on-missing=warn", func(t *testing.T) {
			buf := &bytes.Buffer{}
			output = buf
			t.Setenv("A", "woop")
			cmd := &Env{OnMissing: "warn", EnvFile: "testdata/.env3.ace", Identities: argp.Append{I: &([]string{"testdata/identitynot-found"})}, Command: []string{"sh", "-c", "echo $A $B $C"}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})

		t.Run("identity on-missing=ignore", func(t *testing.T) {
			buf := &bytes.Buffer{}
			output = buf
			t.Setenv("A", "woop")
			cmd := &Env{OnMissing: "ignore", EnvFile: "testdata/.env3.ace", Identities: argp.Append{I: &([]string{"testdata/identitynot-found"})}, Command: []string{"sh", "-c", "echo $A $B $C"}}
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
			cmd := &Set{EnvFile: "testdata/.env4.ace", Recipients: argp.Append{I: &([]string{})}, RecipientFiles: argp.Append{I: &([]string{"testdata/recipients1.txt", "testdata/recipients2.txt"})}, EnvPairs: []string{"A=1", "B=2", "C=1 2 3 "}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
		{
			cmd := &Set{EnvFile: "testdata/.env4.ace", Recipients: argp.Append{I: &([]string{})}, RecipientFiles: argp.Append{I: &([]string{"testdata/recipients1.txt"})}, EnvPairs: []string{"A=2", "D=3"}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
		{
			cmd := &Set{EnvFile: "testdata/.env4.ace", Recipients: argp.Append{I: &([]string{})}, RecipientFiles: argp.Append{I: &([]string{"testdata/recipients2.txt"})}, EnvPairs: []string{"C=333 "}}
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
