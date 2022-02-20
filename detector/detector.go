package detector

import (
	"bufio"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-github/v41/github"
	"github.com/hashicorp/go-multierror"
)

const NoName = "[No name]"
const NoDesc = "[No description]"
const NoFormula = "[Could not find a formula or a cask]"

var brewRe = regexp.MustCompile(`^brew\s+["']([^"',]+)["']`)
var caskRe = regexp.MustCompile(`^cask\s+["']([^"',]+)["']`)
var nameRe = regexp.MustCompile(`^\s+name\s+["'](.+)["']`)
var descRe = regexp.MustCompile(`^\s+desc\s+["'](.+)["']`)

type Detector struct {
	client *github.Client
}

func New(client *github.Client) *Detector {
	return &Detector{
		client: client,
	}
}

type Formula struct {
	Name string
	Desc string
}

func (d *Detector) Detect(ctx context.Context, line string) (*Formula, error) {
	// brew
	if strings.HasPrefix(line, "brew") {
		m := brewRe.FindStringSubmatch(line)
		if len(m) < 1 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}
		splitted := strings.Split(m[1], "/")
		switch {
		case len(splitted) == 1:
			// homebrew-core
			return d.parseFormula(ctx, "Homebrew", "homebrew-core", fmt.Sprintf("Formula/%s.rb", splitted[0]))
		case len(splitted) == 3:
			// tap repositories
			f, err := d.parseFormula(ctx, splitted[0], fmt.Sprintf("homebrew-%s", splitted[1]), fmt.Sprintf("Formula/%s.rb", splitted[2]))
			if err == nil {
				return f, nil
			}
			f, err2 := d.parseFormula(ctx, splitted[0], fmt.Sprintf("homebrew-%s", splitted[1]), fmt.Sprintf("%s.rb", splitted[2]))
			if err2 == nil {
				return f, nil
			}
			var merr error
			merr = multierror.Append(merr, err)
			merr = multierror.Append(merr, err2)
			return nil, merr
		default:
			return nil, fmt.Errorf("invalid line: %s", line)
		}
	}
	// cask
	if strings.HasPrefix(line, "cask") {
		m := caskRe.FindStringSubmatch(line)
		if len(m) < 1 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}
		splitted := strings.Split(m[1], "/")
		switch {
		case len(splitted) == 1:
			// homebrew-cask
			f, err := d.parseFormula(ctx, "Homebrew", "homebrew-cask", fmt.Sprintf("Casks/%s.rb", splitted[0]))
			if err == nil {
				return f, nil
			}
			// homebrew-cask-fonts
			f, err2 := d.parseFormula(ctx, "Homebrew", "homebrew-cask-fonts", fmt.Sprintf("Casks/%s.rb", splitted[0]))
			if err2 == nil {
				return f, nil
			}
			// homebrew-cask-drivers
			f, err3 := d.parseFormula(ctx, "Homebrew", "homebrew-cask-drivers", fmt.Sprintf("Casks/%s.rb", splitted[0]))
			if err3 == nil {
				return f, nil
			}
			var merr error
			merr = multierror.Append(merr, err)
			merr = multierror.Append(merr, err2)
			merr = multierror.Append(merr, err3)
			return nil, merr
		default:
			return nil, fmt.Errorf("invalid line: %s", line)
		}
	}

	return nil, nil
}

func (d *Detector) parseFormula(ctx context.Context, owner, repo, path string) (*Formula, error) {
	f, _, _, err := d.client.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, fmt.Errorf("%s is not file", path)
	}
	c, err := f.GetContent()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(strings.NewReader(c))
	formula := &Formula{Name: NoName, Desc: NoDesc}
	for scanner.Scan() {
		line := scanner.Text()
		{
			m := descRe.FindStringSubmatch(line)
			if len(m) == 2 {
				formula.Desc = m[1]
			}
		}
		{
			m := nameRe.FindStringSubmatch(line)
			if len(m) == 2 {
				formula.Name = m[1]
			}
		}
	}
	return formula, nil
}
