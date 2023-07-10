ALTER table tb_blindbox_event ADD COLUMN `event_endpoint` varchar(100) COLLATE utf8mb4_bin DEFAULT "" COMMENT 'url endpoint,例如: fsat.io/collection/biteagle' AFTER `event_name`;
UPDATE tb_blindbox_event SET event_endpoint = 'biteagle' WHERE id = 1;
create unique index uidx_event_endpoit on tb_blindbox_event(event_endpoint) USING BTREE;