# 郵便番号データ

日本郵便が提供する郵便番号データ（CSV形式）をJSON形式に変換するツールです。

## 特徴

- 日本郵便の公式サイトから最新の郵便番号CSVデータを取得し、JSONファイルを生成します。
- Go言語で実装されており、シンプルで高速に動作します。
- 生成されたJSONデータは、郵便番号検索や住所変換などのアプリケーションで利用可能です。

## 使用方法

### 前提条件

- Go 1.16以降がインストールされていること。

### 手順

1. リポジトリをクローンします。

    ```bash
    git clone https://github.com/zafinc/zipcode.git
    cd zipcode
    ```

2. `main.go` を実行してJSONファイルを生成します。

    ```bash
    go run main.go
    ```

   実行後、`data` ディレクトリ内に `{郵便番号の上3桁}.json` が生成されます。

## データソース

- 日本郵便 公式サイト: [https://www.post.japanpost.jp/](https://www.post.japanpost.jp/)

本ツールは、日本郵便が提供するCSV形式の郵便番号データを利用しています。

## 利用方法

* Typescript例

    ```typescript
    async function setZipcode(zipcode: string): void {
      if (zipcode.length != 7) { return }
      const response = await fetch('https://zafinc.github.io/zipcode/data/' + zipcode.substring(0, 3) + '.json')
      if (response.ok) {
        const addresses = await response.json()
        const address = addresses[zipcode]
        if (address !== undefined) {
          // 都道府県コード
          console.log(address.prefectureCode)
          // 市区町村
          console.log(address.city)
          // 町域
          console.log(address.town)
        }
      }
    }
    ```
