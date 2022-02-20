/*
Copyright © 2022 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/k1LoW/brewfile-desc/detector"
	"github.com/k1LoW/brewfile-desc/version"
	"github.com/k1LoW/go-github-client/v41/factory"
	"github.com/spf13/cobra"
)

var (
	inPlace bool
	force   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "brewfile-desc [BREWFILE]",
	Short:   "brewfile-desc add descriptions of formulae to Brewfile",
	Long:    `brewfile-desc add descriptions of formulae to Brewfile.`,
	Args:    cobra.ExactArgs(1),
	Version: version.Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := factory.NewGithubClient()
		if err != nil {
			return err
		}
		d := detector.New(client)
		fp, err := os.Open(args[0])
		if err != nil {
			return err
		}
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		if err := s.Color("cyan"); err != nil {
			return err
		}
		s.Suffix = " Analyzing Brewfile..."
		s.Start()

		data := [][]string{}
		maxlen := 0

		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			line := scanner.Text()
			formula, err := d.Detect(ctx, line)
			if err != nil {
				if force {
					_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
					data = append(data, []string{line, ""})
					continue
				}
				_ = fp.Close()
				return err
			}
			if formula == nil {
				data = append(data, []string{line, ""})
				continue
			}
			if maxlen < len(line) {
				maxlen = len(line)
			}
			if formula.Desc == detector.NoDesc && formula.Name != detector.NoName {
				formula.Desc = formula.Name
			}
			data = append(data, []string{line, fmt.Sprintf("# %s", formula.Desc)})
		}
		s.Stop()

		if err := fp.Close(); err != nil {
			return err
		}

		var out *os.File
		if inPlace {
			out, err = os.OpenFile(args[0], os.O_WRONLY|os.O_TRUNC, 0)
			if err != nil {
				return err
			}
		} else {
			out = os.Stdout
		}

		for _, d := range data {
			if d[1] == "" {
				if _, err := fmt.Fprintf(out, "%s\n", d[0]); err != nil {
					_ = out.Close()
					return err
				}
				continue
			}
			if _, err := fmt.Fprintf(out, fmt.Sprintf("%%-%ds %%s\n", maxlen), d[0], d[1]); err != nil {
				_ = out.Close()
				return err
			}
		}

		if err := out.Close(); err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&force, "force", "f", false, "force generate")
	rootCmd.Flags().BoolVarP(&inPlace, "in-place", "i", false, "edit Brewfile in-place")
}
