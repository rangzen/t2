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

package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rangzen/t2/pkg/backend"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type TranslationService struct {
	Endpoint string
	ApiKey   string
}

type RequestResponse struct {
	Data TranslateTextResponseList `json:"data"`
}

type TranslateTextResponseList struct {
	Translations []TranslateTextResponseTranslation `json:"translations"`
}

type TranslateTextResponseTranslation struct {
	DetectedSourceLanguage string `json:"detectedSourceLanguage"`
	Model                  string `json:"model"`
	Text                   string `json:"translatedText"`
}

type RequestUsage struct {
	CharacterCount int64 `json:"character_count"`
	CharacterLimit int64 `json:"character_limit"`
}

func (d TranslationService) Name() string {
	return "Google"
}

func (d TranslationService) Translate(text string, source string, target string) (backend.TranslationResponse, error) {
	googleConfig := d.prepareGoogleConfig(text, source, target)

	req, err := d.prepareRequest(googleConfig)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		return backend.TranslationResponse{}, errors.New(
			fmt.Sprint("status:", res.StatusCode, " body:", string(body)),
		)
	}

	var dres RequestResponse
	err = json.Unmarshal(body, &dres)
	if err != nil {
		log.Fatal(err)
	}

	sb := strings.Builder{}
	for _, t := range dres.Data.Translations {
		sb.WriteString(t.Text)
	}
	return backend.TranslationResponse{Text: sb.String()}, nil
}

// prepareGoogleConfig creates the DeepL configuration
func (d TranslationService) prepareGoogleConfig(text string, source string, target string) url.Values {
	deeplConfig := url.Values{}
	deeplConfig.Set("q", text)
	checkedTarget := checkGoogleLanguage(target)
	deeplConfig.Set("target", checkedTarget)
	deeplConfig.Set("format", "text")
	checkedSource := checkGoogleLanguage(source)
	deeplConfig.Set("source", checkedSource)
	deeplConfig.Set("key", d.ApiKey)
	return deeplConfig
}

// checkGoogleLanguage will correct if needed the language
func checkGoogleLanguage(source string) string {
	// https://cloud.google.com/translate/docs/languages
	if source == "EN-GB" || source == "EN-US" {
		return "EN"
	}
	source = strings.ToLower(source)
	return source
}

// prepareRequest creates the HTTP Request
func (d TranslationService) prepareRequest(config url.Values) (*http.Request, error) {
	dcEncoded := config.Encode()
	req, err := http.NewRequest(http.MethodPost, d.Endpoint, strings.NewReader(dcEncoded))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(dcEncoded)))
	return req, err
}

func (d TranslationService) Usage() (backend.UsageResponse, error) {
	return backend.UsageResponse{}, errors.New("Check Google Cloud Console for usages.")
}
