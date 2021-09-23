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
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"log"
)

// clipboardCmd represents the clipboard command
var clipboardCmd = &cobra.Command{
	Use:   "clipboard",
	Short: "Use clipboard as input",
	Long: `Use clipboard as input.
Works on Windows, MacOS and Linux/Unix (require xsel or xclip).`,
	Run: func(cmd *cobra.Command, args []string) {
		t, err := clipboard.ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		if err := t2(t); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(clipboardCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clipboardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clipboardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
