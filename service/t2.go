package service

import (
	"fmt"
	"log"
)

type T2Service struct {
	config T2Config
	tr     Deepl
}

type T2Config struct {
	SourceLang string
	PivotLang  string
}

func NewT2(config T2Config, tr Deepl) T2Service {
	return T2Service{config: config, tr: tr}
}

func (t2 T2Service) TraductionTranslation(t string) error {
	fmt.Println("# Original text")
	fmt.Println(t)
	firstPass, err := t2.tr.DeeplTranslate(t, t2.config.SourceLang, t2.config.PivotLang)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("# Pivot text")
	fmt.Println(firstPass.Translations[0].Text)
	secondPass, err := t2.tr.DeeplTranslate(firstPass.Translations[0].Text, t2.config.PivotLang, t2.config.SourceLang)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("# Double translated text")
	fmt.Println(secondPass.Translations[0].Text)
	return nil
}
