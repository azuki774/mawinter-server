# Monthly_Fix_Billing

## Description

<details>
<summary><strong>Table Definition</strong></summary>

```sql
CREATE TABLE `Monthly_Fix_Billing` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `category_id` int(11) NOT NULL,
  `day` int(11) NOT NULL,
  `price` int(11) NOT NULL,
  `type` varchar(64) NOT NULL,
  `memo` varchar(255) NOT NULL,
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=[Redacted by tbls] DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
```

</details>

## Columns

| Name | Type | Default | Nullable | Extra Definition | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | ---------------- | -------- | ------- | ------- |
| id | int(11) |  | false | auto_increment |  |  |  |
| category_id | int(11) |  | false |  |  |  |  |
| day | int(11) |  | false |  |  |  |  |
| price | int(11) |  | false |  |  |  |  |
| type | varchar(64) |  | false |  |  |  |  |
| memo | varchar(255) |  | false |  |  |  |  |
| created_at | datetime | current_timestamp() | true |  |  |  |  |
| updated_at | timestamp | current_timestamp() | false | on update current_timestamp() |  |  |  |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| PRIMARY | PRIMARY KEY | PRIMARY KEY (id) |

## Indexes

| Name | Definition |
| ---- | ---------- |
| PRIMARY | PRIMARY KEY (id) USING BTREE |

## Relations

![er](Monthly_Fix_Billing.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
