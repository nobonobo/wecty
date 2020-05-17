# Wecty

Wecty: フロントエンドツールキット for Go and TinyGo

- Wecty は [Vecty(github.com/gopherjs/vecty)](https://github.com/gopherjs/vecty) のアイディアを元にしています
- Wecty は Vecty よりも単純に実装されました
- Wecty は WASM アーキテクチャのみをサポートします

## 方針

- 基本的な機能を一通り含みます
- リッチな機能や冗長な機能は採用しません
- それにより出力される WASM サイズを小さく保ちます
- 描画の最速を目指したりはしない

## 機能

- SetTitle: ドキュメントタイトルの変更
- AddStylesheet: スタイルシートの追加
- AddScript: スクリプトの追加
- AddMeta: メタヘッダの追加
- RenderBody/Rerender: ボディコンポーネントの描画とコンポーネントの再描画
- ネストされたコンポーネントのサポート
- Mount/Unmount: マウントとアンマウントタイミングのフック
- Tag: HTML タグのマークアップ
- Attr: タグの属性をマークアップ
- Class: タグのクラスをマークアップ
- Event: タグのイベントをマークアップ
- Router: シンプルな SPA ルーター(URL のハッシュ利用)
- Utilities:
  - wecty generate: go-generate 用ツール、HTML 記述から Go のコードを生成
  - wecty server: 開発用サーバー

## 基本の使い方

### ツールのセットアップ

```shell
go get github.com/nobonobo/wecty/cmd/wecty
```

### ファイル群の準備

- 以下のファイル(後述)を作成する
  - top.go
  - top.html
  - main.go
- go mod init sample

top.html

```html
<form @submit="{{c.OnSubmit}}">
  <button>Submit</button>
</form>
```

以下のコマンドを実行した場合、

```shell
wecty generate -c TopView top.html
```

top_gen.go が生成されます

```go
package main

import (
  "github.com/nobonobo/wecty"
)

func (c *TopView) Render() wecty.HTML {
  return wecty.Tag("form",
    wecty.Event("submit", c.OnSubmit),
    wecty.Tag("button", vecty.Text("Submit")),
  )
}
```

以下のような Go コードを手書きで書いておくことでコンポーネントとして利用可能になります。

top.go

```go
package main

import (
  "github.com/nobonobo/wecty"
)

type TopView struct {
  wecty.Core
}

func (c *TopView) OnSubmit(ev js.Value) interface{} {
  println("submit!")
  return nil
}
```

また以下の記述を top.go に加えておくと go generate で wecty generate が自動的に走るようになります。

```go
//go:generate wecty generate -c TopView -p main top.html
```

main.go

```go
package main

import (
  "github.com/nobonobo/wecty"
)

func main() {
  wecty.RenderBody(wecty.Tag("body", &TopView{}))
  select {}
}
```

あとは`wecty server`を起動しておけば、
[http://localhost:8080](http://localhost:8080)をブラウザで開くだけで
top_gen.go が生成され WASM がビルド＆サーブされブラウザで動作を開始する

## 開発の進め方

- 上記のセットアップが完了したならば
- top.html の編集ー＞ブラウザリロードで反映結果を確認できます
- top.go や main.go を修正後、ブラウザリロードで修正後の wasm モジュールの挙動を確認できます

## Q&A

- Q: なぜ Vecty とは別に作ったの？

  A: GopherJS の開発が停滞しつつあること、Vecty は GopherJS と Go 両対応により複雑な実装になっている。

- Q: Router はなぜハッシュベース？

  A: URL を書き換えるスタイルはプロキシサーバーの URL 割り当てと整合をとる必要がある。ハッシュベースは単一の URL を振り向けるだけで動作する。つまり、SPA コンテンツを S3 に置いた場合でも動作する。

- Q: Vecty のように prop や event パッケージを設けないのはなぜ？

  A: 基本のマークアップは wecty generate の出力に任せるのでマークアップの容易さは無用だった。それにそれらのパッケージが WASM サイズの肥大化を招いていた。

- Q: wecty generate のマークアップ機能が足りないのはなぜ？

  A: 頑張っても Go の手書きの自由度を超えることはできない。ユーザーは wecty generate で済ますか Go のコードで細かく書いてコンポーネントを実装するかを使い分けてもらいたい。

- Q: コンポーネントより細かい単位の最適な DOM ツリーの更新をしないのはなぜ？

  A: 軽量な実装で仮装 DOM を実装した。賢くすることで描画更新は早くなるかもしれないが、WASM サイズが膨れてしまうことは避けたかった。速度を極力落とさない作りはユーザーがチャレンジできる。それはコンポーネントの粒度を小さく保つこと。

## コンポーネントに成れる構造体の条件

- wecty.Core を埋め込みした構造体定義
- `Render() wecty.HTML`メソッドを持つこと

## DOM ツリーのマークアップ

`コンポーネント.Render() HTML`メソッドを実装するには
単一の wecty.Tag(...)の戻り値を return するように記述する。

```go
func (c *コンポーネント) Render() HTML {
  return wecty.Tag(...)
}
```

wecty.Tag 関数の定義は以下の通り

`wecty.Tag(tagName, markups ...wecty.Markup) *Node`

tagName には HTML タグ名を渡す。markups には以下の記述を書く。

Markup になれるもの一覧

- wecty.Attr(...)の戻り値（親 Tag の属性値になる）
- wecty.Class{}値（親 Tag の classList 値になる）
- wecty.Event(...)の戻り値（親 Tag にイベントリスナーを追加）
- wecty.Text(...)の戻り値（子ノードを追加）
- wecty.Tag(...)の戻り値（子ノードを追加）
- Component を満たすオブジェクト（子ノードを追加）

## ユーティリティ

wecty いくつかのサブコマンドをもつユーティリティ

- wecty generate: HTML 記述から Wecty 用の Go コードを生成するツール
- wecty server: 開発用簡易 HTTP サーバー

## wecty generate

```
Usage of generate:
  -c string
    	component name
  -o string
    	output filename
  -p string
    	output package name (default "main")
```

- default `output filename`=`basename_gen.go`
- default `package name`=`main`
- require `component name`

### HTML ライクなマークアップの基本

`<body></body>` -> `wecty.Tag("body")`

### 属性マークアップ

`<h1 id="title">Title</h1>`:

```go
wecty.Tag("h1",
  wecty.Attr("id", "title"),
  vecty.Text("Title"),
)
```

### Class マークアップ

`<h1 class="title large">Title</h1>`:

```go
wecty.Tag("h1",
  wecty.Class{
    "title": true,
    "large": true,
  },
  vecty.Text("Title"),
)
```

### イベントマークアップ

```html
<form @submit="{{c.OnSubmit}}">
  <button>Submit</button>
</form>
```

```go
wecty.Tag("form",
  wecty.Event("submit", c.OnSubmit),
  wecty.Tag("button", vecty.Text("Submit")),
)
```

### RAW 記述

```html
<import>github.com/nobonobo/examples/todo/components</import>
<div><raw>&components.Item{}</raw></div>
```

```go
import (
  "github.com/nobonobo/examples/todo/components"
)
...
  return wecty.Tag("div",
    &components.Item{},
  )
...
```

## wecty server

開発用簡易 HTTP サーバー

```
Usage of server:
  -addr string
    	listen address (default ":8080")
  -tinygo
    	use tinygo tool chain
```

### 機能

- 静的コンテンツのサーブ
- "main.wasm"を要求されるとその該当フォルダで go-generate と WASM のビルド、gzip が行われその結果をサーブします
- index.html が無いところで要求されたら標準的な WASM 読み込み用 HTML を返します
- "wasm_exec.js"リソースが要求されたら適切な wasm_exec.js をサーブします

### コマンドオプション

-tinygo: WASM ビルドに tinygo を使う

## 出力サイズ

Todo サンプルのコンパイル事例

| ツール | WASM サイズ | gzipped  |
| ------ | ----------- | -------- |
| Go     | 2.8MiB      | 787.2KiB |
| TinyGo | 374KiB      | 154KiB   |

## 既知の問題

- TinyGo コンパイラはまだ Go-Module サポートがありません。GOPATH、GO111MODULE=off などの環境変数設定が必要です
- TinyGo の WASM 出力は js.finalizeRef が未実装なためメモリーリークが起こりえます
- TinyGo 0.13.1 は log.Print 系を使うとデッドロックするバグがあります（dev 版では修正済み？）

## 今後の機能追加

- 条件別マークアップ機能の提供
- デプロイ用の静的ファイルセットをエクスポートする支援機能の提供
- 新規プロジェクト生成機能の提供
- ~~コマンドツールの統合＆サブコマンドによる多機能化~~
