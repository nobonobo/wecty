package components

import (
	"github.com/nobonobo/wecty"
)

// ItemList ...
type ItemList struct {
	wecty.Core
}

// Items ...
func (c *ItemList) Items() *wecty.Node {
	items := []wecty.Markup{
		&Item{Title: "Item1"},
		&Item{Title: "Item2"},
		&Item{Title: "Item3"},
	}
	return wecty.Tag("ul", items...)
}
