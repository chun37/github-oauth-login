# プロジェクト要件

## アーキテクチャ
- モノレポ構成
- DDD（ドメイン駆動設計）を採用
- SOLIDの原則に従ってコードを書くこと

## バックエンド (Golang + Echo)
- REST APIを実装
- GitHub OAuth認証を実装
- セッションはCookieに保存
- データベースはPostgreSQL
- podmanを使用

## フロントエンド (TypeScript + Next.js)
- SPAを実装
- GitHub OAuthログインフロー

## GitHub OAuth
- アプリケーションは既に作成済
- Client IDとClient Secretは取得済
- 認証情報はプロジェクトルートの.envに記述
- .env.exampleを作成すること

## 機能要件
- GitHubでログイン
- ログイン後に取得したアクセストークンを使用してGitHub Profile APIを呼び出す
- 取得したプロフィール情報を表示
- ログアウト機能は不要

## セッション管理
- Cookie の有効期限: 1年 (365日)
- セッション情報の保存先: PostgreSQL
- Cookie設定:
  - Secure属性: true (HTTPでも動作可能)
  - HttpOnly: true
  - SameSite: Lax

## プロフィール情報
- GitHubから取得したプロフィール情報はデータベースに保存しない
- 都度APIを呼び出して取得する

## 技術要件
- ポート番号:
  - バックエンド: 8080
  - フロントエンド: 3000
  - PostgreSQL: 5432
- 開発環境用とプロダクション用の両方を作成（プロダクションで内容が変わる場合）
- リバースプロキシは不要
- ライブラリは最新バージョンを調査して使用
- ベストプラクティスは都度調査すること

## 技術選定 (2025年1月調査結果)

### バックエンド
- Echo: v4 (最新版)
- PostgreSQL ドライバ: `github.com/jackc/pgx/v5`
- セッション管理: `github.com/alexedwards/scs/v2` + PostgreSQLストア
- OAuth2: `golang.org/x/oauth2`

### フロントエンド
- Next.js: 15 (App Router を使用)
- TypeScript: 最新版

### インフラ
- Podman Compose
- PostgreSQL: 最新版

## DDD実装ガイドライン
- Entities: 可変で識別可能な構造体
- Value Objects: 不変で識別不可能な構造体
- Aggregates: EntitiesとValue Objectsの組み合わせ
- Repositories: Aggregatesの永続化を担当
- Services: リポジトリとサブサービスを組み合わせてビジネスフローを構築
- レイヤー構造:
  - domain: ビジネスロジック
  - application: ユースケース固有の操作
  - infrastructure: データベースアクセスなどの技術的機能

## 実装前の確認
- ユーザに指示された内容を必ずCLAUDE.mdに保存
- ユーザの指示内容とCLAUDE.mdの内容に相違がないことを確認してから実装開始
