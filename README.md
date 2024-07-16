# caplint

Markdownとかテキスト系じゃない、パワポやブラウザのスライドツールでも、サクサクとLinterでチェックできたらいいよね！を実装したツール

# こんなことありませんか？

スライド書いてて、Linterをかけようとしたときに、Markdownとかテキスト形式だったらtextlintをVSCode拡張やショートカットで呼べばいいけど、
パワポとか、Googleスライドなんかで書いている時は、[pptx2md](https://github.com/ssine/pptx2md)辺りで一回、Markdownに変換してからLinterにかけますよね。
となると、ショートカットみたいなのでLinterかけるのにやってやれなくは無いけど、ちょいちょい手間だし、まとめてやるとウンザリするくらい修正箇所出たりする
（自分で書いた文章の修正箇所にウンザリしてもねぇ、、）
それよっか、ページ単位くらいでサクッとLinterかけれないかなぁ、できればテキスト形式関係なく出来ると良いんだけど、、と思ってたらツール書いてた

# 機能

OCRとLinterを組み合わせて文字認識にLinterをかけるツールです！
すまん、手元にWindowsしか無いのでMacでは動かん！！

- ① 任意のファイルにOCRをかけてLinterに流せます
- ② クリップボードにある画像にOCRをかけてLinterに流せます
- ③ ショートカットキー打ち込むとクリップボードにある画像にOCRをかけてLinterに流せます
- ④ ショートカットキー打ち込むと現在のアクティブウィンドウをキャプチャしてOCRをかけてLinterに流せます

④のモードはウィンドウまるごとキャプチャするのでメニューバーとかに文字がある場合はそれも認識しちゃうんで実用性はそこそこ・・・
③が便利なので良いんじゃないですかね、一画面事にバンバンOCR→Linterかけてチェックできるので便利！

# インストール方法

## 前提

### Tesseractのセットアップ

[Tesseract OCR をWindowsにインストールする方法](https://gammasoft.jp/blog/tesseract-ocr-install-on-windows/)

この辺の手順でWindowsにTesseractをインストールしてください。デフォでパスが通らないので↓とかでパス通してください

```
set PATH=%PATH%;c:\Program Files\Tesseract-OCR
```

適当な画像を読ませて動作チェックするとよろしげです

### textlintのセットアップ

[textlintで日本語の校正を行う](https://zenn.dev/yamane_yuta/articles/65886897cefa1e)

この辺の手順でWindowsにtextlintをインストールしてください。あと↓のかんじでtextlintのルールもついでに入れてください

```
npm install -D textlint-rule-preset-ja-spacing textlint-rule-preset-ja-technical-writing textlint-rule-no-mix-dearu-desumasu textlint-rule-ja-no-redundant-expression  @textlint-ja/textlint-rule-no-synonyms sudachi-synonyms-dictionary textlint-rule-ja-no-orthographic-variants
```

適当なテキストファイルを読ませて動作チェックするとよろしけです(npx textlint --initやってなかったとかありがち)

## ツール本体 

ここまで出来たらツールをこのリポジトリからもってきます

```
go get github.com/yasutakatou/caplint
```

それかクローンしてビルドするか

```
git clone https://github.com/yasutakatou/caplint
cd caplint
go build .
```

面倒なら[ここにバイナリファイル置いておくので](https://github.com/yasutakatou/caplint/releases)手元で展開するでもOKです

# アンインストール方法

## Tesseractとtextlint

色々ありそうなのでググってくださいー

## ツール本体

Go言語なのでバイナリファイル消してあげればOK！（流石Goシンプルでサイコー）

# 使い方

- 任意のファイルにOCRをかけてLinterに流せます
- クリップボードにある画像にOCRをかけてLinterに流せます
- ショートカットキー打ち込むとクリップボードにある画像にOCRをかけてLinterに流せます
- ショートカットキー打ち込むと現在のアクティブウィンドウをキャプチャしてOCRをかけてLinterに流せます

# オプション

```
-clipboard
      [-clipboard=input clipboard image (true is enable)]
-config string
      [-config=config file)] (default "caplint.ini")
-debug
      [-debug=debug mode (true is enable)]
-file string
      [-file=exists png file)] (default "text.png")
-log
      [-log=logging mode (true is enable)]
-nodelete
      [-nodelete=no delete temp file mode (true is enable)]
-resize int
      [-resize=resize count (default x2)] (default 2)
-shortcut
      [-shortcut=shortcut key mode (true is enable)]
-shortcutclipboard int
      [-shortcutclipboatrd=input clipboard image when shotcut key mode (default 'z')] (default 90)
-shortcutwindow int
      [-shortcutwindow=input forground window when shotcut key mode (default 'a')] (default 65)
-shortexit int
      [-shortcutexit=shotcut key mode exit (default 'q')] (default 81)
```

## -clipboard
## -config
## -debug
## -file
## -log
## -nodelete
## -resize
## -shortcut
## -shortcutclipboard
## -shortcutwindow
## -shortexit

# ライセンス
Apache License 2.0
MIT License
ISC License


