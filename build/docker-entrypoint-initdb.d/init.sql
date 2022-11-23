USE mawinter;
grant all privileges on *.* to root@"%";

CREATE TABLE `Category` (
  `id` int NOT NULL AUTO_INCREMENT,
  `category_id` int NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime  default current_timestamp,
  `updated_at` timestamp default current_timestamp on update current_timestamp,
  PRIMARY KEY (`id`),
  UNIQUE KEY `category_id` (`category_id`)
);

create table `Record_YYYYMM` (
  `id` int NOT NULL AUTO_INCREMENT,
  `category_id` int NOT NULL,
  `from` varchar(64) NOT NULL,
  `type` varchar(64) NOT NULL,
  `created_at` datetime  default current_timestamp,
  `updated_at` timestamp default current_timestamp on update current_timestamp,
  PRIMARY KEY (`id`),
  index `idx_cat` (`category_id`)
);
