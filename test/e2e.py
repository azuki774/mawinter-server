# pip install mysql-connector-python
import requests
import sys
import mysql.connector

cnx = None

print('# Database setup start')

try:
    cnx = mysql.connector.connect(
        user='root',  # ユーザー名
        password='password',  # パスワード
        host='localhost',  # ホスト名(IPアドレス）
        database='mawinter'  # データベース名
    )

    cursor = cnx.cursor()

    cursor.execute("DROP TABLE IF EXISTS Record_200004")
    cursor.execute("DROP TABLE IF EXISTS Record_200005")
    cursor.execute("DROP TABLE IF EXISTS Record_200006")
    cursor.execute("DROP TABLE IF EXISTS Record_200007")
    cursor.execute("DROP TABLE IF EXISTS Record_200008")
    cursor.execute("DROP TABLE IF EXISTS Record_200009")
    cursor.execute("DROP TABLE IF EXISTS Record_200010")
    cursor.execute("DROP TABLE IF EXISTS Record_200011")
    cursor.execute("DROP TABLE IF EXISTS Record_200012")
    cursor.execute("DROP TABLE IF EXISTS Record_200101")
    cursor.execute("DROP TABLE IF EXISTS Record_200102")
    cursor.execute("DROP TABLE IF EXISTS Record_200103")

    cursor.close()

except Exception as e:
    print(f"Error Occurred: {e}")

print('# Database setup complete')


print('# health check')
url = 'http://localhost:8080/'
response = requests.get(url)
if response.status_code == 200:
    print("[OK] {}".format(url))
else:
    print("[NG] {}".format(url))
    print(response.status_code)
    sys.exit(1)

print('# create table')
url = 'http://localhost:8080/v2/table/2000'
response = requests.post(url)
if response.status_code == 201:
    print("[OK] {}".format(url))
else:
    print("[NG] {}".format(url))
    print(response.status_code)
    sys.exit(1)

print('# create table already exists')
url = 'http://localhost:8080/v2/table/2000'
response = requests.post(url)
if response.status_code == 204:
    print("[OK] {}".format(url))
else:
    print("[NG] {}".format(url))
    print(response.status_code)
    sys.exit(1)

print('# create records')
url = 'http://localhost:8080/v2/record'
data1 = '{"category_id": 100, "datetime": "20000410", "from": "testfrom1", "type": "S1", "price": 1000, "memo": "memo"}'
data2 = '{"category_id": 200, "datetime": "20000415", "from": "testfrom2", "type": "", "price": 2000, "memo": ""}'
headers = { "Content-Type": "application/json" }
response = requests.post(url, data=data1, headers = headers)
if response.status_code == 201:
    print("[OK] {}".format(url))
else:
    print("[NG] {}".format(url))
    print(response.status_code)
    sys.exit(1)

response = requests.post(url, data=data2, headers = headers)
if response.status_code == 201:
    print("[OK] {}".format(url))
else:
    print("[NG] {}".format(url))
    print(response.status_code)
    sys.exit(1)

print('# get records')
url = 'http://localhost:8080/v2/record/200004'
response = requests.get(url)
if response.status_code == 200:
    json_data = response.json()
    want = [{"category_id": 100, "category_name" : "月給", "id" : 1, "datetime": "2000-04-10T00:00:00+09:00", "from": "testfrom1", "type": "S1", "price": 1000, "memo": "memo"}, 
            {"category_id": 200, "category_name" : "家賃", "id" : 2, "datetime": "2000-04-15T00:00:00+09:00", "from": "testfrom2", "type": "", "price": 2000, "memo": ""}
           ]
    if want != json_data:
        print("[NG] {}".format(url))
        print(json_data)
        print(want)
        sys.exit(1)
    print("[OK] {}".format(url))
else:
    print("[NG] {}".format(url))
    print(response.status_code)
    sys.exit(1)
