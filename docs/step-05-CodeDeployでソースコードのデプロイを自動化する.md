## EC2 インスタンスに CodeDeploy エージェントをインストールする

[Amazon Linux または RHEL 用の CodeDeploy エージェントのインストールまたは再インストール](https://docs.aws.amazon.com/ja_jp/codedeploy/latest/userguide/codedeploy-agent-operations-install-linux.html)
[AWS CodeDeploy を使って Laravel アプリをデプロイしてみた](https://qiita.com/kurikazu/items/869e7f4265c073d52211)

マシンに SSH 接続して次のコマンドを実行する

```sh
sudo yum update
sudo yum install ruby
sudo yum install wget
cd /home/ec2-user
wget https://aws-codedeploy-ap-northeast-1.s3.ap-northeast-1.amazonaws.com/latest/install
chmod +x ./install
sudo ./install auto
sudo service codedeploy-agent status
```

##

[CodeDeploy と連動するように Amazon EC2 インスタンスを設定する](https://docs.aws.amazon.com/ja_jp/codedeploy/latest/userguide/instances-ec2-configure.html)
[CodeDeploy を Amazon EC2 Auto Scaling と統合する](https://docs.aws.amazon.com/ja_jp/codedeploy/latest/userguide/integrations-aws-auto-scaling.html)

## CodeDeploy にアタッチする IAM ロールを作成する

[サービスロールの作成 (コンソール)](https://docs.aws.amazon.com/ja_jp/codedeploy/latest/userguide/getting-started-create-service-role.html#getting-started-create-service-role-console)

- IAM > ロール > ロールの作成 > AWS サービス > CodeDeploy > CodeDeploy > 次の手順
- 次のステップ
- 次のステップ
- ロール名: CodeDeployServiceRole
- ロールの作成
- CodeDeployServiceRole > インラインポリシーの追加
  - サービス: EC2
  - アクション: RunInstances
  - リソース: 全てのリソース
  - サービス: EC2
  - アクション: CreateTags
  - リソース: 全てのリソース
  - サービス: IAM
  - アクション: PassRole
  - リソース: 全てのリソース
- CodeDeployAdditionalPolicy > ポリシーの作成

## EC2 にアタッチする IAM ロールを作成する

- IAM > ポリシー > ポリシーの作成 > JSON

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:Get*", "s3:List*"],
      "Resource": ["arn:aws:s3:::aws-codedeploy-ap-northeast-1/*"]
    }
  ]
}
```

- ポリシーの確認 > 名前: CodeDeploy-EC2-Permissions > ポリシーの作成
- IAM > ロール > ロールの作成 > AWS サービス > EC2 > 次の手順
- CodeDeploy-EC2-Permissions > 次のステップ
- 次のステップ
- ロール名: CodeDeploy-EC2-Instance-Profile
- ロールの作成

## Amazon EC2 Auto Scaling グループを作成する

[Amazon EC2 Auto Scaling グループを作成して設定するには (コンソール)](https://docs.aws.amazon.com/ja_jp/codedeploy/latest/userguide/tutorials-auto-scaling-group-create-auto-scaling-group.html)

- EC2 > ATUO SCALING > 起動設定 > 起動設定の作成
- AMI の選択
  - マイ AMI > Blog WebServer
- インスタンスタイプの選択
  - t2.micro > 次の手順
- 詳細設定
  - 名前: CodeDeploy-AS-Configuration
  - IAM ロール: CodeDeploy-EC2-Instance-Profile
- ストレージの追加
  - 次の手順
- セキュリティグループの設定
  - セキュリティグループの割り当て: 既存のセキュリティグループを選択する > Web > 確認
- 確認
  - 起動設定の作成
- 既存のキーペアを選択するか、新しいキーペアを作成します。
  - 既存のキーペアを選択
  - WebServer#1
  - チェックボックス: ON
  - 起動設定の作成
- 起動設定の作成ステータス
  - この起動設定を使用して Auto Scaling グループを作成する
- Auto Scaling グループの作成
  - グループ名: CodeDeploy-AS-Group
  - ネットワーク: blog-production
  - サブネット: 10.0.0.0/24, 10.0.1.0/24
  - 高度な詳細
    - ターゲットグループ: blog-tg
- スケーリングポリシーの設定
  - このグループを初期のサイズに維持する
- 通知の設定
- タグを設定
- 確認

## デプロイする

[デプロイを作成するには (コンソール)](https://docs.aws.amazon.com/ja_jp/codedeploy/latest/userguide/tutorials-auto-scaling-group-create-deployment.html#tutorials-auto-scaling-group-create-deployment-console)

- CodeDeploy > デプロイ > アプリケーション > アプリケーションの作成
  - アプリケーション名: blog-auto-scaling
  - コンピューティングプラットフォーム: EC2/オンプレミス
- デプロイグループ > デプロイグループの作成
  - デプロイグループ名 blog-auto-scaling-deploy-group
  - サービスロール: CodeDeployServiceRole
  - デプロイタイプ: インプレース
  - 環境設定: Amazon EC2 Auto Scaling グループ > CodeDeploy-AS-Group
  - デプロイ設定: CodeDeployDefault.OneAtATimes
  - ロードバランサー
    - ロードバランシングを有効にする: ON
    - Application Load Balancer
    - ターゲットグループの選択: blog-tg
  - デプロイグループの作成
