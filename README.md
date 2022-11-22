# cmd\_supervisor

ブラウザからコマンドの実行を管理するWebサービス。

# 使い方

## 0. 準備

```bash
git clone https://github.com/wakewakame/cmd_supervisor.git
cd cmd_supervisor
cp .env_template .env
```

上記コマンド実行後、 `.env` を適宜書き換える。

## 1. 起動

```
make up
make bash
go run main.go
```

## 2. コマンド実行

`http://localhost` にアクセスし、入力欄からコマンドを実行する。

## 3. 終了

```
make down
```
