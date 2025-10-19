# GitHub OAuth Login Application

GitHub OAuth を使用したログインシステムのサンプルアプリケーションです。

## 技術スタック

### バックエンド
- **言語**: Go 1.23
- **フレームワーク**: Echo v4
- **アーキテクチャ**: DDD (Domain-Driven Design)
- **データベース**: PostgreSQL 16
- **セッション管理**: scs v2 (PostgreSQL store)
- **OAuth**: golang.org/x/oauth2

### フロントエンド
- **フレームワーク**: Next.js 15 (App Router)
- **言語**: TypeScript
- **パッケージマネージャー**: pnpm

### インフラ
- **コンテナ**: Podman / Docker
- **リバースプロキシ**: nginx (本番環境のみ)

## 機能

- GitHub OAuth 認証
- セッション管理 (Cookie + PostgreSQL)
- GitHub プロフィール情報の取得と表示

## セットアップ

このプロジェクトは2つのデプロイ方法をサポートしています：

1. **Docker/Podmanを使用したデプロイ**（推奨、簡単）- 以下の手順を参照
2. **Docker不使用のデプロイ**（本番環境向け）- `DEPLOYMENT.md` を参照

### 前提条件（Docker/Podman使用の場合）

- Podman または Docker
- Podman Compose または Docker Compose
- GitHub OAuth App の作成

### GitHub OAuth App の作成

1. GitHubにログインし、Settings > Developer settings > OAuth Apps に移動
2. "New OAuth App" をクリック
3. 以下の情報を入力:
   - **Application name**: 任意の名前
   - **Homepage URL**: `http://127.0.0.1:8000`
   - **Authorization callback URL**: `http://127.0.0.1:8000/api/auth/callback`
4. Client ID と Client Secret を取得

### 環境変数の設定

1. プロジェクトルートに `.env` ファイルを作成:

```bash
cp .env.example .env
```

2. `.env` ファイルを編集し、GitHub OAuth の認証情報を設定:

```env
GITHUB_CLIENT_ID=your_github_client_id_here
GITHUB_CLIENT_SECRET=your_github_client_secret_here
```

### データベースマイグレーション

初回起動時は、PostgreSQL のマイグレーションを実行する必要があります:

```bash
# コンテナ起動後、バックエンドコンテナに入る
podman exec -it github-oauth-backend-dev sh

# マイグレーションを実行
psql postgresql://postgres:postgres@postgres:5432/github_oauth_app < migrations/001_create_sessions_table.up.sql
```

## 実行方法

```bash
# すべてのサービスをビルドして起動
podman-compose -f compose.yaml up --build

# または Docker Compose を使用する場合
docker-compose -f compose.yaml up --build
```

アクセス:
- アプリケーション: http://127.0.0.1:8000 (nginx経由)

**注意**: nginxをリバースプロキシとして使用します。
- フロントエンド: `http://127.0.0.1:8000/`
- バックエンドAPI: `http://127.0.0.1:8000/api/`
  - nginxが `/api/` へのリクエストを `backend:8080/api/` に転送（パスを保持）
- すべてのアクセスはポート8000経由で行います

## API エンドポイント

### 認証

- `GET /api/auth/login` - GitHub OAuth ログインページへリダイレクト
- `GET /api/auth/callback` - GitHub OAuth コールバック
- `GET /api/auth/check` - 認証状態の確認

### ユーザー

- `GET /api/user/profile` - GitHub プロフィール情報の取得（要認証）

## プロジェクト構造

```
.
├── backend/                 # バックエンド (Go + Echo)
│   ├── cmd/
│   │   └── api/
│   │       └── main.go     # エントリーポイント
│   ├── internal/
│   │   ├── domain/         # ドメイン層
│   │   ├── application/    # アプリケーション層
│   │   ├── infrastructure/ # インフラ層
│   │   └── interfaces/     # インターフェース層
│   ├── migrations/         # データベースマイグレーション
│   ├── Dockerfile
│   └── Dockerfile.dev
├── frontend/               # フロントエンド (Next.js)
│   ├── src/
│   │   ├── app/           # App Router
│   │   ├── components/
│   │   ├── lib/
│   │   └── types/
│   ├── Dockerfile
│   └── Dockerfile.dev
├── nginx/                  # nginx設定 (本番環境)
│   └── nginx.conf         # リバースプロキシ設定
├── compose.yaml            # プロダクション用
├── compose.dev.yaml        # 開発用
└── .env.example           # 環境変数のサンプル
```

## アーキテクチャ

### バックエンド

DDD (Domain-Driven Design) + Clean Architecture を採用:

- **Domain層**: ビジネスロジックとドメインモデル
- **Application層**: ユースケースとDTO
- **Infrastructure層**: データベース、外部API、セッション管理
- **Interfaces層**: HTTPハンドラー、ミドルウェア

SOLID 原則に従った設計:
- **S**ingle Responsibility Principle
- **O**pen/Closed Principle
- **L**iskov Substitution Principle
- **I**nterface Segregation Principle
- **D**ependency Inversion Principle

### フロントエンド

Next.js 15 App Router を使用:
- Server Components
- Client Components (use client)
- TypeScript による型安全性

## セッション管理

- **保存先**: PostgreSQL
- **有効期限**: 1年 (365日)
- **Cookie設定**:
  - `HttpOnly`: true
  - `Secure`: true
  - `SameSite`: Lax
  - `Persist`: true

## デプロイメント

### Docker/Podmanを使用したデプロイ

上記の「実行方法」セクションを参照してください。

### Docker不使用のデプロイ

本番環境へのデプロイ方法の詳細は `DEPLOYMENT.md` を参照してください。

以下のコンポーネントを個別にセットアップする必要があります：

- Go 1.25以上（バックエンド）
- Node.js 20以上（フロントエンド）
- PostgreSQL 16以上
- nginx（リバースプロキシ）

`DEPLOYMENT.md` には以下の詳細な手順が記載されています：

- 各コンポーネントのインストール方法
- データベースのセットアップとマイグレーション
- バックエンド・フロントエンドのビルドとデプロイ
- nginxの設定
- systemdサービスの設定
- SSL/TLS証明書の取得
- トラブルシューティング

## ライセンス

MIT
