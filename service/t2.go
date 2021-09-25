package service

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
)

type T2Service struct {
	config T2Config
	tr     Translation
}

type T2Config struct {
	SourceLang string
	PivotLang  string
}

func NewT2(config T2Config, tr Translation) T2Service {
	return T2Service{config: config, tr: tr}
}

func (t2 T2Service) TraductionTranslation(t string) error {
	fmt.Println("# Original text")
	fmt.Println(t)
	firstPass, err := t2.tr.Translate(t, t2.config.SourceLang, t2.config.PivotLang)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("# Pivot text")
	fmt.Println(firstPass.Text)
	secondPass, err := t2.tr.Translate(firstPass.Text, t2.config.PivotLang, t2.config.SourceLang)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("# Double translated text")
	fmt.Println(secondPass.Text)

	fmt.Println("# Diff version")
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(t, secondPass.Text, false)
	fmt.Println(dmp.DiffPrettyText(diffs))

	return nil
}
