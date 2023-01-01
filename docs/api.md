## GET /
- healthCheck 用

以下、すべてBasic認証が必要

## GET /record/year/{year}
- 年間サマリを取得する。
### response:

    [
        {
            "category_id": 210,
            "category_name": "種類",
            "price": [4,5,6,7,8,9,10,11,12,1,2,3], // 4月から3月までの数値が配列で返る
            "total": 213912
        },
        {
            "category_id": 211,
            "category_name": "種類",
            "price": [4,5,6,7,8,9,10,11,12,1,2,3],
            "total": 4210
        },
    ]
- price のところにはその年の 4,5,6,...,3月 の合計が入る。

## POST /record/
- データを追加する。
### request:

    {
        "category_id" : 120,
        "datetime" : "20210101",
        "from" : "discord",
        "type" : "",
        "price" : 210,
        "memo" : ""
    }

### response:
    {
        "id" : 123, 
        "category_id" : 400,
        "category_name" : "cat1", 
        "date" : "2021-01-01T00:00:00Z",
        "from" : "discord",
        "type" : "",
        "price" : 1234,
        "memo": ""
    }
- `date` フィールドが空のときは現在時刻が入る。

## POST /record/fixmonth/
- テーブル `Fix_Monthly_Billing` のデータを `Record_YYYYMM` データに追加する。

### request:
None

### response
- 201 Created 成功した場合
    ```
    [
        {
            "category_id" : 100, 
            "day" : 2,
            "type" : "type1",
            "price" : 1234,
            "memo" : "memo1",
        },
        {
            "category_id" : 101, 
            "day" : 4,
            "type" : "type2",
            "price" : 10000,
            "memo" : "memo2",
        },
    ]

    ```
- 400 Bad Request その月がすでに追加済の場合

## POST /table/{year}
- FY{year} 用のテーブルを生成する。
- すでに生成済の場合は何もしない
### response:
- 201 Created
