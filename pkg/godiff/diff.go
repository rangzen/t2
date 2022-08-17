package godiff

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

type Diff struct{}

func (Diff) Print(a, b string) string {
	dmp := diffmatchpatch.New()
	d := dmp.DiffMain(a, b, false)
	return dmp.DiffPrettyText(d)
}
