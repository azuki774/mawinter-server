# mawinter-expense
- 家計簿 Rest API サーバ + MariaDB
- kubernetes single node 実行環境

## 使い方
```
kubectl apply -f manifest/prd/
```

## 初回セットアップ
- `docs/setup.md` を参照
    
### データについて
- DBのデータ本体は `/data/maw-pv` が使われる。
- DBのdumpが、 `/data/maw-dump` に出力される。
