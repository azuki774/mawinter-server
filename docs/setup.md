## 必要なkubernetesリソース
### 初回DBセットアップ
```
kubectl exec --stdin --it <db-pod-name>  -- bash
mysql -u root -ppassword
use mawinter

build/docker-entrypoint-initdb.d/init.sql の中身をコピペして入れる
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
