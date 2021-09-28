package service

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type T2 struct {
	config Config
	ts     TranslationService
}

type Config struct {
	SourceLang string
	PivotLang  string
	OnlyDiff   bool
}

func NewT2(c Config, ts TranslationService) T2 {
	return T2{config: c, ts: ts}
}

func (t2 T2) TraductionTranslation(t string) error {
	if !t2.config.OnlyDiff {
		fmt.Println("# Original text")
		fmt.Println(t)
	}

	firstPass, err := t2.ts.Translate(t, t2.config.SourceLang, t2.config.PivotLang)
	if err != nil {
		return err
	}
	if !t2.config.OnlyDiff {
		fmt.Printf("# Pivot text (%s -> %s by %s)\n", t2.config.SourceLang, t2.config.PivotLang, t2.ts.Name())
		fmt.Println(firstPass.Text)
	}

	secondPass, err := t2.ts.Translate(firstPass.Text, t2.config.PivotLang, t2.config.SourceLang)
	if err != nil {
		return err
	}
	if !t2.config.OnlyDiff {
		fmt.Printf("# Double translated text (%s -> %s by %s)\n", t2.config.PivotLang, t2.config.SourceLang, t2.ts.Name())
		fmt.Println(secondPass.Text)
	}

	if !t2.config.OnlyDiff {
		fmt.Println("# Diff version")
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(t, secondPass.Text, false)
	fmt.Println(dmp.DiffPrettyText(diffs))

	return nil
}
