package run

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want *Run
	}{
		{
			name: "empty args",
			args: []string{},
			want: &Run{},
		},
		{
			name: "help flag",
			args: []string{"--help"},
			want: &Run{ShowHelp: true},
		},
		{
			name: "version flag",
			args: []string{"--version"},
			want: &Run{ShowVersion: true},
		},
		{
			name: "url only",
			args: []string{"https://example.com/script.sh"},
			want: &Run{URL: "https://example.com/script.sh"},
		},
		{
			name: "url with script args",
			args: []string{"https://example.com/script.sh", "arg1", "arg2"},
			want: &Run{
				URL:        "https://example.com/script.sh",
				ScriptArgs: []string{"arg1", "arg2"},
			},
		},
		{
			name: "all flags with url",
			args: []string{"+env", "+raw", "+nocache", "https://example.com/script.sh"},
			want: &Run{
				SendEnv: true,
				Raw:     true,
				NoCache: true,
				URL:     "https://example.com/script.sh",
			},
		},
		{
			name: "flags after url are script args",
			args: []string{"https://example.com/script.sh", "+env", "--help"},
			want: &Run{
				URL:        "https://example.com/script.sh",
				ScriptArgs: []string{"+env", "--help"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args)

			if got.ShowHelp != tt.want.ShowHelp {
				t.Errorf("ShowHelp = %v, want %v", got.ShowHelp, tt.want.ShowHelp)
			}
			if got.ShowVersion != tt.want.ShowVersion {
				t.Errorf("ShowVersion = %v, want %v", got.ShowVersion, tt.want.ShowVersion)
			}
			if got.SendEnv != tt.want.SendEnv {
				t.Errorf("SendEnv = %v, want %v", got.SendEnv, tt.want.SendEnv)
			}
			if got.Raw != tt.want.Raw {
				t.Errorf("Raw = %v, want %v", got.Raw, tt.want.Raw)
			}
			if got.NoCache != tt.want.NoCache {
				t.Errorf("NoCache = %v, want %v", got.NoCache, tt.want.NoCache)
			}
			if got.URL != tt.want.URL {
				t.Errorf("URL = %v, want %v", got.URL, tt.want.URL)
			}
			if len(got.ScriptArgs) != len(tt.want.ScriptArgs) {
				t.Errorf("ScriptArgs = %v, want %v", got.ScriptArgs, tt.want.ScriptArgs)
			}
		})
	}
}

func TestExec_Help(t *testing.T) {
	var stdout bytes.Buffer
	cmd := &Run{
		ShowHelp: true,
		Stdout:   &stdout,
		Stderr:   &bytes.Buffer{},
	}

	err := cmd.Exec(context.Background())
	if err != nil {
		t.Errorf("Exec() error = %v", err)
	}

	if !strings.Contains(stdout.String(), "usage:") {
		t.Errorf("Exec() stdout = %q, want usage info", stdout.String())
	}
}

func TestExec_Version(t *testing.T) {
	var stdout bytes.Buffer
	cmd := &Run{
		ShowVersion: true,
		Stdout:      &stdout,
		Stderr:      &bytes.Buffer{},
	}

	err := cmd.Exec(context.Background())
	if err != nil {
		t.Errorf("Exec() error = %v", err)
	}
}

func TestExec_NoURL(t *testing.T) {
	var stdout bytes.Buffer
	cmd := &Run{
		Stdout: &stdout,
		Stderr: &bytes.Buffer{},
	}

	err := cmd.Exec(context.Background())
	if err != nil {
		t.Errorf("Exec() error = %v", err)
	}

	// Should show help when no URL
	if !strings.Contains(stdout.String(), "usage:") {
		t.Errorf("Exec() stdout = %q, want usage info", stdout.String())
	}
}
