# pip install mysql-connector-python
import requests
import sys
import time
import json
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
    cursor.execute("TRUNCATE TABLE Monthly_Fix_Billing")
    cursor.execute("TRUNCATE TABLE Monthly_Fix_Done")
    cursor.execute("INSERT INTO Monthly_Fix_Billing VALUES (1, 100, 10, 1000, '', 'memo1', '2000/01/23', '2000/01/23')")
    cursor.execute("INSERT INTO Monthly_Fix_Billing VALUES (2, 200, 20, 2000, '', 'memo2', '2000/01/23', '2000/01/23')")

    cursor.close()
    cnx.commit()
    cnx.close()

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
data3 = '{"category_id": 300, "datetime": "20000515", "from": "", "type": "", "price": 3000, "memo": ""}'
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

response = requests.post(url, data=data3, headers = headers)
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

print('# get year summary')
url = 'http://localhost:8080/v2/record/summary/2000'
response = requests.get(url)
if response.status_code == 200:
    json_data = response.json()
    want = [{'category_id': 100, 'category_name': '月給', 'count': 1, 'price': [1000, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 1000}, {'category_id': 101, 'category_name': 'ボーナス', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 110, 'category_name': '雑所得', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 200, 'category_name': '家賃', 'count': 1, 'price': [2000, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 2000}, {'category_id': 210, 'category_name': '食費', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 220, 'category_name': '電気代', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 221, 'category_name': 'ガス代', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 222, 'category_name': '水道費', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 230, 'category_name': 'コンピュータリソース', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 231, 'category_name': '通信費', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 240, 'category_name': '生活用品', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 250, 'category_name': '娯楽費', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 251, 'category_name': '交遊費', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 260, 'category_name': '書籍・勉強', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 270, 'category_name': '交通費', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 280, 'category_name': '衣服等費', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 300, 'category_name': '保険・税金', 'count': 1, 'price': [0, 3000, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 3000}, {'category_id': 400, 'category_name': '医療・衛生', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 500, 'category_name': '雑費', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 600, 'category_name': '家賃用貯金', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 601, 'category_name': 'PC用貯金', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 700, 'category_name': 'NISA入出金', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}, {'category_id': 701, 'category_name': 'NISA変動', 'count': 0, 'price': [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], 'total': 0}]
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

print('# insert fixmonth records')
url = 'http://localhost:8080/v2/record/fixmonth?yyyymm=200006'
response = requests.post(url)
if response.status_code == 201:
    json_data = response.json()
    want = [{"category_id": 100, "category_name" : "月給", "id" : 1, "datetime": "2000-06-10T00:00:00+09:00", "from": "fixmonth", "type": "", "price": 1000, "memo": "memo1"}, 
            {"category_id": 200, "category_name" : "家賃", "id" : 2, "datetime": "2000-06-20T00:00:00+09:00", "from": "fixmonth", "type": "", "price": 2000, "memo": "memo2"}
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

print('# insert fixmonth records already done')
url = 'http://localhost:8080/v2/record/fixmonth?yyyymm=200006'
response = requests.post(url)
if response.status_code == 204:
    print("[OK] {}".format(url))
else:
    print("[NG] {}".format(url))
    print(response.status_code)
    sys.exit(1)