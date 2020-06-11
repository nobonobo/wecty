package wecty

import (
	"fmt"
	"strings"
	"syscall/js"
)

var (
	global              = js.Global()
	document            = global.Get("document")
	console             = global.Get("console")
	mountList           []Mounter
	listenerReleaseList []func()
)

// ==================
// interfaces

// Mounter ...
type Mounter interface {
	Mount()
}

// Unmounter ...
type Unmounter interface {
	Unmount()
}

// HTML ...
type HTML interface {
	html() js.Value
}

// Renderer ...
type Renderer interface {
	Render() HTML
}

// Component ...
type Component interface {
	ref() *Core
	Renderer
}

// Applyer ...
type Applyer interface {
	apply(node js.Wrapper)
}

// Markup ...
type Markup interface {
	markup() Applyer
}

// attribute markup
type attribute struct {
	Name  string
	Value interface{}
}

// markup ...
func (a attribute) markup() Applyer {
	return a
}

// apply ...
func (a attribute) apply(node js.Wrapper) {
	node.JSValue().Call("setAttribute", a.Name, a.Value)
}

// Class markup
type Class map[string]bool

// markup ...
func (c Class) markup() Applyer {
	return c
}

// apply ...
func (c Class) apply(node js.Wrapper) {
	classList := node.JSValue().Get("classList")
	for name, ok := range c {
		if ok {
			classList.Call("add", name)
		}
	}
}

// eventMarkup ...
type eventMarkup struct {
	name string
	fn   func(ev js.Value) interface{}
}

func (e eventMarkup) markup() Applyer {
	return e
}

func (e eventMarkup) apply(node js.Wrapper) {
	if n, ok := node.(*Node); ok {
		cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return e.fn(args[0])
		})
		jv := node.JSValue()
		n.ref().listeners = append(n.ref().listeners, func() {
			jv.Call("removeEventListener", e.name, cb, false)
			cb.Release()
		})
		jv.Call("addEventListener", e.name, cb, false)
	}
}

// Core ...
type Core struct {
	last      js.Value
	isNode    bool
	children  []Component
	mount     []bool
	listeners []func()
	update    bool
}

// JSValue ...
func (c *Core) ref() *Core {
	return c
}

func (c *Core) markup() Applyer { return nil }

// html ...
func (c *Core) html() js.Value {
	return c.last
}

// appendChild ...
func (c *Core) appendChild(child Component) {
	c.children = append(c.children, child)
	c.mount = append(c.mount, false)
	c.update = false
}

// textNode ...
type textNode struct {
	Core
	text string
}

// markup ...
func (t *textNode) markup() Applyer { return nil }

func (t *textNode) html() js.Value {
	core := t.Core.ref()
	if !core.update {
		core.last = document.Call("createTextNode", string(t.text))
		core.update = true
	}
	return core.last
}

// Render ...
func (t *textNode) Render() HTML {
	return t
}

// JSValue ...
func (t *textNode) JSValue() js.Value {
	return t.ref().last
}

func (t *textNode) String() string {
	return fmt.Sprintf("%q(%p)", t.text, t)
}

// Node element
type Node struct {
	Core
	ns      string
	tag     string
	markups []Markup
}

func (n *Node) markup() Applyer { return nil }

func (n *Node) String() string {
	return fmt.Sprintf("<%s(%p)>", n.tag, n)
}

// Render ...
func (n *Node) Render() HTML {
	core := n.ref()
	core.children = nil
	core.mount = nil
	markups := []Markup{}
	for _, a := range n.markups {
		if v, ok := a.(Component); ok {
			core.appendChild(v)
		} else {
			markups = append(markups, a)
		}
	}
	n.markups = markups
	return n
}

// JSValue ...
func (n *Node) JSValue() js.Value {
	return n.ref().last
}

// html ...
func (n *Node) html() js.Value {
	core := n.ref()
	if !core.update {
		if len(n.ns) > 0 {
			core.last = document.Call("createElementNS", n.ns, n.tag)
		} else {
			core.last = document.Call("createElement", n.tag)
		}
		for _, a := range n.markups {
			m := a.markup()
			if m != nil {
				m.apply(n)
			}
		}
		for i, c := range n.ref().children {
			if m, ok := c.(Unmounter); ok && n.ref().mount[i] {
				m.Unmount()
				n.ref().mount[i] = false
			}
			html := c.Render()
			if !c.ref().isNode {
				node := html.(*Node)
				node.Render()
				c.ref().children = []Component{node}
			}
			jv := html.html()
			if !jv.IsUndefined() && !jv.IsNull() {
				if m, ok := c.(Mounter); ok {
					n.ref().mount[i] = true
					mountList = append(mountList, m)
				}
				c.ref().last = jv
				core.last.Call("appendChild", jv)
			}
		}
		core.update = true
	}
	return core.last
}

// ==== private funcs ====

func replaceNode(newNode, oldNode js.Value) {
	if newNode.Equal(oldNode) {
		return
	}
	oldNode.Get("parentNode").Call("replaceChild", newNode, oldNode)
}

