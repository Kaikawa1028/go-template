# canly-custom-db-admin-backend

Canly の カスタムDB&API サーバのソースコードです

## ローカル環境起動

Canly 本体の DB に接続するため、あらかじめ Canly 本体(meocloud リポジトリ)のアプリケーションを起動しておく必要があります。  
以下、詳しい手順になります。

以下のコマンドで Shared Network を作成します

```
docker network create canly-shared-network
```

### ここからは meocloud リポジトリで作業してください。

※`make up` を使用せずに Docker を立ち上げた場合、上記の Shared Network が meocloud リポジトリ側で設定されません。必ず `make up` にて meocloud（Canly 本体）の環境を立ち上げてください。
具体的には、以下のように db コンテナの network 設定か異なっています。

- docker-compose.yml

```
db:
（省略）
    networks:
      - app-network
```

- docker-compose.local.yml

```
db:
（省略）
    networks:
      - app-network
      - canly-shared-network
```

`cms_dev`ブランチに切り替えます。（API 用の migration が存在するため、ブランチ変更＋ migrate 実行は必ず行ってください。）

以下コマンドを実行し環境を起動します。

```
make up
```

以下のコマンドを実行し、マイグレーションを行います（DB のスキーマが更新されます）

```
# シェル起動
make app-bash
# マイグレーション実行
php artisan migrate
```

### ここからは本リポジトリ(canly-api-server)で作業してください。  
`web`フォルダで以下のコマンドを実行し、`.env`を作成します。

```
cp .env.example .env
cp .env.integration.example .env.integration
```

以下のコマンドを実行し、コンテナを起動します。

```
docker-compose up -d
```

以下のコマンドでコードを生成します

```
docker-compose exec web make generate
```

以下のコマンドで動作確認できます。  
200 OK が返ってくれば成功です。

```
curl -XGET -I 'http://localhost:3005/v1/healthcheck'
```

## インテグレーションテストの実行方法

Canly 本体の DB に接続するため、あらかじめ Canly 本体(meocloud リポジトリ)のアプリケーションを起動しておく必要があります。  
以下、詳しい手順になります。

meocloud リポジトリの MySQL に接続し、`meo_api_testing`という名前でデータベースを作成します

meocloud リポジトリの`cms/v1.2`ブランチで以下のコマンドを実行します

```
make app-bash
# スキーマ構築
DB_DATABASE=meo_api_testing php artisan migrate
# テストデータ追加
DB_DATABASE=meo_api_testing php artisan db:seed --class=ApiIntegrationTestSeeder
```

※DB をまるごと作り直したい場合は `migrate` の部分を `migrate:fresh` とすると DB を削除＆作成した後に migration を実行します

本リポジトリ(canly-api-server)で以下のコマンドを実行します

```
docker-compose up
docker-compose exec web make integration_tests
```

### VSCodeで開発する方法

[コチラ](https://www.notion.so/canlyhp/VSCode-503f644ce3114e0e97515eab085b38b3) を参照下さい

### カバレッジの確認方法

以下のコマンドを実行します

```
docker-compose up -d
docker-compose exec web make coverage
```

生成された web/coverage.html を開くとファイル毎の詳細なカバレッジを確認できます

また、以下のコマンドで関数毎のカバレッジが確認できます（標準出力に出力されます）

```
docker-compose exec web make coverage_func
```

## master を最新化した際に確認すべき事

- `web/.env.example` が更新されている場合は、手元の `web/.env` へ反映させる（`.env.integration`も同様）


## コミットルール

- Backlogのチケット番号( CANLY-XXXX ) から開始するようにしてください

[コチラ](https://github.com/marketplace/actions/backlog-notify)を参考にしてください。

Ex.)
```
CANLY-1234 : 店舗IDのバリデーションを修正しました
```


### コミットテンプレートの適用方法
```
git config --local commit.template .gitmessage
```

