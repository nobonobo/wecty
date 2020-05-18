package wecty

import (
	"log"
	"strings"
	"syscall/js"
)

var (
	global     = js.Global()
	document   = global.Get("document")
	console    = global.Get("console")
	mountList  []Mounter
	deleteNode []js.Value
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

type eventReleaser struct {
	node  js.Value
	event string
	cb    js.Func
}

func (er *eventReleaser) Release() {
	er.node.Call("removeEventListener", er.event, er.cb)
	er.cb.Release()
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
	n := node.(*Node)
	core := n.ref()
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return e.fn(args[0])
	})
	core.listeners = append(core.listeners, cb)
	/*
		core.listeners = append(core.listeners, &eventReleaser{
			node:  node.JSValue(),
			event: e.name,
			cb:    cb,
		})
	*/
	node.JSValue().Call("addEventListener", e.name, cb, false)
}

// Core ...
type Core struct {
	current   Component
	parent    Component
	last      js.Value
	isNode    bool
	children  []Component
	listeners []js.Func
	//listeners []*eventReleaser
	update bool
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

func (c *Core) cleanup() {
	for _, child := range c.children {
		child.ref().cleanup()
		if last := child.ref().last; !last.IsNull() {
			deleteNode = append(deleteNode, last)
		}
	}
	for _, l := range c.listeners {
		l.Release()
	}
	c.listeners = nil
}

// appendChild ...
func (c *Core) appendChild(child Component) {
	c.children = append(c.children, child)
	c.update = false
}

// textNode ...
type textNode struct {
	Core
	text string
}

// markup ...
func (t *textNode) markup() Applyer { return nil }

// Render ...
func (t *textNode) Render() HTML {
	return nil
}

// JSValue ...
func (t *textNode) JSValue() js.Value {
	return t.ref().last
}

// html ...
func (t *textNode) html() js.Value {
	core := t.Core.ref()
	if !core.update {
		core.last = document.Call("createTextNode", string(t.text))
		core.update = true
	}
	return core.last
}

// Node element
type Node struct {
	Core
	tag     string
	markups []Markup
}

func (n *Node) markup() Applyer { return nil }

// Render ...
func (n *Node) Render() HTML {
	core := n.ref()
	for _, a := range n.markups {
		if v, ok := a.(Component); ok {
			v.Render()
			core.appendChild(v)
		}
	}
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
		core.current = n
		core.last = document.Call("createElement", n.tag)
		for _, a := range n.markups {
			m := a.markup()
			if m != nil {
				m.apply(n)
			}
		}
		for _, c := range n.ref().children {
			jv := renderDOMNode(c)
			if !jv.IsUndefined() {
				if m, ok := c.ref().current.(Mounter); ok {
					mountList = append(mountList, m)
				}
				c.ref().last = jv
				core.last.Call("appendChild", jv)
				c.ref().parent = n
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

// renderDOMNode ...
func renderDOMNode(c Component) js.Value {
	if !c.ref().last.IsUndefined() {
		return c.ref().last
	}
	c.ref().current = c
	html := c.Render()
	if !c.ref().isNode && html != nil {
		node := html.(*Node)
		c.ref().children = []Component{node}
		if !node.ref().last.IsUndefined() {
			return node.ref().last
		}
		d := node.JSValue()
		if d.IsUndefined() {
			node.Render()
			return node.html()
		}
		return d
	}
	if html != nil {
		return html.html()
	}
	if r, ok := c.(HTML); ok {
		return r.html()
	}
	panic("invalid component")
}

func requestAnimationFrame(callback func(float64)) int {
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

func finalize() {
	//requestAnimationFrame(nil)
	for _, v := range deleteNode {
		if parent := v.Get("parentNode"); !parent.IsNull() {
			parent.Call("removeChild", v)
		}
	}
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
	return &textNode{
		text: s,
	}
}

// Tag Node
func Tag(name string, markups ...Markup) *Node {
	return &Node{Core: Core{isNode: true}, tag: name, markups: markups}
}

// Render ...
func Render(c Component) HTML {
	if m, ok := c.(Mounter); ok {
		mountList = append(mountList, m)
	}
	html := c.Render()
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
	c.ref().cleanup()
	if u, ok := c.ref().current.(Unmounter); ok {
		u.Unmount()
	}
	act := document.Get("activeElement").Get("id").String()
	target := renderDOMNode(c)
	html := Render(c)
	newNode := html.html()
	c.ref().last = newNode
	if n, ok := c.ref().children[0].(*Node); ok {
		n.ref().last = newNode
	}
	replaceNode(newNode, target)
	if f := document.Call("getElementById", act); !f.IsNull() {
		f.Call("focus")
	}
	finalize()
}

// RenderBody ...
func RenderBody(c Component) {
	if u, ok := c.ref().current.(Unmounter); ok {
		u.Unmount()
	}
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
	c.ref().last = body
	replaceNode(body, target)
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

// AddScript adds an external script to the document.
func AddScript(url string) {
	script := document.Call("createElement", "script")
	script.Set("src", url)
	document.Get("head").Call("appendChild", script)
}

// PrintTree ...
func PrintTree(c Component, indent int) {
	log.Printf("%s%[2]T(%[2]p)", strings.Repeat("  ", indent), c)
	log.Printf("%s%[2]T(%#[2]v)", strings.Repeat("  ", indent), c.ref().last)
	for _, child := range c.ref().children {
		PrintTree(child, indent+1)
	}
}
