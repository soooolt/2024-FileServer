# ファイル整理ウェブサイトの作成

## 1. ファイルの種類

- **画像ファイル**: png, jpeg, jpg, gif, webp
- **動画ファイル**: mp4
- **音声ファイル**: mp3

> **注**: 文書やコードはGitHubがあるので扱う必要はありません。

## 2. ファイルの整理

整理とは、検索や閲覧が容易であることを指します。具体的には以下のような機能を提供します。

- **タグをつけて検索ができる**: 例えば、場所や日付などの情報をタグとして追加します。
- **ファイル拡張子で検索ができる**: ファイル拡張子もタグとして扱うことができます。
- **よく使うファイルとそうでないファイルを分ける**: 閲覧回数等を保持し、よく使うファイルを優先的に表示します。

## 3. ウェブサイトの構成

### 3.1 ファイル検索画面

| 機能 | 説明 |
| --- | --- |
| 検索フォーム | タグを用いたOR,AND検索を行う |
| ファイル追加フォーム | 新しいファイルの追加, タグ付け |
| タグ追加ボタン | タグを別途追加する |
| フィルタボタン | 検索フォーム以外に基礎的な項目をANDに用いる |

### 3.2 閲覧画面

| 機能 | 説明 |
| --- | --- |
| 見やすい画面 | できる限りファイルを大きめに表示 |
| タグ追加ボタン | タグを別途追加する |

## 4. システム構造

以下にシステムのファイル構成を示します。

```
.
├── contents
│   ├── image
│   ├── animated
│   ├── movie
│   ├── audio
│   └── FileInfo.db
├── page
│   ├── index.html
│   └── view.html
├── static
│   ├── css
│   │   ├── base.css
│   │   ├── index.css
│   │   └── view.css
│   └── js
│       └── script.js
├── main.go
└── library
    ├── dbutills
    │   └── dbutill.go
    └── fileops
        └── fileops.go
```

- **contents**: 各種メディアファイルとファイル情報データベースが格納されています。
- **page**: ウェブサイトのHTMLファイルが格納されています。
- **static**: CSSとJavaScriptファイルが格納されています。
- **main.go**: メインのGoファイルです。
- **library**: ユーティリティ関数が格納されているGoファイルが格納されています。

## 5.データベースの構成
![image](https://github.com/soooolt/2024-FileServer/assets/126924993/761e8788-1bd9-4e83-ba30-728bf6e78a80)
```erDiagram
erDiagram
    Media ||--|{ MediaTag : has
    Tag ||--|{ MediaTag : has
    MediaType ||--|| Media : "is a type of"

    Media {
        int ID
        string NAME
        int MediaType_ID
    }

    MediaType {
        int ID
        string TYPE
    }

    Tag {
        int ID
        string TAG
    }

    MediaTag {
        int Media_ID
        int Tag_ID
    }
```
