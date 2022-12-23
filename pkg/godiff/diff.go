/*
Copyright © 2021 Cedric L'homme <public@l-homme.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

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
