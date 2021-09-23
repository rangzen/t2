package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const deeplEndpointUsage = "https://api-free.deepl.com/v2/usage"

type Deepl struct {
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

func (d Deepl) DeeplTranslate(text string, source string, target string) (DeeplRequestResponse, error) {
	client := &http.Client{}
	deeplConfig := url.Values{}
	deeplConfig.Set("text", text)
	deeplConfig.Set("source_lang", source)
	deeplConfig.Set("target_lang", target)
	dcEncoded := deeplConfig.Encode()
	req, err := http.NewRequest(http.MethodPost, d.Endpoint, strings.NewReader(dcEncoded))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "DeepL-Auth-Key "+d.ApiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(dcEncoded)))

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		return DeeplRequestResponse{}, errors.New(
			fmt.Sprint("status:", res.StatusCode, " body:", string(body)),
		)
	}
	var dres DeeplRequestResponse
	err = json.Unmarshal(body, &dres)
	if err != nil {
		log.Fatal(err)
	}
	return dres, nil
}

func (d Deepl) DeeplUsage() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, deeplEndpointUsage, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "DeepL-Auth-Key "+d.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var dres DeeplRequestUsage
	err = json.Unmarshal(body, &dres)
	if err != nil {
		log.Fatal(err)
	}
	result := fmt.Sprintf("Usage: %d/%d\n", dres.CharacterCount, dres.CharacterLimit)
	return result, nil
}
