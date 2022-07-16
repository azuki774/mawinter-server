import sys
import base64
import json
import requests
from requests.auth import HTTPBasicAuth

auth = HTTPBasicAuth("test", "test")


def main():
    print("test start")

    # 1-1 health check
    url = "http://localhost:8080/"
    res = requests.get(url)
    if res.status_code != 200:
        print("1-1. health check failed")
        sys.exit(1)
    print("1-1. health check ok")

    # 1-2 get summary
    url = "http://localhost:8080/summary/year/2021"
    res = requests.get(url, auth=auth)
    if res.status_code != 200:
        print("1-2. get summary failed")
        print(res.status_code)
        sys.exit(1)
    print("1-2. get summary ok")

    # 1-3 get recent
    url = "http://localhost:8080/record/recent/"
    res = requests.get(url, auth=auth)
    if res.status_code != 200:
        print("1-3. get recent failed")
        print(res.status_code)
        sys.exit(1)
    print("1-3. get recent ok")

    # 2-1 add record
    url = "http://localhost:8080/record/"
    json_data = {"category_id": 210, "price": 1}
    res = requests.post(url, json.dumps(json_data), auth=auth)
    if res.status_code != 201:
        print("2-1. add record")
        print(res.status_code)
        sys.exit(1)

    res_data = res.json()

    if res_data["id"] != 1:
        print("2-1. add record failed (id)")
        print(res_data["id"])
        sys.exit(1)
    if res_data["category_id"] != 210:
        print("2-1. add record failed (category_id)")
        print(res_data["category_id"])
        sys.exit(1)
    if res_data["category_name"] != "食費":
        print("2-1. add record failed (category_name)")
        print(res_data["category_name"])
        sys.exit(1)
    if res_data["price"] != 1:
        print("2-1. add record failed (price)")
        print(res_data["price"])
        sys.exit(1)
    if res_data["memo"] != "":
        print("2-1. add record failed (memo)")
        print(res_data["memo"])
        sys.exit(1)

    print("2-1. add record ok")

    # 2-2 add old record
    url = "http://localhost:8080/record/"
    json_data = {"category_id": 220, "price": 1234, "memo": "test", "date": "19990101"}
    res = requests.post(url, json.dumps(json_data), auth=auth)
    if res.status_code != 201:
        print("2-2. add record")
        print(res.status_code)
        sys.exit(1)

    res_data = res.json()

    if res_data["id"] != 2:
        print("2-2. add record failed (id)")
        print(res_data["id"])
        sys.exit(1)
    if res_data["category_id"] != 220:
        print("2-2. add record failed (category_id)")
        print(res_data["category_id"])
        sys.exit(1)
    if res_data["category_name"] != "電気代":
        print("2-2. add record failed (category_name)")
        print(res_data["category_name"])
        sys.exit(1)
    if res_data["date"] != "1999-01-01T00:00:00+09:00":
        print("2-2. add record failed (date)")
        print(res_data["date"])
        sys.exit(1)
    if res_data["price"] != 1234:
        print("2-2. add record failed (price)")
        print(res_data["price"])
        sys.exit(1)
    if res_data["memo"] != "test":
        print("2-2. add record failed (memo)")
        print(res_data["memo"])
        sys.exit(1)

    print("2-2. add old record ok")

    # 3-1 delete record
    url = "http://localhost:8080/record/1"
    res = requests.delete(url, auth=auth)
    if res.status_code != 204:
        print("3-1. delete record")
        print(res.status_code)
        sys.exit(1)

    sys.exit(0)


if __name__ == "__main__":
    main()
