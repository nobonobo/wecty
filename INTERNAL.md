# 内部情報

### 型の互換性

| type         | ret.(alias)     | isMarkup | isApplyer | isRenderer | isHTML | isWrapper |
| ------------ | --------------- | -------- | --------- | ---------- | ------ | --------- |
| Attr()       | Markup          | true     | true      | false      | false  | false     |
| Class        | map[string]bool | true     | true      | false      | false  | false     |
| Event()      | Markup          | true     | true      | false      | false  | false     |
| Text()       | Markup          | true     | true      | true       | true   | true      |
| Tag()        | \*Node          | true     | false     | true       | true   | true      |
| ユーザー定義 | 　 Component    | true     | false     | true       | false  | true      |

- Markup は markup()により Aplyer に変換
- Applyer は apply()により Wrapper(DOM ノード)を加工する
- Renderer は Render()により HTML に変換（ユーザー定義可能）
- HTML は html()により js.Value に変換（実際の DOM ノード）

コンポーネントとしての基礎機能は wecty.Core 構造体に実装済み。

ユーザー定義構造体をコンポーネントにするためには wecty.Core を埋め込んでおき、
あとは「Render()HTML」メソッドをサポートするだけ。

### マークアップ関数

- Tag(tagName string, markups ...Markup) \*Node
- Attr(key string, value interface{}) Markup
- Class{key: true/false, ...} as Markup
- Event(name string, fn func(ev js.Value) interface{}) Markup
- Text(text string) Markup

#### 描画関数

- RenderBody(c Component): Body ノードを差し替える
- Rerender(c Component): 予め RenderBody で描画済みのコンポーネントを再描画する

#### ユーティリティ関数

- SetTitle(title string)
- AddMeta(name, content string)
- AddStylesheet(url string)
- LoadScript(url string)
- LoadModule(names []string, url string) <-chan js.Value

##### LoadModule 例

```go
names := []string{"Server", "Client"}
defs := map[string]js.Value{}
for obj := range wecty.LoadModule(names, "jsonrpclib.js") {
  defs[names[0]] = obj
  names = names[1:]
}
```

以下の Javascript(module) 記述と等価

```javascript
import { Server , Client } as defs from "jsonrpclib.js";
```

### ルーター

- NewRouter() \*Router
- type Router
  - Handle(path string, fn func(path string))
  - Start() error
  - Navigate(path string) error
  - Current() \*url.URL

##　 HTML マークアップの方針

- HTML ライクな記述を Go の記述に変換
- ややこしい定義を書くところは最後には Go の記述に対応する構文で書かせて解決
- 属性値を得る記述とタグ部分定義を得る記述する機能を提供するだけ
- あとは HTML 記述をストレートに Go の記述に変換するだけ

### HTML マークアップ（実装済み）

- `<import as="cmp">github.com/nobonobo/wecty/examples/todo/components</import>`
  - 他所のパッケージを参照する必要がある場合に記述
  - `import (cmp "github.com/nobonobo/wecty/examples/todo/components")`に置換
  - as 属性は省略可能
- `<raw>wecty.Markup値を得るGo記述を直接記述</raw>`
  - Node を記述する際に Go の記述を使う（主にコンポーネントインスタンスを書く）
- `"{{Goの記述として解釈して値ソースにする}}"`
  - TextNode: 文字列化して wecty.Text でラップ
  - 属性値: 文字列化して属性値に利用

## Conditional マークアップ（WIP）

`<if cond="">...</if>`条件の値が true じゃない場合は空の Markup を返す。

- ただし、cond や直接 Go 記述の内容が文法的に正しいのかは確認しない。

####　 Class マークアップ（WIP）

以下の２記述は同じ結果を出力する（順序不同）
ただし、後者は動的な値を使ってオンオフできる

```html
<div class={{
  hoge: true,
  moge: true,
}}></div>
```

```html
<div class="hoge moge"></div>
```

例：

```html
<div class={{
  hoge: c.PropHoge,
  moge: c.PropMoge,
}}></div>
```
