package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/slaskis/ace/internal/test"
)

func TestAce(t *testing.T) {
	t.Run("single recipient", func(t *testing.T) {
		os.Remove("testdata/.env1.ace")
		{
			cmd := &Set{EnvFile: "testdata/.env1.ace", Recipients: []string{"testdata/recipients1.txt"}, EnvPairs: []string{"A=1", "B=2", "C=1 2 3 "}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}

		{
			cmd := &Set{EnvFile: "testdata/.env1.ace", Recipients: []string{"testdata/recipients1.txt"}, Input: strings.NewReader("X=1\nY=2\nZ=3")}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
		{
			buf := &bytes.Buffer{}
			cmd := &Get{EnvFile: "testdata/.env1.ace", Identity: "testdata/identity1", Output: buf}
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
			cmd := &Set{EnvFile: "testdata/.env2.ace", Recipients: []string{"testdata/recipients12.txt"}, EnvPairs: []string{"A=1", "B=2", "C=1 2 3 "}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
		{
			cmd := &Set{EnvFile: "testdata/.env2.ace", Recipients: []string{"testdata/recipients1.txt"}, EnvPairs: []string{"A=2", "D=3"}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
		{
			cmd := &Set{EnvFile: "testdata/.env2.ace", Recipients: []string{"testdata/recipients13.txt"}, EnvPairs: []string{"E=5"}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}
		t.Run("identity1", func(t *testing.T) {
			buf := &bytes.Buffer{}
			cmd := &Get{EnvFile: "testdata/.env2.ace", Identity: "testdata/identity1", Output: buf}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})
		t.Run("identity2", func(t *testing.T) {
			buf := &bytes.Buffer{}
			cmd := &Get{EnvFile: "testdata/.env2.ace", Identity: "testdata/identity2", Output: buf}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})
		t.Run("identity3", func(t *testing.T) {
			buf := &bytes.Buffer{}
			cmd := &Get{EnvFile: "testdata/.env2.ace", Identity: "testdata/identity3", Output: buf}
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
			cmd := &Set{EnvFile: "testdata/.env3.ace", Recipients: []string{"testdata/recipients1.txt"}, EnvPairs: []string{"A=1", "B=2", "C=1 2 3 "}}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Run("identity1", func(t *testing.T) {
			buf := &bytes.Buffer{}
			cmd := &Env{EnvFile: "testdata/.env3.ace", Identity: "testdata/identity1", Command: []string{"sh", "-c", "echo $A $B $C"}, Output: buf}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})
		t.Run("identity2", func(t *testing.T) {
			buf := &bytes.Buffer{}
			cmd := &Env{EnvFile: "testdata/.env3.ace", Identity: "testdata/identity2", Command: []string{"sh", "-c", "echo $A $B $C"}, Output: buf}
			err := cmd.Run()
			if err != nil {
				t.Fatal(err)
			}
			test.Snapshot(t, buf.Bytes())
		})
	})
}
