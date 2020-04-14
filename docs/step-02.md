# 冗長構成の EC2・RDS アプリケーションを構築する

## 参考

- [Word press preinstall-iam 対応版-aws 体験ハンズオン-セキュア＆スケーラブルウェブサービス構築編](https://www.slideshare.net/KamedaHarunobu/word-press-preinstalliamaws)
- [AWS 体験ハンズオン〜セキュア&スケーラブルウェブサービス構築〜](https://aws-ref.s3.amazonaws.com/handson/building_3tier_on_vpc_ver9.2_ja.pdf)
- [さくらの VPS に Go の環境を構築してみた](https://saitodev.co/article/%E3%81%95%E3%81%8F%E3%82%89%E3%81%AEVPS%E3%81%ABGo%E3%81%AE%E7%92%B0%E5%A2%83%E3%82%92%E6%A7%8B%E7%AF%89%E3%81%97%E3%81%A6%E3%81%BF%E3%81%9F)
- [SCP を使ってローカルから AWS EC2 サーバーへファイルをアップロードするメモ](https://qiita.com/Buddychen/items/51c23b105497f48e3da5)

## VPC の作成（サブネットも同時に一つ作成）

- リージョンを東京にする
- VPC > VPC ウィザードの起動 > 1 個のパブリックサブネットを持つ VPC
- VPC 名を入力
- AZ を選択
- VPC を作成
- VPC とサブネットを確認

## サブネットの追加

- 名前タグは「パブリックサブネット」「プライベートサブネット」のどちらか

| No. | 名前タグ               | VPC             | AZ              | CIDR ブロック |
| :-- | :--------------------- | :-------------- | :-------------- | :------------ |
| 1   | パブリックサブネット   | blog-production | ap-northeast-1a | 10.0.0.0/24   |
| 2   | パブリックサブネット   | blog-production | ap-northeast-1c | 10.0.1.0/24   |
| 3   | プライベートサブネット | blog-production | ap-northeast-1a | 10.0.2.0/24   |
| 4   | プライベートサブネット | blog-production | ap-northeast-1c | 10.0.3.0/24   |

- パブリックサブネットのルートテーブルを編集
- ルートテーブルの選択肢が 2 つあるので、インターネットゲートウェイ(0.0.0.0/0)がある方を選択

## EC2 インスタンス作成

- EC2 > インスタンス > インスタンスの作成 > Amazon Linux 2 (x86) > t2.micro > インスタンスの詳細の設定
  - ネットワーク: blog-production
  - サブネット: パブリックサブネット 1
  - 自動割り当てパブリック IP: 有効
  - ストレージの追加 > タグ付け > タグを追加 Name; WebServer #1
  - セキュリティグループの設定 > 新しいセキュリティグループを作成する > セキュリティグループ名: Web > 説明 > Web > ルールの追加 > タイプ: HTTP > ソース: 任意の場所 > 確認と作成 > 作成
  - 新しいキーペアを作成 WebServer#1 キーペアのダウンロード
  - インスタンスの作成
  - [Linux インスタンスへの接続方法](https://docs.aws.amazon.com/console/ec2/instances/connect/docs)

## Elastic IP (EIP) の取得

- インスタンス削除後もまた同じ IP を利用できるようにする
- EC2 > ネットワーク&セキュリティ > Elastic IP > Elastic IP アドレスの割り当て > Amazon の IPv4 アドレスプール > 割り当て
- 作成された IP をクリック > Elastic IP アドレスの関連付け > インスタンス: WebServer#1 > 関連づける > 作成されたパブリック IP アドレスをメモする

## EC2 に SSH 接続

```sh
chmod 600 ~/Downloads/WebServer1.pem
ssh -i ~/Downloads/WebServer1.pem ec2-user@3.114.176.52
```

```sh
amazon-linux-extras list
amazon-linux-extras list | grep go
wget https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz
echo $PATH
sudo tar -C /usr/local -xzf go1.14.2.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
exec $SHELL -l
go version
```

```sh
pwd
mkdir ~/webapp
```

```sh
scp -i ~/Downloads/WebServer1.pem ./playground/main.go ec2-user@$WEB_SERVER_1_IP:~/
```

```sh
go run main.go
```

```sh
2020/04/12 14:02:22 listen tcp :80: bind: permission denied
```

調べたところポート番号 1024 以下はルート権限で実行する必要があるとのこと。

```sh
go build -o app
sudo ./app
```

```sh
curl http://パブリックIP/hello
```

## プライベートサブネットにオーロラを配置してみた

### DB 用セキュリティグループの作成

- EC2 > ネットワーク&セキュリティ > セキュリティグループ > セキュリティグループの作成
  - セキュリティグループ名: DB
  - 説明: DB For DB For Aurora
  - VPC: blog-production
- インバウンド > ルールの追加
  - タイプ: MYSQL/Aurora
  - 送信元: カスタム IP
    - Web セキュリティグループ
- セキュリティグループの作成

## DB サブネットグループの作成

- RDS > サブネットグループ > DB サブネットグループの作成
  - 名前: db subnet
  - 説明: RDS for MySQL/Aurora
  - VPC: blog-production
  - サブネットの追加
    - 1
      - ap-northeast-1a
      - 10.0.2.0/24
      - サブネットを追加します
    - c
      - ap-northeast-1b
      - 10.0.3.0/24
      - サブネットを追加します
  - 作成

## DB インスタンスの作成

- RDS > データベース > データベースの作成
- データベース作成方法: 標準作成 > エンジンのタイプ: Amazon Aurora > エディション: MySQL との互換性を持つ Amazon Aurora > バージョン: MySQL 5.7 2.07.2 > テンプレート: 開発/テスト
- DB クラスター識別子: blog-database
- マスターユーザー名: admin
- パスワードの自動生成: true
- DB インスタンスクラス: バースト可能クラス db.t2.small
- マルチ AZ 配置: Aurora レプリカを作成しない
- Virtual Private Cloud (VPC): blog-production
- 追加の接続設定
- サブネットグループ: db subnet
- パブリックアクセス可能: なし
- VPC セキュリティグループ: 既存の選択
- 既存の VPC セキュリティグループ: DB
- アベイラビリティゾーン: ap-northeast-1a
- データベースポート: 3306
- 追加設定
- DB インスタンス識別子: blog-database-instance-1
- 最初のデータベース名: blog
- DB クラスターのパラメータグループ: 初期値
- DB パラメータグループ: 初期値
- フェイルオーバー優先順位: 指定なし
- バックアップ保存期間: 1 日間
- スナップショットにタグをコピー: true
- 暗号を有効化: false
- バックトラックを有効にする: false
- 拡張モニタリングの有効化: true
- 詳細度: 60 秒
- モニタリングロール: default
- マイナーバージョン自動アップグレードの有効化: true
- メンテナンスウィンドウ: 設定なし
- 削除保護の有効化: false
- データベースの作成
- 認証情報の詳細
  - user
  - pw
  - endpoint
  - をメモする

## RDS に SSH ポートフォワーディングを利用して接続してみた

- [AWS クラウドでの Linux 踏み台ホスト: クイックスタートリファレンスデプロイ](https://docs.aws.amazon.com/ja_jp/quickstart/latest/linux-bastion/welcome.html)
- [【AWS VPC 入門】4.NatGateway/Bastion](https://qiita.com/_Yasuun_/items/4fc50d94baafd6e38af2)

### EC2 を作成する

- 普通に EC2 インスタンスを作成する
- パブリックサブネットに配置
- 踏み台用のセキュリティグループを作る
- ソースをマイ IP にする
- キーペアを新規作成する
- DB セキュリティグループのインバウンドルールに Step セキュリティグループからのアクセス許可を追加する

## EC2 から DB に接続する

### データベースにデータを投入する

```

```

## CodeDeploy で EC2 にデプロイしてみた

- [チュートリアル: WordPress を Amazon EC2 インスタンス (Amazon Linux または Red Hat Enterprise Linux および Linux, macOS, or Unix) にデプロイする](https://docs.aws.amazon.com/ja_jp/codedeploy/latest/userguide/tutorials-wordpress.html)

EC2 インスタンスに接続する

```sh
ssh -i ~/Downloads/WebServer1.pem ec2-user@$$WEB_SERVER_1_IP
```

CodeDeploy エージェントがインストールされているか確認する

```sh
sudo service codedeploy-agent status
```
