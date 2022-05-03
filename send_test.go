package kakscmd

import (
	"bytes"
	"encoding/hex"
	"os"
	"strings"
	"testing"
)

func TestWrite(t *testing.T) {
	cmd := "eval -client client0 echo hello world"
	for _, cmd := range []string{cmd, cmd + "\n"} {
		want := "022f000000260000006576616c202d636c69656e7420636c69656e7430206563686f2068656c6c6f20776f726c640a"

		b := &bytes.Buffer{}

		n, err := Write(b, cmd)
		if err != nil {
			t.Fatalf("should not error: %s", err.Error())
		}
		cmdNoNewline := strings.TrimSuffix(cmd, "\n")
		if n != len(cmdNoNewline)+10 {
			t.Errorf("unexpected amount of bytes: want %d got %d", len(cmdNoNewline), n)
		}

		res := hex.EncodeToString(b.Bytes())
		if len(res) != len(want) {
			t.Errorf("unexpected result: want / got \n%+v\n%+v\n", want, res)
		}

		for i := range res {
			if res[i] != want[i] {
				t.Fatalf("unexpected result: want / got \n%+v\n%+v\nfirst difference %d",
					want, res, i)
			}
		}
	}
}

func TestSocketPath(t *testing.T) {
	tests := []struct {
		name          string
		session       string
		want          string
		xdgRuntimeDir string
	}{
		{
			name:          "XDG",
			session:       "foo",
			want:          "/run/user/1000/kakoune/foo",
			xdgRuntimeDir: "/run/user/1000",
		},
		{
			name:          "tmp",
			session:       "foo",
			want:          "/tmp/kakoune/foo",
			xdgRuntimeDir: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("XDG_RUNTIME_DIR", tt.xdgRuntimeDir)
			if got := SocketPath(tt.session); got != tt.want {
				t.Errorf("SocketPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
