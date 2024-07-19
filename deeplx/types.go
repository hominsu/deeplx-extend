package deeplx

type Text struct {
	Text                string `json:"text"`
	RequestAlternatives int    `json:"requestAlternatives"`
}

type Lang struct {
	SourceLangUserSelected string `json:"source_lang_user_selected"`
	TargetLang             string `json:"target_lang"`
}

type CommonJobParams struct {
	WasSpoken       bool   `json:"wasSpoken"`
	TranscribeAS    string `json:"transcribe_as"`
	RegionalVariant string `json:"regionalVariant,omitempty"`
}

type Params struct {
	Texts           []Text          `json:"texts"`
	Splitting       string          `json:"splitting"`
	Lang            Lang            `json:"lang"`
	Timestamp       int64           `json:"timestamp"`
	CommonJobParams CommonJobParams `json:"commonJobParams"`
}

type PostData struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      int64  `json:"id"`
	Params  Params `json:"params"`
}

type TranslationResult struct {
	Code         int
	ID           int64
	Message      string
	Data         string
	Alternatives []string
	SourceLang   string
	TargetLang   string
	Method       string
}

type PayloadAPI struct {
	Text       []string `json:"text"`
	TargetLang string   `json:"target_lang"`
	SourceLang string   `json:"source_lang"`
}

type Translation struct {
	Text string `json:"text"`
}

type TranslationResponse struct {
	Translations []Translation `json:"translations"`
}

type DeepLUsageResponse struct {
	CharacterCount int `json:"character_count"`
	CharacterLimit int `json:"character_limit"`
}
