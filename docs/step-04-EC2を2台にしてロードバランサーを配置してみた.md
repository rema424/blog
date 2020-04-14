# EC2 を 2 台にしてロードバランサーを配置してみた

## Web サーバーの AMI (Amazon Mashine Image) を作成

- EC2 > インスタンス > WebSserver#1 右クリック > イメージ > イメージの作成
- イメージ名: Blog WebServer
- イメージの作成
- EC2 > AMI > 確認

## AMI から 2 個目の WebServer EC2 インスタンスを作成する

- EC2 > AMI > Blog WebServer 右クリック > 起動
- t2.micro > 次のステップ
- 情報を入力して次のステップ
  - ネットワーク: blog-production vpc
  - サブネット: パブリックサブネット | ap-northeast-1c
  - 自動割り当てパブリック IP: 有効
- ストレージは変更せず次のステップ
- タグの追加 > キー: Name, 値: WebServer #2 > 次のステップ
- 既存のセキュリティグループを選択する > Web > 確認と作成
- 起動
- 既存のキーペアの選択 > WebServer#1 > チェックボックスにチェック > 作成
- インスタンスの表示

## ELB を作成

- EC2 > ロードバランサー > ロードバランサーの作成 > Application Load Balancer > 作成
- 手順 1: ロードバランサーの設定
  - 名前: blog-lb
  - スキーム: インターネット向け
  - IP アドレスタイプ: ipv4
  - アベイラビリティゾーン
    - VPC: blog-production
    - アベイラビリティゾーン
      - ap-northeast-1a: パブリックサブネット
      - ap-northeast-1c: パブリックサブネット
  - 次の手順
- 手順 2: セキュリティ設定の構成
  - 次の手順
- 手順 3: セキュリティグループの設定
  - セキュリティグループの割り当て: 新しいセキュリティグループを作成する
  - セキュリティグループ名: LB
  - 説明: LB
  - タイプ: HTTP
  - 次の手順
- 手順 4: ルーティングの設定
  - ターゲットグループ
    - ターゲットグループ: 新しいターゲットグループ
    - 名前: blog-tg
    - ターゲットの種類: インスタンス
    - プロトコル: HTTP
    - ポート: 80
  - ヘルスチェック
    - プロトコル: HTTP
    - パス: /
  - 次の手順
- 手順 5: ターゲットの登録
  - WebServer #1, WebServer #2 > 登録済みに追加
  - 次の手順
- 手順 6: 確認
  - 作成
- ロードバランサー作成状況
  - 閉じる

## 接続確認

- EC2 > ロードバランサー > DNS 名 > ブラウザでアクセス
- 502 Bad Gateway
- ヘルスチェックに合格しないとターゲットが設定されないらしい（"/"は登録されていない）
- ヘルスチェックのパスを"/people"に変更
- 再度アクセス > 502
- EC2 が稼働していないらしい

```sh
make ssh-web-1
export PATH=$PATH:/usr/local/go/bin
export DB_USER=developer
export DB_PASS=Passw0rd!
export DB_NAME=blog
export DB_HOST=blog-database.cluster-cll9xfuraffh.ap-northeast-1.rds.amazonaws.com
export DB_PORT=3306
export DB_NET=tcp
export DB_ADDR=''
sudo -E ./app &
pgrep -alf app
exit
```

## Web セキュリティグループの設定変更

- EC2 > セキュリティグループ > Web > インバウンドルール > インバウンドルールの編集
  - 既存の設定を全て削除
  - タイプ: HTTP, ソース: カスタム-セキュリティグループ-LB
  - タイプ: HTTP, ソース: マイ IP
  - タイプ: SSH, ソース: 任意の場所
- ルールの保存

```sh
make ssh-web-2
export PATH=$PATH:/usr/local/go/bin
export DB_USER=developer
export DB_PASS=Passw0rd!
export DB_NAME=blog
export DB_HOST=blog-database.cluster-cll9xfuraffh.ap-northeast-1.rds.amazonaws.com
export DB_PORT=3306
export DB_NET=tcp
export DB_ADDR=''
sudo -E ./app &
pgrep -alf app
pkill
```

再度 `http://ロードバランサーのDNS名/people` にアクセス > 成功
