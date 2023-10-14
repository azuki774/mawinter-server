-- +migrate Up
ALTER TABLE `Record_YYYYMM` RENAME TO `Record`;

-- +migrate Down
ALTER TABLE `Record` RENAME TO `Record_YYYYMM`;
