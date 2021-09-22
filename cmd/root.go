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
	"github.com/rangzen/double-deepl/service"
	"github.com/spf13/cobra"
	"log"
	"os"

	"github.com/spf13/viper"
)

var cfgFile string
var sourceLang string
var pivotalLang string
var usage bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "double-deepl [flags] \"Text to translate.\"",
	Example: "double-deepl -s EN-US \"I will treat my wound.\"",
	Short: "Double translation using deepl.com",
	Long: `Use deepl.com translation services to translate from
a source language to a pivotal one and translate again
to the source language.
In the process, most of the obvious errors are corrected.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		doubleDeepl(cmd, args)
	},
}

func doubleDeepl(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("need text to translate")
	}

	text := args[0]
	d := service.Deepl{
		Endpoint: viper.GetString("Endpoint"),
		ApiKey:   viper.GetString("ApiKey"),
	}

	fmt.Println("# Original text")
	fmt.Println(text)
	firstPass, err := d.DeeplTranslate(text, sourceLang, pivotalLang)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("# Pivotal text")
	fmt.Println(firstPass.Translations[0].Text)
	secondPass, err := d.DeeplTranslate(firstPass.Translations[0].Text, pivotalLang, sourceLang)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("# Double translated text")
	fmt.Println(secondPass.Translations[0].Text)

	if !usage {
		return nil
	}
	deeplUsage, err := d.DeeplUsage()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Usage: %d/%d\n", deeplUsage.CharacterCount, deeplUsage.CharacterLimit)
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.double-deepl.yaml)")

	rootCmd.Flags().StringVarP(&pivotalLang, "pivotal", "p", "FR", "pivotal language")
	rootCmd.Flags().StringVarP(&sourceLang, "source", "s", "EN", "source language")
	rootCmd.Flags().BoolVarP(&usage, "usage", "u", false, "display usage at the end")
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
		viper.SetConfigName(".double-deepl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
