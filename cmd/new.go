/*
Copyright © 2025 hiifong <f@ilo.nz>

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
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/hiifong/gh-tea/config"
)

func init() {
	rootCmd.AddCommand(newCmd)
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "New a gitea host named <name> for the gitea at <Host>",
	Long:  `usage: gh tea new <name> <Host>, this command like git remote add <name> <Host> command`,
	Run:   newRun,
}

func newRun(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.Fatal(`incorrect parameters, the new command only requires name and url parameters,
Please use gh tea new <name> <Host> instead`)
	}
	name := args[0]
	url := args[1]
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		log.Fatal("Host must start with http:// or https://")
	}
	log.Debugf("new %s %s", name, url)
	cfg.Tea[config.TeaName(name)] = config.TeaItem{
		Name: name,
		Host: url,
	}
	log.Infof("tea: %+v", cfg)
	config.WriteConfig(cfg.Tea)
}
