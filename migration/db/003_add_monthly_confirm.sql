-- +migrate Up
CREATE TABLE `Monthly_Confirm` (
  `yyyymm` varchar(6) NOT NULL,
  `confirm` tinyint(1) NOT NULL,
  `confirm_datetime` datetime,
  `created_at` datetime default current_timestamp,
  `updated_at` timestamp default current_timestamp on update current_timestamp,
  PRIMARY KEY (`yyyymm`)
);

-- +migrate Down
DROP TABLE `Monthly_Confirm`;
