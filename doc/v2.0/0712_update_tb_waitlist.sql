ALTER table tb_waitlist ADD COLUMN `event_id` bigint NOT NULL COMMENT '活动id';
UPDATE tb_waitlist SET event_id = 1;