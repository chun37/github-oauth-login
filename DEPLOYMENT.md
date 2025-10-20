# サーバーデプロイガイド (Docker不使用)

このドキュメントでは、Docker/Podmanを使用せずに、サーバー上に直接アプリケーションをデプロイする方法を説明します。

## 目次

1. [前提条件](#前提条件)
2. [PostgreSQLのセットアップ](#postgresqlのセットアップ)
3. [バックエンドのデプロイ](#バックエンドのデプロイ)
4. [フロントエンドのデプロイ](#フロントエンドのデプロイ)
5. [nginxの設定](#nginxの設定)
6. [動作確認](#動作確認)
7. [トラブルシューティング](#トラブルシューティング)

## 前提条件

以下のソフトウェアがサーバーにインストールされている必要があります：

- **Go**: 1.25以上
- **Node.js**: 20以上
- **pnpm**: 最新版
- **PostgreSQL**: 16以上
- **nginx**: 最新版

### ソフトウェアのインストール

#### Go 1.25以上

```bash
# 最新のGoをダウンロード（バージョンは適宜変更）
wget https://go.dev/dl/go1.25.0.linux-amd64.tar.gz

# 既存のGoを削除（ある場合）
sudo rm -rf /usr/local/go

# 解凍してインストール
sudo tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz

# パスを設定（~/.bashrcまたは~/.zshrcに追加）
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# バージョン確認
go version
```

#### Node.js 20とpnpm

```bash
# Node.js 20のインストール（nvmを使用）
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install 20
nvm use 20

# pnpmのインストール
corepack enable
corepack prepare pnpm@latest --activate

# バージョン確認
node --version
pnpm --version
```

#### PostgreSQL 16

```bash
# Ubuntu/Debianの場合
sudo apt update
sudo apt install postgresql-16 postgresql-client-16

# CentOS/RHELの場合
sudo dnf install postgresql16-server postgresql16
sudo postgresql-setup --initdb
sudo systemctl enable postgresql
sudo systemctl start postgresql

# バージョン確認
psql --version
```

#### nginx

```bash
# Ubuntu/Debianの場合
sudo apt install nginx

# CentOS/RHELの場合
sudo dnf install nginx

# サービスの有効化
sudo systemctl enable nginx
```

## PostgreSQLのセットアップ

### 1. データベースとユーザーの作成

```bash
# PostgreSQLユーザーに切り替え
sudo -u postgres psql

# PostgreSQL内で以下のSQLを実行
CREATE DATABASE github_oauth_app;
CREATE USER postgres WITH PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON DATABASE github_oauth_app TO postgres;

# PostgreSQL 15以降の場合、スキーマ権限も付与
\c github_oauth_app
GRANT ALL ON SCHEMA public TO postgres;

# 終了
\q
```

### 2. PostgreSQLの接続設定

PostgreSQLの設定ファイルを編集して、ローカル接続を許可します。

```bash
# pg_hba.confファイルを編集
sudo nano /etc/postgresql/16/main/pg_hba.conf

# 以下の行を追加または変更
# local   all             all                                     md5
# host    all             all             127.0.0.1/32            md5
```

PostgreSQLを再起動：

```bash
sudo systemctl restart postgresql
```

### 3. マイグレーションの実行

プロジェクトのマイグレーションファイルを使用してテーブルを作成します。

```bash
# プロジェクトディレクトリに移動
cd /path/to/github-oauth-login

# マイグレーションを実行
psql -U postgres -d github_oauth_app < backend/migrations/001_create_sessions_table.up.sql
```

## バックエンドのデプロイ

### 1. プロジェクトのクローン

```bash
# デプロイ先ディレクトリに移動（例: /opt）
cd /opt

# Gitリポジトリをクローン
git clone <your-repository-url> github-oauth-login
cd github-oauth-login
```

### 2. 環境変数の設定

```bash
# .envファイルを作成
cp .env.example .env
nano .env
```

`.env`ファイルを編集：

```env
# GitHub OAuth Settings
GITHUB_CLIENT_ID=your_actual_github_client_id
GITHUB_CLIENT_SECRET=your_actual_github_client_secret
GITHUB_REDIRECT_URL=https://yourdomain.com/api/auth/callback

# Database Settings
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=github_oauth_app
DB_SSLMODE=disable

# Session Settings
SESSION_SECRET=your_random_session_secret_min_32_characters_or_more

# Cookie Settings
COOKIE_DOMAIN=10.11.22.112

# Application Settings
BACKEND_PORT=8080
FRONTEND_URL=https://yourdomain.com
BACKEND_URL=https://yourdomain.com

# Environment
ENV=production
```

### 3. バックエンドのビルド

```bash
cd backend

# 依存関係のダウンロード
go mod download

# バイナリのビルド
CGO_ENABLED=0 GOOS=linux go build -o github-oauth-backend ./cmd/api

# バイナリを適切な場所に配置
sudo mkdir -p /opt/github-oauth-app
sudo cp github-oauth-backend /opt/github-oauth-app/
sudo cp -r migrations /opt/github-oauth-app/
```

### 4. systemdサービスの作成（推奨）

バックエンドをsystemdサービスとして設定します。

```bash
sudo nano /etc/systemd/system/github-oauth-backend.service
```

以下の内容を記述：

```ini
[Unit]
Description=GitHub OAuth Backend Service
After=network.target postgresql.service
Requires=postgresql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/github-oauth-app
EnvironmentFile=/opt/github-oauth-login/.env
ExecStart=/opt/github-oauth-app/github-oauth-backend
Restart=on-failure
RestartSec=5s

# セキュリティ設定
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

サービスを有効化して起動：

```bash
# サービスの有効化
sudo systemctl enable github-oauth-backend.service

# サービスの起動
sudo systemctl start github-oauth-backend.service

# ステータスの確認
sudo systemctl status github-oauth-backend.service

# ログの確認
sudo journalctl -u github-oauth-backend.service -f
```

## フロントエンドのデプロイ

### 1. フロントエンドのビルド

```bash
cd /opt/github-oauth-login/frontend

# 依存関係のインストール
pnpm install --frozen-lockfile

# 環境変数を設定してビルド
NEXT_PUBLIC_BACKEND_URL=/api pnpm build
```

### 2. ビルド成果物の配置

Next.jsのstandaloneモードでビルドされたファイルを配置します。

```bash
# ビルド成果物を適切な場所に配置
sudo mkdir -p /opt/github-oauth-app/frontend
sudo cp -r .next/standalone/* /opt/github-oauth-app/frontend/
sudo cp -r .next/static /opt/github-oauth-app/frontend/.next/
sudo cp -r public /opt/github-oauth-app/frontend/
```

### 3. systemdサービスの作成（推奨）

フロントエンドをsystemdサービスとして設定します。

```bash
sudo nano /etc/systemd/system/github-oauth-frontend.service
```

以下の内容を記述：

```ini
[Unit]
Description=GitHub OAuth Frontend Service
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/github-oauth-app/frontend
Environment="NODE_ENV=production"
Environment="PORT=3000"
ExecStart=/usr/bin/node server.js
Restart=on-failure
RestartSec=5s

# セキュリティ設定
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

サービスを有効化して起動：

```bash
# サービスの有効化
sudo systemctl enable github-oauth-frontend.service

# サービスの起動
sudo systemctl start github-oauth-frontend.service

# ステータスの確認
sudo systemctl status github-oauth-frontend.service

# ログの確認
sudo journalctl -u github-oauth-frontend.service -f
```

## nginxの設定

### 重要: Cookie設定について

バックエンドが正しくCookieを設定できるよう、`.env`ファイルに`COOKIE_DOMAIN`を必ず設定してください。

- **開発環境**: `COOKIE_DOMAIN=127.0.0.1`
- **本番環境**: `COOKIE_DOMAIN=10.11.22.112` (実際のサーバーのIPアドレスまたはドメイン名)

この設定がないと、GitHubからのOAuthリダイレクト後にCookieが送信されず、認証が失敗します。

### 1. nginx設定ファイルの作成

```bash
sudo nano /etc/nginx/sites-available/github-oauth-app
```

以下の内容を記述：

```nginx
upstream frontend {
    server 127.0.0.1:3000;
}

upstream backend {
    server 127.0.0.1:8080;
}

server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;

    # HTTPS にリダイレクト（SSL証明書設定後に有効化）
    # return 301 https://$server_name$request_uri;
}

server {
    # HTTPSを使用する場合（SSL証明書を設定後に有効化）
    # listen 443 ssl http2;
    # server_name yourdomain.com www.yourdomain.com;

    # SSL証明書の設定（Let's Encryptなどで取得）
    # ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    # ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    # ssl_protocols TLSv1.2 TLSv1.3;
    # ssl_ciphers HIGH:!aNULL:!MD5;

    # 開発/テスト環境の場合はHTTPのまま使用
    listen 80;
    server_name yourdomain.com www.yourdomain.com;

    # フロントエンド
    location / {
        proxy_pass http://frontend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Cookie の転送を有効化
        proxy_pass_request_headers on;
        proxy_set_header Cookie $http_cookie;
        proxy_pass_header Set-Cookie;
    }

    # バックエンドAPI
    location /api/ {
        proxy_pass http://backend/api/;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Cookie の転送を有効化
        proxy_pass_request_headers on;
        proxy_set_header Cookie $http_cookie;
        proxy_pass_header Set-Cookie;
    }
}
```

### 2. nginx設定の有効化

```bash
# シンボリックリンクを作成
sudo ln -s /etc/nginx/sites-available/github-oauth-app /etc/nginx/sites-enabled/

# デフォルト設定を無効化（必要に応じて）
sudo rm /etc/nginx/sites-enabled/default

# 設定ファイルのテスト
sudo nginx -t

# nginxを再起動
sudo systemctl restart nginx
```

### 3. ファイアウォールの設定

```bash
# HTTP/HTTPSポートを開放
sudo ufw allow 'Nginx Full'

# または個別に
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
```

### 4. SSL証明書の取得（本番環境推奨）

Let's Encryptを使用して無料のSSL証明書を取得します。

```bash
# Certbotのインストール
sudo apt install certbot python3-certbot-nginx

# SSL証明書の取得
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# 自動更新の設定（cron）
sudo certbot renew --dry-run
```

証明書取得後、nginx設定ファイルのHTTPS部分のコメントを解除してください。

## 動作確認

### 1. サービスの起動確認

```bash
# バックエンドの確認
sudo systemctl status github-oauth-backend.service
curl http://127.0.0.1:8080/api/auth/check

# フロントエンドの確認
sudo systemctl status github-oauth-frontend.service
curl http://127.0.0.1:3000

# nginxの確認
sudo systemctl status nginx
curl http://127.0.0.1
```

### 2. アプリケーションへのアクセス

ブラウザで以下のURLにアクセス：

- HTTP: `http://yourdomain.com`
- HTTPS: `https://yourdomain.com`（SSL証明書設定後）

### 3. GitHub OAuth設定の確認

GitHubのOAuth App設定で、以下を確認：

- **Homepage URL**: `https://yourdomain.com`
- **Authorization callback URL**: `https://yourdomain.com/api/auth/callback`

## トラブルシューティング

### バックエンドが起動しない

```bash
# ログを確認
sudo journalctl -u github-oauth-backend.service -n 50

# PostgreSQLの接続を確認
psql -U postgres -d github_oauth_app -h 127.0.0.1

# 環境変数を確認
sudo systemctl show github-oauth-backend.service --property=Environment
```

### フロントエンドが起動しない

```bash
# ログを確認
sudo journalctl -u github-oauth-frontend.service -n 50

# Node.jsのバージョンを確認
node --version

# ビルド成果物を確認
ls -la /opt/github-oauth-app/frontend/
```

### nginxでエラーが発生する

```bash
# nginxのエラーログを確認
sudo tail -f /var/log/nginx/error.log

# アクセスログを確認
sudo tail -f /var/log/nginx/access.log

# 設定ファイルをテスト
sudo nginx -t
```

### データベース接続エラー

```bash
# PostgreSQLの状態を確認
sudo systemctl status postgresql

# 接続設定を確認
cat /etc/postgresql/16/main/pg_hba.conf

# PostgreSQLのログを確認
sudo tail -f /var/log/postgresql/postgresql-16-main.log
```

## アップデート手順

アプリケーションを更新する場合：

```bash
# 1. リポジトリを更新
cd /opt/github-oauth-login
git pull

# 2. バックエンドを再ビルド
cd backend
go build -o github-oauth-backend ./cmd/api
sudo cp github-oauth-backend /opt/github-oauth-app/

# 3. フロントエンドを再ビルド
cd ../frontend
pnpm install --frozen-lockfile
NEXT_PUBLIC_BACKEND_URL=/api pnpm build
sudo cp -r .next/standalone/* /opt/github-oauth-app/frontend/
sudo cp -r .next/static /opt/github-oauth-app/frontend/.next/
sudo cp -r public /opt/github-oauth-app/frontend/

# 4. サービスを再起動
sudo systemctl restart github-oauth-backend.service
sudo systemctl restart github-oauth-frontend.service
```

## セキュリティのベストプラクティス

1. **ファイアウォールの設定**: 必要なポート（80、443）のみ開放
2. **SSL/TLSの使用**: 本番環境では必ずHTTPSを使用
3. **定期的なアップデート**: OSとパッケージを常に最新に保つ
4. **強力なパスワード**: データベースとセッションシークレットに強力なパスワードを使用
5. **最小権限の原則**: サービスは専用ユーザー（www-data）で実行
6. **ログの監視**: 定期的にログを確認し、異常を検知

## 参考リソース

- [Go公式ドキュメント](https://go.dev/doc/)
- [Next.js Deployment](https://nextjs.org/docs/deployment)
- [nginx Documentation](https://nginx.org/en/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Let's Encrypt](https://letsencrypt.org/)
