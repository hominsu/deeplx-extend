package deeplx

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/abadojack/whatlanggo"
	"github.com/andybalholm/brotli"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
)

func initPayload(sourceLang, targetLang string) *PostData {
	hasRegionalVariant := false
	targetLangParts := strings.Split(targetLang, "-")

	// targetLang can be "en", "pt", "pt-PT", "pt-BR"
	// targetLangCode is the first part of the targetLang, e.g. "pt" in "pt-PT"
	targetLangCode := targetLangParts[0]
	if len(targetLangParts) > 1 {
		hasRegionalVariant = true
	}

	commonJobParams := CommonJobParams{
		WasSpoken:    false,
		TranscribeAS: "",
	}
	if hasRegionalVariant {
		commonJobParams.RegionalVariant = targetLang
	}

	return &PostData{
		Jsonrpc: "2.0",
		Method:  "LMT_handle_texts",
		Params: Params{
			Splitting: "newlines",
			Lang: Lang{
				SourceLangUserSelected: sourceLang,
				TargetLang:             targetLangCode,
			},
			CommonJobParams: commonJobParams,
		},
	}
}

func translateByOfficialAPI(text string, sourceLang string, targetLang string, authKey string, client *fasthttp.Client) (string, error) {
	freeURL := "https://api-free.deepl.com/v2/translate"
	textArray := strings.Split(text, "\n")

	payload := PayloadAPI{
		Text:       textArray,
		TargetLang: targetLang,
		SourceLang: sourceLang,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI(freeURL)
	req.SetBody(payloadBytes)

	req.Header.Set("Authorization", "DeepL-Auth-Key "+authKey)
	req.Header.Set("Content-Type", "application/json")

	err = client.Do(req, resp)
	if err != nil {
		return "", err
	}

	// Parsing the response
	var translationResponse TranslationResponse
	err = json.Unmarshal(resp.Body(), &translationResponse)
	if err != nil {
		return "", err
	}

	// Concatenating the translations
	var sb strings.Builder
	for _, translation := range translationResponse.Translations {
		sb.WriteString(translation.Text)
	}

	return sb.String(), nil
}

func translateByDeepLX(sourceLang, targetLang, translateText, authKey string, client *fasthttp.Client) (TranslationResult, error) {
	id := getRandomNumber()
	if sourceLang == "" {
		lang := whatlanggo.DetectLang(translateText)
		sourceLang = strings.ToUpper(lang.Iso6391())
	}

	if targetLang == "" {
		targetLang = "EN"
	}

	if translateText == "" {
		return TranslationResult{
			Code:    http.StatusNotFound,
			Message: "No translateText to translate",
		}, nil
	}

	www2 := "https://www2.deepl.com/jsonrpc"
	id = id + 1
	payload := initPayload(sourceLang, targetLang)
	text := Text{
		Text:                translateText,
		RequestAlternatives: 3,
	}
	payload.ID = id
	payload.Params.Texts = append(payload.Params.Texts, text)
	payload.Params.Timestamp = getTimeStamp(getICount(translateText))

	reqBody, _ := json.Marshal(payload)
	bodyStr := string(reqBody)

	// Adding spaces to the JSON string based on the ID to adhere to DeepL's request formatting rules
	if (id+5)%29 == 0 || (id+3)%13 == 0 {
		bodyStr = strings.Replace(bodyStr, "\"method\":\"", "\"method\" : \"", -1)
	} else {
		bodyStr = strings.Replace(bodyStr, "\"method\":\"", "\"method\": \"", -1)
	}

	// Creating a new HTTP POST request with the JSON data as the body
	reqBody = []byte(bodyStr)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI(www2)
	req.SetBody(reqBody)

	// Setting HTTP headers to mimic a request from the DeepL iOS App
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("x-app-os-name", "iOS")
	req.Header.Set("x-app-os-version", "16.3.0")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("x-app-device", "iPhone13,2")
	req.Header.Set("User-Agent", "DeepL-iOS/2.9.1 iOS 16.3.0 (iPhone13,2)")
	req.Header.Set("x-app-build", "510265")
	req.Header.Set("x-app-version", "2.9.1")
	req.Header.Set("Connection", "keep-alive")

	err := client.Do(req, resp)
	if err != nil {
		return TranslationResult{
			Code:    http.StatusServiceUnavailable,
			Message: "DeepL API request failed",
		}, nil
	}

	// Handling potential Brotli compressed response body
	var bodyReader io.Reader
	switch string(resp.Header.Peek("Content-Encoding")) {
	case "br":
		bodyReader = brotli.NewReader(resp.BodyStream())
	default:
		bodyReader = resp.BodyStream()
	}

	// Reading the response body and parsing it with gjson
	body, _ := io.ReadAll(bodyReader)
	// body, _ := io.ReadAll(resp.Body)
	res := gjson.ParseBytes(body)

	// Handling various response statuses and potential errors
	if res.Get("error.code").String() == "-32600" {
		log.Println(res.Get("error").String())
		return TranslationResult{
			Code:    http.StatusNotAcceptable,
			Message: "Invalid target language",
		}, nil
	}

	if resp.StatusCode() == http.StatusTooManyRequests && authKey != "" {
		authKeyArray := strings.Split(authKey, ",")
		for _, authKey := range authKeyArray {
			validity, err := checkUsageAuthKey(authKey, client)
			if err != nil {
				continue
			} else {
				if validity {
					translatedText, err := translateByOfficialAPI(translateText, sourceLang, targetLang, authKey, client)
					if err != nil {
						return TranslationResult{
							Code:    http.StatusTooManyRequests,
							Message: "Too Many Requests",
						}, nil
					}
					return TranslationResult{
						Code:       http.StatusOK,
						Message:    "Success",
						ID:         1000000,
						Data:       translatedText,
						SourceLang: sourceLang,
						TargetLang: targetLang,
						Method:     "Official API",
					}, nil
				}
			}

		}
	} else {
		var alternatives []string
		res.Get("result.texts.0.alternatives").ForEach(func(key, value gjson.Result) bool {
			alternatives = append(alternatives, value.Get("text").String())
			return true
		})
		if res.Get("result.texts.0.text").String() == "" {
			return TranslationResult{
				Code:    http.StatusServiceUnavailable,
				Message: "Translation failed, API returns an empty result.",
			}, nil
		} else {
			return TranslationResult{
				Code:         http.StatusOK,
				ID:           id,
				Message:      "Success",
				Data:         res.Get("result.texts.0.text").String(),
				Alternatives: alternatives,
				SourceLang:   sourceLang,
				TargetLang:   targetLang,
				Method:       "Free",
			}, nil
		}
	}

	return TranslationResult{
		Code:    http.StatusServiceUnavailable,
		Message: "Uknown error",
	}, nil
}

func checkUsageAuthKey(authKey string, client *fasthttp.Client) (bool, error) {
	url := "https://api-free.deepl.com/v2/usage"

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI(url)

	req.Header.Add("Authorization", "DeepL-Auth-Key "+authKey)

	err := client.Do(req, resp)
	if err != nil {
		return false, err
	}

	var response DeepLUsageResponse
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return false, err
	}
	return response.CharacterCount < 499900, nil
}
