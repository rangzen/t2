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

package deepl

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

const deeplEndpointUsage = "https://api-free.deepl.com/v2/usage"

type TranslationService struct {
	Endpoint string
	ApiKey   string
}

type RequestResponse struct {
	Translations []RequestResponseTranslation
}

type RequestResponseTranslation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

type RequestUsage struct {
	CharacterCount int64 `json:"character_count"`
	CharacterLimit int64 `json:"character_limit"`
}

func (d TranslationService) Name() string {
	return "DeepL"
}

func (d TranslationService) Translate(text string, source string, target string) (backend.TranslationResponse, error) {
	deeplConfig := d.prepareDeeplConfig(text, source, target)

	req, err := d.prepareRequest(deeplConfig)
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
	for _, t := range dres.Translations {
		sb.WriteString(t.Text)
	}
	return backend.TranslationResponse{Text: sb.String()}, nil
}

// prepareDeeplConfig creates the DeepL configuration
func (d TranslationService) prepareDeeplConfig(text string, source string, target string) url.Values {
	deeplConfig := url.Values{}
	deeplConfig.Set("text", text)
	checkedSource := checkDeeplSource(source)
	deeplConfig.Set("source_lang", checkedSource)
	deeplConfig.Set("target_lang", target)
	return deeplConfig
}

// checkDeeplSource will correct if needed the source language
func checkDeeplSource(source string) string {
	// DeepL accept EN-GB and EN-US in target language but not as source language.
	// https://www.deepl.com/docs-api/translating-text/request/
	if source == "EN-GB" || source == "EN-US" {
		return "EN"
	}
	return source
}

// prepareRequest creates the HTTP Request
func (d TranslationService) prepareRequest(deeplConfig url.Values) (*http.Request, error) {
	dcEncoded := deeplConfig.Encode()
	req, err := http.NewRequest(http.MethodPost, d.Endpoint, strings.NewReader(dcEncoded))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "DeepL-Auth-Key "+d.ApiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(dcEncoded)))
	return req, err
}

func (d TranslationService) Usage() (backend.UsageResponse, error) {
	req, err := http.NewRequest(http.MethodGet, deeplEndpointUsage, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "DeepL-Auth-Key "+d.ApiKey)

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

	var dres RequestUsage
	err = json.Unmarshal(body, &dres)
	if err != nil {
		log.Fatal(err)
	}
	return backend.UsageResponse{
		Used:  dres.CharacterCount,
		Limit: dres.CharacterLimit,
	}, nil
}
