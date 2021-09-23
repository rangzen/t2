package service

import (
	"fmt"
	"log"
)

type t2Service struct {
	config T2Config
	tr     Deepl
}

type T2Config struct {
	PrintUsage bool
	SourceLang string
	PivotLang  string
}

func NewT2(config T2Config, tr Deepl) t2Service {
	return t2Service{config: config, tr: tr}
}

func (t2 t2Service) TraductionTranslation(t string) error {
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

	if !t2.config.PrintUsage {
		return nil
	}
	deeplUsage, err := t2.tr.DeeplUsage()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Usage: %d/%d\n", deeplUsage.CharacterCount, deeplUsage.CharacterLimit)
	return nil
}
