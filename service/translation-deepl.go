package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const deeplEndpointUsage = "https://api-free.deepl.com/v2/usage"

type TranslationDeepl struct {
	Endpoint string
	ApiKey   string
}

type DeeplRequestResponse struct {
	Translations []DeeplRequestResponseTranslation
}

type DeeplRequestResponseTranslation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

type DeeplRequestUsage struct {
	CharacterCount int64 `json:"character_count"`
	CharacterLimit int64 `json:"character_limit"`
}

func (d TranslationDeepl) Translate(text string, source string, target string) (TranslationResponse, error) {
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
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		return TranslationResponse{}, errors.New(
			fmt.Sprint("status:", res.StatusCode, " body:", string(body)),
		)
	}

	var dres DeeplRequestResponse
	err = json.Unmarshal(body, &dres)
	if err != nil {
		log.Fatal(err)
	}

	sb := strings.Builder{}
	for _, t := range dres.Translations {
		sb.WriteString(t.Text)
	}
	return TranslationResponse{Text: sb.String()}, nil
}

// prepareDeeplConfig creates the DeepL configuration
func (d TranslationDeepl) prepareDeeplConfig(text string, source string, target string) url.Values {
	deeplConfig := url.Values{}
	deeplConfig.Set("text", text)
	checkedSource := checkSource(source)
	deeplConfig.Set("source_lang", checkedSource)
	deeplConfig.Set("target_lang", target)
	return deeplConfig
}

// checkSource will correct if needed the source language
func checkSource(source string) string {
	// DeepL accept EN-GB and EN-US in target language but not as source language.
	// https://www.deepl.com/docs-api/translating-text/request/
	if source == "EN-GB" || source == "EN-US" {
		return "EN"
	}
	return source
}

// prepareRequest creates the HTTP Request
func (d TranslationDeepl) prepareRequest(deeplConfig url.Values) (*http.Request, error) {
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

func (d TranslationDeepl) Usage() (UsageResponse, error) {
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
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var dres DeeplRequestUsage
	err = json.Unmarshal(body, &dres)
	if err != nil {
		log.Fatal(err)
	}
	return UsageResponse{
		Used:  dres.CharacterCount,
		Limit: dres.CharacterLimit,
	}, nil
}
