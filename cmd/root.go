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
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/hiifong/gh-tea/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	pkgErr "github.com/hiifong/gh-tea/pkg/errors"
)

var (
	cfgFile string
	cfg     config.Config
	use     string
	debug   bool

	err error

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "gh tea",
		Short: "A brief description of your application",
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: rootRun,
	}
)

func init() {
	cobra.OnInitialize(initLog)
	cobra.OnInitialize(initConfig)

	defaultCfgFile, err := os.UserHomeDir()
	pkgErr.Check(err)
	defaultCfgFile = path.Join(defaultCfgFile, ".config/gh-tea/config.yaml")

	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config",
		"c",
		"",
		fmt.Sprintf("config file (default is %s)", defaultCfgFile),
	)

	rootCmd.PersistentFlags().StringVar(
		&use,
		"use",
		"",
		"choose a gitea",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&debug,
		"debug",
		"d",
		false,
		"enable debug mode",
	)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		pkgErr.Check(err)
		home = path.Join(home, ".config/gh-tea")
		if _, err = os.Stat(home); os.IsNotExist(err) {
			err = os.MkdirAll(home, os.ModePerm)
			pkgErr.Check(err)
		}
		// Search config in home directory with name ".gh-tea" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	viper.SetEnvPrefix("TEA")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// If a config file is found, read it in.
	if err = viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) && cfgFile == "" {
			home, err := os.UserHomeDir()
			pkgErr.Check(err)
			home = path.Join(home, ".config/gh-tea")
			cfgFile = path.Join(home, "config.yaml")
			if _, err = os.Stat(cfgFile); os.IsNotExist(err) {
				_, err = os.OpenFile(cfgFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
				pkgErr.Check(err)
				viper.SetDefault("tea", config.Tea{
					config.Default: config.TeaItem{
						User: map[config.Username]config.UserItem{},
					},
				})
				log.Debugf("using config file: %s", cfgFile)
				err = viper.WriteConfigAs(cfgFile)
				pkgErr.Check(err)
				err = viper.Unmarshal(&cfg)
				pkgErr.Check(err)
			}
		}
	} else {
		log.Debugf("using config file: %s", viper.ConfigFileUsed())
		err = viper.Unmarshal(&cfg)
		pkgErr.Check(err)
	}
	log.Debugf("loaded config: %v", cfg)
}

func initLog() {
	log.SetReportCaller(true)
	log.SetPrefix("[Tea]")
	if debug {
		log.SetLevel(log.DebugLevel)
		log.SetPrefix("[Tea debug]")
	} else {
		log.SetOutput(os.Stderr)
		log.SetLevel(log.InfoLevel)
	}
}

func rootRun(cmd *cobra.Command, args []string) {
	log.Info("hi world, this is the gh-tea extension!")
	client, err := api.DefaultRESTClient()
	if err != nil {
		log.Fatal(err)
	}
	response := struct{ Login string }{}
	err = client.Get("user", &response)
	pkgErr.Check(err)
	log.Infof("running as %s", response.Login)

	var repos []struct {
		Name        string
		Private     bool
		Description string
		CloneURL    string `json:"clone_url"`
		Topics      []string
		Visibility  string
	}
	err = client.Get(fmt.Sprintf("users/%s/repos?sort=updated&per_page=10", response.Login), &repos)
	pkgErr.Check(err)
	log.Infof("repos: \n%+v", repos)
}
