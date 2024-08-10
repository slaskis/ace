package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

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

func TestIntegration(t *testing.T) {
	if os.Getenv("ACE_TESTBIN") == "" {
		t.Skip("Not running integration tests")
	}
	tests := []struct {
		Args  []string
		Stdin io.Reader
	}{
		{[]string{"ace"}, nil},
		{[]string{"ace", "version"}, nil},
		{[]string{"ace", "set", "-e=testdata/.env.invalid.ace", "A=1", "B=2"}, nil},
		{[]string{"ace", "get", "-e=testdata/.env.invalid.ace", "-i=testdata/nonexistent_identity.txt"}, nil},
		{[]string{"ace", "set", "-e=testdata/.env1.ace", "-r=invalid"}, nil},
		{[]string{"ace", "env", "-e=testdata/.env.invalid.ace", "-i=testdata/identity1", "--on-missing=warn", "--", "sh", "-c", "echo $A"}, nil},
		{[]string{"ace", "env", "-e=testdata/.env.invalid.ace", "-i=testdata/identity1", "--on-missing=ignore", "--", "sh", "-c", "echo $A"}, nil},
		{[]string{"ace", "env", "-e=testdata/.env.invalid.ace", "--", "sh", "-c", "echo $A"}, nil},

		{[]string{"rm", "-f", "testdata/.envi1.ace"}, nil},
		{[]string{"ace", "set", "-e=testdata/.envi1.ace", "-R=testdata/recipients1.txt"}, strings.NewReader("X=1\nY=2\nZ=3\n# comment\ninvalid line")},
		{[]string{"ace", "set", "-e=testdata/.envi1.ace", "-r=age10sunh5mqv3jw7audxcylw3s9redgjfhqenkuhm4v4hetg84q835qamk6x6"}, strings.NewReader("X=1\nY=2\nZ=3\n# comment\ninvalid line")},
		{[]string{"ace", "get", "-e=testdata/.envi1.ace", "-i=testdata/identity1"}, nil},
		{[]string{"ace", "env", "-e=testdata/.envi1.ace", "-i=testdata/identity1", "--", "sh", "-c", "echo $X"}, nil},
		{[]string{"ace", "env", "-e=testdata/.envi1.ace", "--on-missing=warn", "--", "sh", "-c", "echo $A"}, nil},

		{[]string{"rm", "-f", "testdata/.envi3.ace"}, nil},
		{[]string{"ace", "set", "-e=testdata/.envi3.ace", "-R=testdata/recipients1.txt", "A=1", "B=2", "C=1 2 3 "}, nil},
		{[]string{"ace", "get", "-e=testdata/.envi3.ace", "-i=testdata/identity1", "A"}, nil},

		{[]string{"rm", "-f", "testdata/.envi4.ace"}, nil},
		{[]string{"ace", "set", "-e=testdata/.envi4.ace", "-R=testdata/recipients1.txt", "-R=testdata/recipients2.txt", "A=1", "B=2", "C=1 2 3 "}, nil},
		{[]string{"ace", "set", "-e=testdata/.envi4.ace", "-R=testdata/recipients1.txt", "A=2", "D=3"}, nil},
		{[]string{"ace", "set", "-e=testdata/.envi4.ace", "-R=testdata/recipients2.txt", "C=333 "}, nil},
		{[]string{"ace", "get", "-e=testdata/.envi4.ace", "-i=testdata/identity1"}, nil},
		{[]string{"ace", "get", "-e=testdata/.envi4.ace", "-i=testdata/identity2"}, nil},
		{[]string{"ace", "get", "-e=testdata/.envi4.ace", "-i=testdata/identity1", "-i=testdata/identity2"}, nil},
		{[]string{"ace", "get", "-e=testdata/.envi4.ace", "-i=testdata/identity2", "-i=testdata/identity1"}, nil},
	}
	coverDir := ".coverdata/" + strconv.FormatInt(time.Now().Unix(), 10)
	os.MkdirAll(coverDir, 0755)
	for _, tt := range tests {
		t.Run(strings.ReplaceAll(strings.Join(tt.Args, " "), "/", "_"), func(t *testing.T) {
			if tt.Args[0] == "ace" {
				tt.Args[0] = os.Getenv("ACE_TESTBIN")
			}
			cmd := exec.Command(tt.Args[0], tt.Args[1:]...)
			cmd.Stdin = tt.Stdin
			cmd.Env = []string{
				"GOCOVERDIR=" + coverDir,
				"PATH=" + os.Getenv("PATH"),
				"HOME=/tmp",
			}
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Log(err)
			}
			test.Snapshot(t, out)
		})
	}
	t.Run("coverage", func(t *testing.T) {
		out, err := exec.Command("go", "tool", "covdata", "func", "-i="+coverDir).CombinedOutput()
		if err != nil {
			t.Log(err)
		}
		test.Snapshot(t, out)
	})
}
