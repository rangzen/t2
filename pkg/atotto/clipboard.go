package atotto

import "github.com/atotto/clipboard"

type Clipboard struct{}

func (Clipboard) Write(t string) error {
	return clipboard.WriteAll(t)
}
