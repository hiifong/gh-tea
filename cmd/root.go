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
	"os"

	"github.com/charmbracelet/log"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	debug   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh tea",
	Short: "A brief description of your application",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(_ *cobra.Command, args []string) {
		log.Info("hi world, this is the gh-tea extension!")
		client, err := api.DefaultRESTClient()
		if err != nil {
			log.Fatal(err)
		}
		response := struct{ Login string }{}
		err = client.Get("user", &response)
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("running as %s", response.Login)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLog)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file (default is $HOME/.gh-tea.yaml)",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&debug,
		"debug",
		"d",
		false,
		"enable debug mode",
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gh-tea" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gh-tea")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Using config file:", viper.ConfigFileUsed())
	}
}

func initLog() {
	log.SetPrefix("[Tea]")
	if debug {
		log.SetLevel(log.DebugLevel)
		log.SetPrefix("[Tea debug]")
	} else {
		log.SetOutput(os.Stderr)
		log.SetLevel(log.InfoLevel)
	}
}
