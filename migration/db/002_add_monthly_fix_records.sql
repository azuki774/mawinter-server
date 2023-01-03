-- +migrate Up
CREATE TABLE `Monthly_Fix_Billing` (
  `id` int NOT NULL AUTO_INCREMENT,
  `category_id` int NOT NULL,
  `day` int NOT NULL,
  `price` int NOT NULL,
  `type` varchar(64) NOT NULL,
  `memo` varchar(255) NOT NULL,
  `created_at` datetime default current_timestamp,
  `updated_at` timestamp default current_timestamp on update current_timestamp,
  PRIMARY KEY (`id`)
);

CREATE TABLE `Monthly_Fix_Done` (
  `yyyymm` varchar(6) NOT NULL,
  `done` tinyint(1) NOT NULL,
  `created_at` datetime default current_timestamp,
  `updated_at` timestamp default current_timestamp on update current_timestamp,
  PRIMARY KEY (`yyyymm`)
);

-- +migrate Down
DROP TABLE `Monthly_Fix_Billing`;
DROP TABLE `Monthly_Fix_Done`;
