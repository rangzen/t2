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

package t2

import (
	"errors"
	"fmt"
	"github.com/rangzen/t2/pkg/backend"
	"github.com/rangzen/t2/pkg/backend/deepl"
	"github.com/rangzen/t2/pkg/backend/google"
)

// Backend is the interface that wraps the translation backend methods.
type Backend interface {
	Name() string
	Translate(text string, source string, pivot string) (backend.TranslationResponse, error)
	Usage() (backend.UsageResponse, error)
}

// Config is the configuration of the package.
type Config struct {
	SourceLang      string
	PivotLang       string
	DiffOnly        bool
	CopyToClipboard bool
}

// Diff is the interface that wraps the pretty print of the difference between
// the original text and the double translated text.
type Diff interface {
	Print(a, b string) string
}

// Clipboard is the interface that wraps the copy to clipboard functionality.
type Clipboard interface {
	Read() (string, error)
	Write(t string) error
}

// T2 is the main struct of the package.
type T2 struct {
	config    Config
	backend   Backend
	diff      Diff
	clipboard Clipboard
}

// NewT2 returns a new T2 struct.
func NewT2(config Config, backend Backend, diff Diff, clipboard Clipboard) T2 {
	return T2{
		config:    config,
		backend:   backend,
		diff:      diff,
		clipboard: clipboard,
	}
}

// Translate is the main function of the package.
// It translates the text from the source language to the pivot language,
// then back to the source language.
// It then prints the diff between the original text and the double translated text.
// If the copyToClipboard flag is set, it also copies the double translated text to the clipboard.
func (t T2) Translate(text string) error {
	if !t.config.DiffOnly {
		fmt.Println("# Original text")
		fmt.Println(text)
	}

	firstPass, err := t.backend.Translate(text, t.config.SourceLang, t.config.PivotLang)
	if err != nil {
		return err
	}
	if !t.config.DiffOnly {
		fmt.Printf("# Pivot text (%s -> %s by %s)\n", t.config.SourceLang, t.config.PivotLang, t.backend.Name())
		fmt.Println(firstPass.Text)
	}

	secondPass, err := t.backend.Translate(firstPass.Text, t.config.PivotLang, t.config.SourceLang)
	if err != nil {
		return err
	}
	if !t.config.DiffOnly {
		fmt.Printf("# Double translated text (%s -> %s by %s)\n", t.config.PivotLang, t.config.SourceLang, t.backend.Name())
		fmt.Println(secondPass.Text)
	}

	if !t.config.DiffOnly {
		fmt.Println("# Diff version")
	}
	prettyPrint := t.diff.Print(text, secondPass.Text)
	fmt.Println(prettyPrint)

	if t.config.CopyToClipboard {
		if err := t.clipboard.Write(secondPass.Text); err != nil {
			return err
		}
	}

	return nil
}

// SelectBackend returns the translation service implementation to use.
func SelectBackend(backend, endPoint, apiKey string) (Backend, error) {
	if endPoint == "" || apiKey == "" {
		return nil, errors.New("missing or incomplete configuration file (.t2.yaml)")
	}

	switch backend {
	case "deepl":
		return deepl.TranslationService{
			Endpoint: endPoint,
			ApiKey:   apiKey,
		}, nil
	case "google":
		return google.TranslationService{
			Endpoint: endPoint,
			ApiKey:   apiKey,
		}, nil
	default:
		return nil, errors.New("unknown translation service")
	}
}
