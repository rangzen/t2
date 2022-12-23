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

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// usageCmd represents the usage command
var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Display usage of the translation service",
	Long: `Display usage and limit of the selected translation service
if the service provide such informations.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := printUsage(); err != nil {
			log.Fatal(err)
		}
	},
}

func printUsage() error {
	ts, err := selectBackend()
	if err != nil {
		return err
	}

	u, err := ts.Usage()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Usage: %d/%d\n", u.Used, u.Limit)
	return nil
}

func init() {
	rootCmd.AddCommand(usageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
