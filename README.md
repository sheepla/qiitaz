<div align="right">

![CI](https://github.com/sheepla/fzwiki/actions/workflows/ci.yml/badge.svg)
![Relase](https://github.com/sheepla/fzwiki/actions/workflows/release.yml/badge.svg)

<a href="https://github.com/sheepla/qiitaz/releases/latest">

![Latest Release](https://img.shields.io/github/v/release/sheepla/qiitaz?style=flat-square)

</a>

</div>

<div align="center">

# qiitaz

</div>

<div align="center">

[Qiita](https://qiita.com)の記事を検索／プレビューできるコマンドラインツール

</div>

## 使い方

```
Usage:
  qiitaz [OPTIONS] QUERY...

Application Options:
  -V, --version  Show version
  -s, --sort=    Sort key to search e.g. "created", "like", "stock", "rel",
                 (default: "rel")
  -o, --open     Open URL in your web browser

Help Options:
  -h, --help     Show this help message
```

1. 引数に検索したいキーワードを指定してコマンドを実行します。
1. 検索結果をfuzzyfinderで絞り込みます。`Ctrl-N`, `Ctrl-P` または `Ctrl-J`, `Ctrl-K` でフォーカスを移動します。 `Tab`キーで選択し `Enter` キーで確定します。
1. 選択した記事のURLが出力されます。次のオプションを指定することで、選択した記事をブラウザで開いたりターミナル上でプレビューしたりすることができます。

### ブラウザで記事のページを開く

`-o`, `--open`オプションを付けるとデフォルトのブラウザが起動し、選択した記事のページが開きます。

### 記事をプレビューする

`-p`, `--preview` オプションを付けると、ターミナル上で記事をプレビューすることができます。
lessなどのページャを使うと読みやすいです。

**例**: `qiitaz -p QUERY... | less -R`

### 高度な検索

クエリ引数に次のオプションや演算子を指定することで、条件を詳細に指定して検索することができます。

**例**: `qiitaz title:Go created:\>2022-03-01`

|オプション          |説明                              |
|--------------------|----------------------------------|
|`title:{{string}}`  |タイトルにそのキーワードが含まれる|
|`body:{{string}}`   |本文にそのキーワードが含まれる    |
|`code:{{string}}`   |コードにそのキーワードが含まれる  |
|`tag:{{string}}`    |記事に付けられているタグ          |
|`-tag:{{string}}`   |除外するタグ                      |
|`user:{{string}}`   |ユーザー名                        |
|`stocks:\>{{int}}`  |ストック数                        |
|`created:\>{{date}}`|作成日がその日以降                |
|`updated:\>{{date}}`|更新日がその日以降                |


|演算子                |説明  |
|----------------------|------|
|`{{条件}} OR {{条件}}`|OR条件|

また、`-s`, `--sort` オプションを指定することで、ソート条件を変更することができます。

**例**: `qiitaz -s like Go`

|値       |説明              |
|---------|------------------|
|`rel`    |関連度順          |
|`like`   |LGTM数の多い順    |
|`stock`  |ストック数の多い順|
|`created`|作成日順          |

## インストール

リリースページから実行可能なバイナリをダウンロード可能です。

> [Latest Release](https://github.com/sheepla/qiitaz/releases/latest)

ソースからビルドする場合は、このリポジトリをクローンして `go install` を実行してください。
`v1.17.8 linux/amd64`で開発しています。

## 関連

- [sheepla/fzwiki](https://github.com/sheepla/fzwiki)
- [sheepla/fzenn](https://github.com/sheepla/fzenn)

