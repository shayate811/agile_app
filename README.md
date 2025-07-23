# TODOアプリ

このアプリは、コマンドラインで動作するシンプルなTODO管理ツールです。タスクの追加・一覧表示・完了・削除が可能で、各タスクにはスプリント番号とタスクの重み（ウェイト）を設定できます。

## 機能

- タスクの追加（タイトル、スプリント番号、タスクウェイトを指定）
- タスクの一覧表示
- 割当者の追加・削除
- タスクの完了
- タスクの削除
- スプリントタイマーの設定
- スプリントタイマーの開始
- プロジェクトの進捗確認
- 貢献度の確認

## 前提条件

このアプリを使用するには、以下がインストールされている必要があります：

- Go 1.18

Go のインストール方法は [公式サイト](https://golang.org/doc/install) を参照してください。

## インストール

### 方法1: go install を使用（推奨）

```bash
go install github.com/shayate811/agile_app@latest
```

この方法でインストールすると、`agile_app` コマンドとして実行できます。

**注意**: `$GOPATH/bin` または `$HOME/go/bin` がPATHに追加されていることを確認してください。

### 方法2: ソースからビルド

```bash
# リポジトリをクローン
git clone https://github.com/shayate811/agile_app.git
cd agile_app

# 依存関係をインストール
go mod tidy

# ビルド
go build -o agile_app

# 実行
./agile_app
```

## 使い方

### 1. タスクの追加

```
agile_app add <タイトル> <スプリント番号> <タスクウェイト>
```

例:
```
agile_app add "shiryou_sakusei" 1 3
agile_app add "kaigi_junbi" 2 2
```

### 2. タスクの一覧表示

```
agile_app list
```

### 3. 割当者の追加・削除
```
# 追加
agile_app assign <タスクID> <割当者の名前>
# 削除
agile_app assign <タスクID>
```

例:
```
# 追加
agile_app assign 2 "hanako"
agile_app assign 1 "taro"
# 削除
agile_app assign 2
```

### 4. タスクの完了

```
agile_app complete <タスクID>
```

例:
```
agile_app complete 2
```

### 5. タスクの削除

```
agile_app delete <タスクID>
```

例:
```
agile_app delete 3
```

### 6. スプリントタイマーの設定

スプリントの各フェーズにかかる時間を設定できます。

```
agile_app timersetting <計画時間> <開発時間> <レビュー時間>
```

例:
```
agile_app timersetting 30 120 60
```

### 7. スプリントタイマーの開始

設定された時間でスプリントタイマーを開始します。

```
agile_app timerstart
```

### 8. プロジェクトの進捗確認

現在のプロジェクトの進捗状況を表示します。

```
agile_app progress
```

### 9. 貢献度の確認

チームメンバーの貢献度を表示します。

```
agile_app contribution
```

## データ保存

タスク情報は `todo.json` ファイルに保存されます。

## コマンド一覧

| コマンド | 説明 | 使用例 |
|---------|------|--------|
| add | タスクを追加 | `agile_app add "shiryou_sakusei" 1 3` |
| list | タスク一覧を表示 | `agile_app list` |
| assign | 割当者を設定/削除 | `agile_app assign 2 "hanako"` |
| complete | タスクを完了 | `agile_app complete 2` |
| delete | タスクを削除 | `agile_app delete 3` |
| timersetting | タイマー設定 | `agile_app timersetting 30 120 60` |
| timerstart | タイマー開始 | `agile_app timerstart` |
| progress | 進捗確認 | `agile_app progress` |
| contribution | 貢献度確認 | `agile_app contribution` |

## トラブルシューティング

### `agile_app: command not found` エラーが出る場合

1. Go がインストールされているか確認:
   ```bash
   go version
   ```

2. GOPATH/bin がPATHに追加されているか確認:
   ```bash
   echo $PATH | grep go
   ```

3. 必要に応じてPATHを追加（`.bashrc` または `.zshrc` に追加）:
   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

---

## 注意事項

- スプリント番号とタスクウェイトは数値で指定してください。
- 存在しないIDを指定した場合は「task not found」と表示されます。
- 割当者を削除する場合は、名前を指定せずにassignコマンドを実行してください。
- タイマー設定の時間は分単位で指定してください。
- タスクタイトルや割当者名は英数字（ローマ字）で指定してください。