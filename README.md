# Addy

## Let's 実行

1. `ADYEN_API_KEY`と`ADYEN_MERCHANT_ACCOUNT`を環境変数に入れるか、`main.go`に直接入れておく。
2. `GO111MODULE=on go run main.go`

## Let's Test
1. http://localhost:8080/ にアクセス。
2. [Test card numbers](https://docs.adyen.com/developers/test-cards/test-card-numbers)のいずれを使ってみる。

## 問題
- 以下のエラーがJSコンソールにどうしても出て、フォームをサブミットできぬ。
```
WARNING: securedFields:: the encryption algorithm is not present. It will not be possible to encrypt input fields
```
- JSインポートの問題に見えるが、問題はアカウントにある可能性もありそう。
[“the encryption algorithm is not present. It will not be possible to encrypt input fields” - Google 検索](https://www.google.co.jp/search?q=%22the+encryption+algorithm+is+not+present.+It+will+not+be+possible+to+encrypt+input+fields%22&oq=%22the+encryption+algorithm+is+not+present.+It+will+not+be+possible+to+encrypt+input+fields%22&aqs=chrome..69i57.948j0j7&sourceid=chrome&ie=UTF-8)
- Adyenのカスタマサポートに連絡した方が速いかも。