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

type TranslationGoogle struct {
	Endpoint string
	ApiKey   string
}

type GoogleRequestResponse struct {
	Data GoogleTranslateTextResponseList `json:"data"`
}

type GoogleTranslateTextResponseList struct {
	Translations []GoogleTranslateTextResponseTranslation `json:"translations"`
}

type GoogleTranslateTextResponseTranslation struct {
	DetectedSourceLanguage string `json:"detectedSourceLanguage"`
	Model                  string `json:"model"`
	Text                   string `json:"translatedText"`
}

type GoogleRequestUsage struct {
	CharacterCount int64 `json:"character_count"`
	CharacterLimit int64 `json:"character_limit"`
}

func (d TranslationGoogle) Translate(text string, source string, target string) (TranslationResponse, error) {
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
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		return TranslationResponse{}, errors.New(
			fmt.Sprint("status:", res.StatusCode, " body:", string(body)),
		)
	}

	var dres GoogleRequestResponse
	err = json.Unmarshal(body, &dres)
	if err != nil {
		log.Fatal(err)
	}

	sb := strings.Builder{}
	for _, t := range dres.Data.Translations {
		sb.WriteString(t.Text)
	}
	return TranslationResponse{Text: sb.String()}, nil
}

// prepareGoogleConfig creates the DeepL configuration
func (d TranslationGoogle) prepareGoogleConfig(text string, source string, target string) url.Values {
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
func (d TranslationGoogle) prepareRequest(config url.Values) (*http.Request, error) {
	dcEncoded := config.Encode()
	req, err := http.NewRequest(http.MethodPost, d.Endpoint, strings.NewReader(dcEncoded))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(dcEncoded)))
	return req, err
}

func (d TranslationGoogle) Usage() (UsageResponse, error) {
	return UsageResponse{}, errors.New("Check Google Cloud Console for usages.")
}
