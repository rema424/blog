# 冗長構成の EC2・RDS アプリケーションを構築する

## 参考

- [Word press preinstall-iam 対応版-aws 体験ハンズオン-セキュア＆スケーラブルウェブサービス構築編](https://www.slideshare.net/KamedaHarunobu/word-press-preinstalliamaws)
- [AWS 体験ハンズオン〜セキュア&スケーラブルウェブサービス構築〜](https://aws-ref.s3.amazonaws.com/handson/building_3tier_on_vpc_ver9.2_ja.pdf)

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
