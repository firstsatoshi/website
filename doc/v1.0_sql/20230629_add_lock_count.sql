alter table website.tb_blindbox_event ADD COLUMN `lock_count` int not NULL COMMENT '锁定数量';

UPDATE website.tb_blindbox_event SET lock_count = 200 WHERE id = 1;