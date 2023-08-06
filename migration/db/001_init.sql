-- +migrate Up

CREATE TABLE `Category` (
  `id` int NOT NULL AUTO_INCREMENT,
  `category_id` int NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime  default current_timestamp,
  `updated_at` timestamp default current_timestamp on update current_timestamp,
  PRIMARY KEY (`id`),
  UNIQUE KEY `category_id` (`category_id`)
);


INSERT IGNORE INTO `Category` (`id`, `category_id`, `name`) VALUES
(1,100,'月給'),
(2,101,'ボーナス'),
(3,110,'雑所得'),
(4,200,'家賃'),
(5,210,'食費'),
(6,220,'電気代'),
(7,221,'ガス代'),
(8,222,'水道費'),
(9,230,'コンピュータリソース'),
(10,231,'通信費'),
(11,240,'生活用品'),
(12,250,'娯楽費'),
(13,251,'交遊費'),
(14,260,'書籍・勉強'),
(15,270,'交通費'),
(16,280,'衣服等費'),
(17,300,'保険・税金'),
(18,400,'医療・衛生'),
(19,500,'雑費'),
(20,600,'家賃用貯金'),
(21,601,'PC用貯金'),
(22,700,'NISA入出金'),
(23,701,'NISA変動');

create table `Record_YYYYMM` (
  `id` int NOT NULL AUTO_INCREMENT,
  `category_id` int NOT NULL,
  `datetime` datetime NOT NULL default current_timestamp,
  `from` varchar(64) NOT NULL,
  `type` varchar(64) NOT NULL, -- Ignore the check if 'D' is included
  `price` int NOT NULL,
  `memo` varchar(255) NOT NULL,
  `created_at` datetime default current_timestamp,
  `updated_at` timestamp default current_timestamp on update current_timestamp,
  PRIMARY KEY (`id`),
  index `idx_cat` (`category_id`),
  index `idx_date` (`datetime`)
);


-- +migrate Down
DROP TABLE `Category`;
DROP TABLE `Record_YYYYMM`;
