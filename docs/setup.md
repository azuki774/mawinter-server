### DBセットアップ
```
build/docker-entrypoint-initdb.d の `init.sql` にある、`Record_YYYYMM` テーブルをテンプレートに、
使う年月(YYYYMM)分、`Record_YYYYMM` テーブルを作っておく。
```

### Basic認証
```
apiVersion: v1
kind: Secret
metadata:
  name: maw-basic-auth
type: kubernetes.io/basic-auth
stringData:
  username: ******
  password: ******
```
