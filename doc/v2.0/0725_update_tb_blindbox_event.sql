ALTER table tb_blindbox_event ADD COLUMN `custome_mint` after `only_whitelist` tinyint(1) DEFAULT '0' COMMENT '是否是自定义mint的项目,类似bitfish可以自定义合成';


