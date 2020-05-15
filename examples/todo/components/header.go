package components

import (
	"github.com/nobonobo/wecty"
)

//go:generate wecty generate -c Header -p components header.html

// Header ...
type Header struct {
	wecty.Core
}
