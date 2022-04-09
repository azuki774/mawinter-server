## DB 保存場所について
- DBのデータ本体は `/data/maw-pv` が使われる。
- DBのdumpが、 `/data/maw-dump` に出力される

## DB 手動操作
    kubectl exec --stdin --it <db-pod-name>  -- bash
    mysql -u root -ppassword
    use mawinter

## rootの権限
    mysql> grant all privileges on *.* to root@"%";
    MariaDB [mawinter]> show grants for root;
    +------------------------------------------------------------------------
    --------------------------------------------------------+
    | Grants for root@%                                                      
                                                            |
    +------------------------------------------------------------------------
    --------------------------------------------------------+
    | GRANT ALL PRIVILEGES ON *.* TO `root`@`%` IDENTIFIED BY PASSWORD '*2470
    C0C06DEE42FD1618BB99005ADCA2EC9D1E19' WITH GRANT OPTION |
    +------------------------------------------------------------------------
    --------------------------------------------------------+

でなければならない。

---
# DB Schema
## types

    MariaDB [mawinter]> SHOW COLUMNS FROM types;
    +-------+--------------+------+-----+---------+----------------+
    | Field | Type         | Null | Key | Default | Extra          |
    +-------+--------------+------+-----+---------+----------------+
    | id    | int(11)      | NO   | PRI | NULL    | auto_increment |
    | type  | varchar(255) | YES  | UNI | NULL    |                |
    +-------+--------------+------+-----+---------+----------------+

- 値は
1. income
2. outgoing
3. saving
4. investment

- `categories` に紐づく

## categories

    MariaDB [mawinter]> SHOW COLUMNS FROM categories;
    +-------------+--------------+------+-----+---------+----------------+
    | Field       | Type         | Null | Key | Default | Extra          |
    +-------------+--------------+------+-----+---------+----------------+
    | id          | int(11)      | NO   | PRI | NULL    | auto_increment |
    | category_id | int(11)      | NO   | UNI | NULL    |                |
    | name        | varchar(255) | YES  |     | NULL    |                |
    | type        | varchar(255) | YES  | MUL | NULL    |                |
    +-------------+--------------+------+-----+---------+----------------+

## records

    MariaDB [mawinter]> SHOW COLUMNS FROM records;
    +-------------+--------------+------+-----+---------------------+----------------+
    | Field       | Type         | Null | Key | Default             | Extra          |
    +-------------+--------------+------+-----+---------------------+----------------+
    | id          | int(11)      | NO   | PRI | NULL                | auto_increment |
    | category_id | int(11)      | NO   | MUL | NULL                |                |
    | date        | datetime     | YES  |     | current_timestamp() |                |
    | price       | int(11)      | NO   |     | NULL                |                |
    | memo        | varchar(255) | YES  |     | NULL                |                |
    +-------------+--------------+------+-----+---------------------+----------------+



