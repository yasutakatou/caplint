# caplint

Markdownとかテキスト系じゃない、パワポやブラウザのスライドツールでも、サクサクとLinterでチェックできたらいいよね！を実装したツール

# こんなことありませんか？

スライド書いてて、Linterをかけようとしたときに、Markdownとかテキスト形式だったらtextlintをVSCode拡張やショートカットで呼べばいいけど、<br>
パワポとか、Googleスライドなんかで書いている時は、[pptx2md](https://github.com/ssine/pptx2md)辺りで一回、Markdownに変換してからLinterにかけますよね。<br>
となると、ショートカットみたいなのでLinterかけるのにやってやれなくは無いけど、ちょいちょい手間だし、まとめてやるとウンザリするくらい修正箇所出たりする<br>
（自分で書いた文章の修正箇所にウンザリしてもねぇ、、）<br>
それよっか、ページ単位くらいでサクッとLinterかけれないかなぁ、できればテキスト形式関係なく出来ると良いんだけど、、と思ってたらツール書いてた<br>

# 機能

OCRとLinterを組み合わせて文字認識にLinterをかけるツールです！<br>
すまん、手元にWindowsしか無いのでMacでは動かん！！

- ① 任意のファイルにOCRをかけてLinterに流せます
- ② クリップボードにある画像にOCRをかけてLinterに流せます
- ③ ショートカットキー打ち込むとクリップボードにある画像にOCRをかけてLinterに流せます
- ④ ショートカットキー打ち込むと現在のアクティブウィンドウをキャプチャしてOCRをかけてLinterに流せます

④のモードはウィンドウまるごとキャプチャするのでメニューバーとかに文字がある場合はそれも認識しちゃうんで実用性はそこそこ・・・<br>
③が便利なので良いんじゃないですかね、一画面事にバンバンOCR→Linterかけてチェックできるので便利！

# 動作画面

- ① 任意のファイルにOCRをかけてLinterに流せます

![image](https://github.com/user-attachments/assets/22b5cb69-6750-4e52-ba2e-50f5bd386d35)

- ② クリップボードにある画像にOCRをかけてLinterに流せます

![image](https://github.com/user-attachments/assets/daf24849-3cc7-4bac-9859-912167a6b4b0)

- ③ ショートカットキー打ち込むとクリップボードにある画像にOCRをかけてLinterに流せます



- ④ ショートカットキー打ち込むと現在のアクティブウィンドウをキャプチャしてOCRをかけてLinterに流せます

![1](https://github.com/user-attachments/assets/32ca80fe-7fba-44ab-bc0a-dee23b3d1757)

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

なお、内部処理としてOCRの認識精度を上げるために読み込んだ画像をリサイズした画像を生成します。デフォルトは倍のサイズのpngファイルが出力されます

# 設定ファイル

OCRとLinterのチューニングとカスタマイズが出来るよう。設定ファイルから外部アプリを呼び出すようにしています。<br>
デフォルトは　caplint.ini　というファイル名です。

## [tesseract]　セクション

OCRにあたるアプリを定義します。{}の部分にテンポラリで生成されるファイル名が入ります（置換して実行されます）

## [textlint]　セクション

Linterにあたるアプリを定義します。{}の部分にテンポラリで生成されるファイル名が入ります（置換して実行されます）

```
[tesseract]
tesseract.exe {}.png {} -l jpn --psm 6

[textlint]
npx textlint {}.txt --format pretty-error
```

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

trueを指定すると今クリップボードにある画像イメージを元にOCRをかけてLinterかけます

## -config

コンフィグファイルをデフォルトの　caplint.ini　から変えたいときにこれでコンフィグファイル名を指定します

## -debug

デバッグモードで動作します。色々出力されます

## -file

pngファイル名を指定すると、その画像ファイルを元にOCRかけてLinterしてくれるモードです

## -log

デバッグモードで出たログを出力するオプションです

## -nodelete

これを指定するとテンポラリで出力した画像ファイルなりを消さずに動作します

## -resize

リサイズする倍率を変えます。あまり大きくし過ぎるとTesseract読み込み時にエラーになるので注意

## -shortcut

ショートカット動作モードです。このオプションで起動するとバックグラウンドでアプリが動作し続け、特定のショートカットキー入力で動作するようになります

## -shortcutclipboard

ショートカット動作の時にクリップボードの画像からキャプチャさせる時のキーアサインです。デフォルトは z (Shift+Ctrl+z)です

## -shortcutwindow

ショートカット動作の時にフォアグラインドのウィンドウからキャプチャさせる時のキーアサインです。デフォルトは a (Shift+Ctrl+a)です

## -shortexit

ショートカット動作の時にアプリを終了させる時のキーアサインです。デフォルトは q (Shift+Ctrl+q)です

# その他

OCRの読み取りうまくいかんわ、精度チューニングしたいわ、だったら↓のリンクが参考になります

[Tesseractの日本語精度を上げる](https://scrapbox.io/villagepump/Tesseract%E3%81%AE%E6%97%A5%E6%9C%AC%E8%AA%9E%E7%B2%BE%E5%BA%A6%E3%82%92%E4%B8%8A%E3%81%92%E3%82%8B)

# ライセンス
Apache License 2.0<br>
MIT License<br>
ISC License


