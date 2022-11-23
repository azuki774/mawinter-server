## DB 保存場所について
- DBのデータ本体は `/data/maw-pv` が使われる。
- DBのdumpが、 `/data/maw-dump` に出力される

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
## Category
    +-------------+--------------+------+-----+---------+----------------+
    | Field       | Type         | Null | Key | Default | Extra          |
    +-------------+--------------+------+-----+---------+----------------+
    | id          | int(11)      | NO   | PRI | NULL    | auto_increment |
    | category_id | int(11)      | NO   | UNI | NULL    |                |
    | name        | varchar(255) | YES  |     | NULL    |                |
    +-------------+--------------+------+-----+---------+----------------+


## Record_YYYYMM
    +-------------+-------------+------+-----+---------------------+-------------------------------+
    | Field       | Type        | Null | Key | Default             | Extra                         |
    +-------------+-------------+------+-----+---------------------+-------------------------------+
    | id          | int(11)     | NO   | PRI | NULL                | auto_increment                |
    | category_id | int(11)     | NO   | MUL | NULL                |                               |
    | from        | varchar(64) | NO   |     | NULL                |                               |
    | type        | varchar(64) | NO   |     | NULL                |                               |
    | created_at  | datetime    | YES  |     | current_timestamp() |                               |
    | updated_at  | timestamp   | NO   |     | current_timestamp() | on update current_timestamp() |
    +-------------+-------------+------+-----+---------------------+-------------------------------+
