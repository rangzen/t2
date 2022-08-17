package t2

type TranslationResponse struct {
	Text string
}

type UsageResponse struct {
	Used  int64
	Limit int64
}

type TranslationService interface {
	Name() string
	Translate(text string, source string, target string) (TranslationResponse, error)
	Usage() (UsageResponse, error)
}

type Clipboard interface {
	Write(t string) error
}

type Diff interface {
	Print(a, b string) string
}

type Config struct {
	SourceLang      string
	PivotLang       string
	DiffOnly        bool
	CopyToClipboard bool
}
