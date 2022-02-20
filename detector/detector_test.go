package detector

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/k1LoW/go-github-client/v41/factory"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		line string
		want *Formula
	}{
		{`brew "act"`, &Formula{Name: NoName, Desc: "Run your GitHub Actions locally 🚀"}},
		{`brew "mtr"`, &Formula{Name: NoName, Desc: "'traceroute' and 'ping' in a single tool"}},
		{`brew "k1LoW/tap/tbls"`, &Formula{Name: NoName, Desc: "tbls is a CI-Friendly tool for document a database, written in Go."}},
		{`brew "k1low/tap/tsocks", args: ["HEAD"]"`, &Formula{Name: NoName, Desc: "[No description]"}},
		{`brew "act" # hello hello hello`, &Formula{Name: NoName, Desc: "Run your GitHub Actions locally 🚀"}},
		{`cask "iterm2"`, &Formula{Name: "iTerm2", Desc: "Terminal emulator as alternative to Apple's Terminal app"}},
		{`cask "font-noto-sans-cjk-jp"`, &Formula{Name: "Noto Sans CJK JP", Desc: NoDesc}},
	}
	c, err := factory.NewGithubClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()
	d := New(c)
	for _, tt := range tests {
		got, err := d.Detect(ctx, tt.line)
		if err != nil {
			t.Error(err)
			continue
		}
		if diff := cmp.Diff(got, tt.want, nil); diff != "" {
			t.Errorf("%s", diff)
		}
	}
}
