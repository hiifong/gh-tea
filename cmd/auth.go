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
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(refreshCmd)
}

var (
	// authCmd represents the auth command
	authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Authenticate gh-tea with Gitea",
	}

	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Authenticate with a Gitea host.",
		Run:   loginRun,
	}

	logoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "Remove authentication for a Gitea account.",
		Run:   loginRun,
	}

	refreshCmd = &cobra.Command{
		Use:   "refresh",
		Short: "Expand or fix the permission scopes for stored credentials for active account.",
		Run:   refreshRun,
	}
)

func loginRun(cmd *cobra.Command, args []string) {
	// TODO
}

func logoutRun(cmd *cobra.Command, args []string) {
	// TODO
}

func refreshRun(cmd *cobra.Command, args []string) {
	// TODO
}
