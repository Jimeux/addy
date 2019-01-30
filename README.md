# Addy

## 実装の流れ
1. https://checkout-test.adyen.com/v40/paymentSession にリクエストを送り、`paymentSession`のデータをビューに入れ込む。
2. JSで`paymentSession`のデータを取得しSDKの`chckt.checkout`に渡し、フォームを初期化する。
3. クレジットカードの情報を入れ、サブミットする
4. 3が成功した場合`recurringDetailReference`がサーバに投げられる
5. 4の値で`/payments`エンドポイントを叩く

## Let's 実行

1. `ADYEN_API_KEY`と`ADYEN_MERCHANT_ACCOUNT`を環境変数に入れるか、`main.go`に直接入れておく。
2. `GO111MODULE=on go run main.go`

## Let's Test
1. http://localhost:8080/test にアクセス。
2. [Test card numbers](https://docs.adyen.com/developers/test-cards/test-card-numbers)のいずれを使ってみる。

## 参考
- [AdyenからのPHPサンプルレポジトリ](https://github.com/Adyen/adyen-web-sdk-sample-code)
- [Web SDKのドキュメント](https://docs.adyen.com/developers/checkout/web-sdk)