func cleanup(c Component) {
	core := c.ref()
	for i, child := range core.children {
		if len(core.mount) > 0 {
			if u, ok := child.(Unmounter); ok && core.mount[i] {
				u.Unmount()
				core.mount[i] = false
			}
		}
		cleanup(child)
	}
	core.children = nil
	for _, v := range core.listeners {
		v()
	}
}

func finalize() {
	for _, v := range mountList {
		v.Mount()
	}
	mountList = nil
}

// ==== global funcs ====

// Attr ...
func Attr(key string, value interface{}) Markup {
	return attribute{Name: key, Value: value}
}

// Text Node
func Text(s string) Markup {
	t := &textNode{
		text: s,
	}
	t.Core.isNode = true
	return t
}

// Tag Node
func Tag(name string, markups ...Markup) *Node {
	return &Node{Core: Core{isNode: true}, tag: name, markups: markups}
}

// TagWithNS Node with NameSpace
func TagWithNS(name, ns string, markups ...Markup) *Node {
	return &Node{Core: Core{isNode: true}, tag: name, ns: ns, markups: markups}
}

// Render ...
func Render(c Component) HTML {
	if m, ok := c.(Mounter); ok {
		mountList = append(mountList, m)
	}
	html := c.Render()
	if !c.ref().isNode {
		node := html.(*Node)
		c.ref().children = []Component{node}
	}
	for !c.ref().isNode {
		r, ok := html.(Component)
		if ok {
			html = Render(r)
			c = r
			continue
		}
		break
	}
	return html
}

// Rerender ...
func Rerender(c Component) {
	if u, ok := c.(Unmounter); ok && !c.ref().last.IsUndefined() {
		u.Unmount()
	}
	target := c.ref().last
	if target.IsNull() || target.IsUndefined() {
		panic("rerender not renderd node")
	}
	cleanup(c)
	act := document.Get("activeElement").Get("id").String()
	newNode := Render(c).html()
	replaceNode(newNode, target)
	if f := document.Call("getElementById", act); !f.IsNull() {
		f.Call("focus")
	}
	target.Call("remove")
	c.ref().last = newNode
	finalize()
}

// RenderBody ...
func RenderBody(c Component) {
	target := document.Call("querySelector", "body")
	html := Render(c)
	if !c.ref().isNode && html != nil {
		node := html.(*Node)
		c.ref().children = []Component{node}
	}
	body := html.html()
	if strings.ToLower(body.Get("tagName").String()) != "body" {
		panic("top level element must be 'body'")
	}
	replaceNode(body, target)
	c.ref().last = body
	finalize()
}

// Event ...
func Event(name string, fn func(ev js.Value) interface{}) Markup {
	return eventMarkup{name, fn}
}

// SetTitle sets the title of the document.
func SetTitle(title string) {
	document.Set("title", title)
}

// AddMeta ...
func AddMeta(name, content string) {
	meta := document.Call("createElement", "meta")
	meta.Set("name", name)
	meta.Set("content", content)
	document.Get("head").Call("appendChild", meta)
}

// AddStylesheet adds an external stylesheet to the document.
func AddStylesheet(url string) {
	link := document.Call("createElement", "link")
	link.Set("rel", "stylesheet")
	link.Set("href", url)
	document.Get("head").Call("appendChild", link)
}

// LoadScript ...
func LoadScript(url string) {
	ch := make(chan bool)
	script := Tag("script",
		Attr("src", url),
		Event("load", func(js.Value) interface{} {
			println("load!")
			close(ch)
			return nil
		}),
	).Render().html()
	document.Get("head").Call("appendChild", script)
	<-ch
}

// LoadModule ...
func LoadModule(names []string, url string) <-chan js.Value {
	ch := make(chan js.Value)
	js.Global().Set("__wecty_send__", Callback1(func(obj js.Value) interface{} {
		ch <- obj
		return nil
	}))
	js.Global().Set("__wecty_close__", Callback0(func() interface{} {
		close(ch)
		return nil
	}))
	lines := []string{}
	for _, name := range names {
		lines = append(lines, fmt.Sprintf("__wecty_send__(%s);", name))
	}
	lines = append(lines, "__wecty_close__();")
	script := Tag("script",
		Attr("type", "module"),
		Text(fmt.Sprintf("import { %s } from %q;\n%s",
			strings.Join(names, ", "),
			url,
			strings.Join(lines, "\n"),
		)),
	).Render().html()
	document.Get("head").Call("appendChild", script)
	return ch
}

// RequestAnimationFrame ...
func RequestAnimationFrame(callback func(float64)) int {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		cb.Release()
		if callback != nil {
			callback(args[0].Float())
		}
		return js.Undefined()
	})
	return global.Call("requestAnimationFrame", cb).Int()
}

/*
// PrintTree ...
func PrintTree(c Component, indent int) {
	log.Printf("%s%v", strings.Repeat("  ", indent), c)
	log.Printf("%s%[2]T(%#[2]v)", strings.Repeat("  ", indent), c.ref().last)
	for _, child := range c.ref().children {
		PrintTree(child, indent+1)
	}
}
*/
