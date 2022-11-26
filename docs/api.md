## GET /
- healthCheck 用

以下、すべてBasic認証が必要

## GET /summary/year/{year}
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

## GET /record/recent/
- 最新20件のレコードを表示する。
### response:
    [
        {
            "id" : 123, 
            "category_id" : 400,
            "category_name" : "cat1", 
            "date" : "2021-01-01T00:00:00Z",
            "from" : "discord",
            "type" : "",
            "price" : 1234,
            "memo": ""
        },
        {
            "id" : 124, 
            "category_id" : 500,
            "category_name" : "cat2", 
            "date" : "2021-01-01T00:00:00Z",
            "from" : "discord",
            "type" : "",
            "price" : 1234,
            "memo": ""
        }
    ]

## DELETE /record/{id}
- id の record を削除する
### response:
- 成功 .. 204 No Contents
- 失敗
    - データがない場合 .. 404
    - 何らかの場合で失敗 .. 500
