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
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/hiifong/gh-tea/pkg/config"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove a host",
	Long:  `Usage: gh tes remove <host>, remove a <host> from the list`,
	Run:   removeRun,
}

func removeRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal(`incorrect parameters, Please use gh tea remove [name] instead`)
	}
	name := args[0]
	if _, ok := cfg.Tea[config.TeaName(name)]; !ok {
		log.Fatalf("%s not exist", name)
	}
	if v, ok := cfg.Tea[config.Default]; ok {
		if v.Name == name {
			log.Fatal("can't remove default name")
		}
	}
	delete(cfg.Tea, config.TeaName(name))
	config.WriteConfig(cfg.Tea)
}
