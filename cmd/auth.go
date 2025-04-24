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
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"runtime"

	"code.gitea.io/sdk/gitea"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/hiifong/gh-tea/pkg/config"
)

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(refreshCmd)
}

var (
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
		Run:   logoutRun,
	}

	refreshCmd = &cobra.Command{
		Use:   "refresh",
		Short: "Expand or fix the permission scopes for stored credentials for active account.",
		Run:   refreshRun,
	}
)

func loginRun(cmd *cobra.Command, args []string) {
	var ocfg oauth2.Config
	var v config.TeaItem
	var isDefault, ok bool

	ocfg.ClientID = config.GiteaOAuth2.ClientID

	if use != "" {
		if v, ok = cfg.Tea[config.TeaName(use)]; ok {
			ocfg.Endpoint.AuthURL = fmt.Sprintf(config.GiteaOAuth2.Endpoint.AuthURL, v.Host)
			ocfg.Endpoint.TokenURL = fmt.Sprintf(config.GiteaOAuth2.Endpoint.TokenURL, v.Host)
		}
		if !ok {
			v = cfg.Tea[config.Default]
			if v.Name == "" || v.Host == "" {
				log.Fatal("no name specified")
			}
			isDefault = true
			ocfg.Endpoint.AuthURL = fmt.Sprintf(config.GiteaOAuth2.Endpoint.AuthURL, v.Host)
			ocfg.Endpoint.TokenURL = fmt.Sprintf(config.GiteaOAuth2.Endpoint.TokenURL, v.Host)
		}
	} else {
		v = cfg.Tea[config.Default]
		if v.Name == "" {
			log.Fatal("no name specified")
		}
		v = cfg.Tea[config.TeaName(v.Name)]
		if v.Host == "" {
			log.Fatal("no host specified")
		}
		isDefault = true
		ocfg.Endpoint.AuthURL = fmt.Sprintf(config.GiteaOAuth2.Endpoint.AuthURL, v.Host)
		ocfg.Endpoint.TokenURL = fmt.Sprintf(config.GiteaOAuth2.Endpoint.TokenURL, v.Host)
	}
	token, err := getToken(cmd.Context(), ocfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("token: %+v", token)
	client, err := gitea.NewClient(v.Host, gitea.SetToken(token.AccessToken))
	if err != nil {
		log.Fatal(err)
	}
	info, resp, err := client.GetMyUserInfo()
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatal(err)
	}
	log.Infof("info: %+v", info)
	v.User = make(map[config.Username]config.UserItem)
	v.User[config.Username(info.UserName)] = token
	if isDefault {
		cfg.Tea[config.TeaName(v.Name)] = v
	} else {
		cfg.Tea[config.TeaName(use)] = v
	}
	config.WriteConfig(cfg.Tea)
	log.Infof("Successfully logged in: %s", info.UserName)
}

func logoutRun(cmd *cobra.Command, args []string) {
	// TODO
}

func refreshRun(cmd *cobra.Command, args []string) {
	// TODO
}

var template string = `<!DOCTYPE html>
<html lang="en">
<head>
	<title>gh-tea</title>
</head>
<body>
<p>Success. You may close this page and return to gh or gh-tea.</p>
<p style="font-style: italic">&mdash;<a href="https://github.com/hiifong/gh-tea">gh-tea</a></p>
</body>
</html>`

func getToken(ctx context.Context, c oauth2.Config) (*oauth2.Token, error) {
	state := oauth2.GenerateVerifier()
	queries := make(chan url.Values)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: consider whether to show errors in browser or command line
		queries <- r.URL.Query()
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(template))
	})
	var server *httptest.Server
	if c.RedirectURL == "" {
		server = httptest.NewServer(handler)
		c.RedirectURL = server.URL
	} else {
		server = httptest.NewUnstartedServer(handler)
		url, err := url.Parse(c.RedirectURL)
		if err != nil {
			log.Fatal(err)
		}
		origHostname := url.Hostname()
		if url.Port() == "" {
			url.Host += ":0"
		}
		l, err := net.Listen("tcp", url.Host)
		if err != nil {
			log.Fatal(err)
		}
		server.Listener = l
		server.Start()
		url.Host = l.Addr().String()
		if url.Hostname() != origHostname {
			// restore original hostname such as 'localhost'
			url.Host = fmt.Sprintf("%s:%s", origHostname, url.Port())
		}
		c.RedirectURL = url.String()
	}
	defer server.Close()
	verifier := oauth2.GenerateVerifier()
	authCodeURL := c.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
	log.Infof("Please complete authentication in your browser...\n%s\n", authCodeURL)
	var open string
	var p []string
	switch runtime.GOOS {
	case "windows":
		open = "rundll32"
		p = append(p, "url.dll,FileProtocolHandler")
	case "darwin":
		open = "open"
	default:
		open = "xdg-open"
	}
	p = append(p, authCodeURL)
	// TODO: wait for server to start before opening browser
	if _, err = exec.LookPath(open); err == nil {
		err = exec.Command(open, p...).Run()
		if err != nil {
			log.Fatalf("Unable to open browser using '%s': %s\n", open, err)
		}
	}
	query := <-queries
	server.Close()
	if query.Get("state") != state {
		return nil, fmt.Errorf("state mismatch")
	}
	code := query.Get("code")
	return c.Exchange(ctx, code, oauth2.VerifierOption(verifier))
}
