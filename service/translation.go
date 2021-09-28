package service

type TranslationResponse struct {
	Text string
}

type UsageResponse struct {
	Used  int64
	Limit int64
}

type Translation interface {
	Name() string
	Translate(text string, source string, target string) (TranslationResponse, error)
	Usage() (UsageResponse, error)
}
