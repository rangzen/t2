/*
Copyright © 2021 Cedric L'homme <public@l-homme.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"errors"
	"fmt"
	"github.com/rangzen/t2/service"
	"github.com/spf13/cobra"
	"log"
	"os"

	"github.com/spf13/viper"
)

var cfgFile string
var sourceLang string
var pivotLang string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "t2 [flags] \"Text to translate.\"",
	Example: "t2 --source EN-US \"I will treat my wound.\"",
	Short:   "Double translation",
	Long: `Use online translation services to translate from
a source language to a pivot language, then translate back
to the source language.
During this process, the most obvious errors are corrected.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := t2(args[0]); err != nil {
			log.Fatal(err)
		}
	},
}

func t2(t string) error {
	endpoint := viper.GetString("Endpoint")
	apiKey := viper.GetString("ApiKey")
	if endpoint == "" || apiKey == "" {
		return errors.New(".t2.yaml seems missing or empty")
	}

	d := service.TranslationDeepl{
		Endpoint: endpoint,
		ApiKey:   apiKey,
	}

	c := service.T2Config{
		SourceLang: sourceLang,
		PivotLang:  pivotLang,
	}

	t2 := service.NewT2(c, d)
	return t2.TraductionTranslation(t)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.t2.yaml)")

	rootCmd.Flags().StringVarP(&pivotLang, "pivot", "p", "FR", "pivot language")
	rootCmd.Flags().StringVarP(&sourceLang, "source", "s", "EN-US", "source language")
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

		// Search config in home directory with name ".double-deepl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".t2")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, errPrint := fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		if errPrint != nil {
			log.Fatal(errPrint)
		}
	}
}
