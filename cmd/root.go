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
	"io"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/hiifong/gh-tea/pkg/config"
	pkgErr "github.com/hiifong/gh-tea/pkg/errors"
	"github.com/hiifong/gh-tea/pkg/global"
	"github.com/hiifong/gh-tea/ui"
)

var (
	cfgFile string
	cfg     config.Config
	use     string
	debug   bool

	home string
	err  error

	rootCmd = &cobra.Command{
		Use:    "gh tea",
		Short:  "A brief description of your application",
		PreRun: preRun,
		Run:    rootRun,
	}
)

func init() {
	home, err = os.UserHomeDir()
	pkgErr.Check(err)
	home = path.Join(home, ".config/gh-tea")
	if _, err = os.Stat(home); os.IsNotExist(err) {
		err = os.MkdirAll(home, os.ModePerm)
		pkgErr.Check(err)
	}

	cobra.OnInitialize(initLog)
	cobra.OnInitialize(initConfig)

	defaultCfgFile := path.Join(home, "config.yaml")

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

func Execute() {
	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("TEA")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err = viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) && cfgFile == "" {
			cfgFile = path.Join(home, "config.yaml")
			if _, err = os.Stat(cfgFile); os.IsNotExist(err) {
				_, err = os.OpenFile(cfgFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
				pkgErr.Check(err)
				viper.SetDefault("tea", config.Tea{
					config.Default: config.TeaItem{
						Name: string(config.Default),
						Host: "https://gitea.com",
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
	file := &lumberjack.Logger{
		Filename:   path.Join(home, "gh-tea.log"),
		MaxSize:    100,
		MaxAge:     30,
		MaxBackups: 10,
		LocalTime:  true,
		Compress:   true,
	}
	global.Writer = io.MultiWriter(file, os.Stdout)
	level := log.InfoLevel
	if debug {
		level = log.DebugLevel
	}
	logger := log.NewWithOptions(global.Writer, log.Options{
		Prefix:       "[Tea]",
		ReportCaller: true,
		Level:        level,
	})
	log.SetDefault(logger)
}

func preRun(cmd *cobra.Command, args []string) {
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

func rootRun(cmd *cobra.Command, args []string) {
	err = ui.Run()
	if err != nil {
		log.Fatal(err)
	}
}
