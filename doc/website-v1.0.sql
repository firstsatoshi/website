SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `tb_waitlist`;
CREATE TABLE `tb_waitlist` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '邮箱',
  `btc_address` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'BTC的P2TR格式地址',
  `referee_id` bigint DEFAULT '0' COMMENT '邀请人id',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`) USING BTREE,
	UNIQUE KEY `email` (`email`) USING BTREE,
	UNIQUE KEY `btc_address` (`btc_address`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='waitlist数据表';

DROP TABLE IF EXISTS `tb_blindbox`;
CREATE TABLE `tb_blindbox` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `name` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '名称',
  `description` varchar(200) COLLATE utf8mb4_bin DEFAULT '' COMMENT '描述',
  `category` varchar(20) COLLATE utf8mb4_bin NOT NULL COMMENT '分类: bald,punk,rich,elite',
  `img_url` varchar(500) COLLATE utf8mb4_bin DEFAULT "" COMMENT '图片url',
  `is_active` tinyint(1) DEFAULT '1' COMMENT '是否激活',
  `is_locked` tinyint(1) DEFAULT '0' COMMENT '是否锁定',
  `status` varchar(20) COLLATE utf8mb4_bin DEFAULT 'NOTMINT'  COMMENT '状态,NOTMINT,MINTING,MINT',
  `commit_txid` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'commit_txid',
  `reveal_txid` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '铭文交易id',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_category` (`category`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='盲盒表';

DROP TABLE IF EXISTS `tb_blindbox_event`;
CREATE TABLE `tb_blindbox_event` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `event_name` varchar(100) COLLATE utf8mb4_bin DEFAULT "" COMMENT '名称',
  `event_description` varchar(200) COLLATE utf8mb4_bin DEFAULT "" COMMENT '描述,富文本',
  `background_img_url` varchar(500) COLLATE utf8mb4_bin DEFAULT "" COMMENT '背景图片url',
  `roadmap_description` varchar(1000) COLLATE utf8mb4_bin DEFAULT "" COMMENT '路线图描述,富文本',
  `roadmap_list` varchar(1000) COLLATE utf8mb4_bin DEFAULT "" COMMENT '路线图;按照 title1;title2;title3 的格式',
  `website_url` varchar(200) COLLATE utf8mb4_bin DEFAULT "" COMMENT '官网url',
  `whitepaper_url` varchar(200) COLLATE utf8mb4_bin DEFAULT "" COMMENT '白皮书url',
  `twitter_url` varchar(200) COLLATE utf8mb4_bin DEFAULT "" COMMENT 'twitter url',
  `discord_url` varchar(200) COLLATE utf8mb4_bin DEFAULT "" COMMENT 'discord url',
  `price_sats` int not NULL COMMENT '价格',
  `is_active` tinyint(1) DEFAULT '1' COMMENT '是否激活',
  `payment_token` varchar(20) NOT NULL COMMENT '支付币种',
  `img_url_list` varchar(500) COLLATE utf8mb4_bin DEFAULT "" COMMENT '图片url列表,按照url1;url2;url3格式',
  `average_image_bytes` int not NULL COMMENT '平均图片大小(字节数)',
  `supply` int not NULL COMMENT '供应量',
  `avail` int not NULL COMMENT '当前可用',
  `lock_count` int DEFAULT 0 COMMENT '锁定数量',
  `mint_limit` int DEFAULT 2 COMMENT '单个地址限购数量',
  `only_whitelist` tinyint(1) DEFAULT '0' COMMENT '是否只有白名单',
  `start_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
  `end_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '结束时间',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='盲盒活动表';




DROP TABLE IF EXISTS `tb_order`;
CREATE TABLE `tb_order` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `order_id` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '订单id',
  `event_id` bigint NOT NULL COMMENT '活动id',
  `count` int not NULL COMMENT '数量',
  `deposit_address` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '充值地址',
  `inscription_data` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '铭刻内容',
  `fee_rate` int NOT NULL COMMENT '费率 n/sat',
  `txfee_amount_sat` int NOT NULL COMMENT '矿工费',
  `service_fee_sat` int NOT NULL COMMENT '服务费',
  `price_sat` int NOT NULL COMMENT '价格',
  `total_amount_sat` int NOT NULL COMMENT '总费用sat',
  `receive_address` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '铭刻内容接收地址',
  `order_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '订单状态: NOTPAID未支付;PAYPENDING支付确认中;PAYSUCCESS支付成功;PAYTIMEOUT超时未支付;MINTING铭刻交易等待确认中;ALLSUCCESS订单成功',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间(进入内存池的时间)',
  `pay_txid` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '付款交易id(支持批量支付,即一笔交易多个输出到我们平台的收款地址,所以不必设置为唯一索引)',
  `pay_confirmed_time` datetime DEFAULT NULL COMMENT '付款交易确认时间',
  `pay_from_address` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '付款地址',
  `real_fee_sat` int DEFAULT '0' COMMENT '实际矿工费',
  `real_change_sat` int DEFAULT '0' COMMENT '实际找零(收入)',
  `version` bigint NOT NULL DEFAULT '0' COMMENT '版本号',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `order_id` (`order_id`) USING BTREE,
  UNIQUE KEY `tb_order_deposit_address_UIDX` (`deposit_address`) USING BTREE,
  KEY `tb_order_receive_address_IDX` (`receive_address`) USING BTREE,
  KEY `tb_order_order_status_IDX` (`order_status`) USING BTREE,
  KEY `tb_order_pay_from_address_IDX` (`pay_from_address`) USING BTREE,
  KEY `tb_order_pay_txid_IDX` (`pay_txid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='盲盒订单表';


DROP TABLE IF EXISTS `tb_inscribe_order`;
CREATE TABLE `tb_inscribe_order` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `order_id` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '订单id',
  `count` int not NULL COMMENT '数量(批量)',
  `deposit_address` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '充值地址',
  `fee_rate` int NOT NULL COMMENT '费率 n/sat',
  `data_bytes` int NOT NULL COMMENT '数据大小(字节数)',
  `txfee_amount_sat` int NOT NULL COMMENT '矿工费',
  `service_fee_sat` int NOT NULL COMMENT '服务费',
  `total_amount_sat` int NOT NULL COMMENT '总费用sat',
  `receive_address` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '铭刻内容接收地址',
  `order_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '订单状态: NOTPAID未支付;PAYPENDING支付确认中;PAYSUCCESS支付成功;PAYTIMEOUT超时未支付;MINTING铭刻交易等待确认中;ALLSUCCESS订单成功',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间(进入内存池的时间)',
  `pay_txid` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '付款交易id(支持批量支付,即一笔交易多个输出到我们平台的收款地址,所以不必设置为唯一索引)',
  `pay_confirmed_time` datetime DEFAULT NULL COMMENT '付款交易确认时间',
  `pay_from_address` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '付款地址',
  `real_fee_sat` int DEFAULT '0' COMMENT '实际矿工费',
  `real_change_sat` int DEFAULT '0' COMMENT '实际找零(收入)',
  `version` bigint NOT NULL DEFAULT '0' COMMENT '版本号',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `order_id` (`order_id`) USING BTREE,
  UNIQUE KEY `tb_inscribe_order_deposit_address_UIDX` (`deposit_address`) USING BTREE,
  KEY `tb_inscribe_order_receive_address_IDX` (`receive_address`) USING BTREE,
  KEY `tb_inscribe_order_order_status_IDX` (`order_status`) USING BTREE,
  KEY `tb_inscribe_order_pay_from_address_IDX` (`pay_from_address`) USING BTREE,
  KEY `tb_inscribe_order_pay_txid_IDX` (`pay_txid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='铭文订单表';


DROP TABLE IF EXISTS `tb_inscribe_data`;
CREATE TABLE `tb_inscribe_data` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `order_id` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '订单id',
  `content_type` varchar(20) COLLATE utf8mb4_bin NOT NULL COMMENT '类型: 如 image/img',
  `file_name` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '文件名',
  `data` mediumblob NOT NULL COMMENT '铭文数据',
  `status` varchar(20) COLLATE utf8mb4_bin DEFAULT 'NOTMINT'  COMMENT '状态,NOTMINT,MINTING,MINT',
  `commit_txid` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'commit_txid',
  `reveal_txid` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '铭文交易id',
  `deleted` tinyint(1) DEFAULT '0' COMMENT '逻辑删除',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='铭文数据表';



DROP TABLE IF EXISTS `tb_address`;
CREATE TABLE `tb_address` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `address` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '地址',
  `coin_type`varchar(10) COLLATE utf8mb4_bin NOT NULL COMMENT '地址类型,BTC,ETH,USDT',
  `account_index` bigint NOT NULL COMMENT 'account_index',
  `address_index` bigint NOT NULL COMMENT 'address_index',
  `bussines_type` varchar(20) COLLATE utf8mb4_bin DEFAULT 'BLINDBOX' COMMENT '业务类型: BLINDBOX, INSCRIBE',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `tb_address_uindex` (`address`) USING BTREE,
  UNIQUE KEY `tb_coin_type_account_index_address_index_IDX` ( `coin_type`, `account_index`,`address_index`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='收款地址表';



DROP TABLE IF EXISTS `tb_deposit`;
CREATE TABLE `tb_deposit` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `coin_type`varchar(20) COLLATE utf8mb4_bin NOT NULL COMMENT '地址类型,BTC,ETH,USDT',
  `from_address` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT 'from地址,如果是btc归集充值,显示输入的第一个地址',
  `to_address` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT 'to地址',
  `txid` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT 'txid',
  `amount` int NOT NULL COMMENT '金额(最小单位)',
  `decimals` int NOT NULL COMMENT '精度(BTC: 8, ETH: 18, USDT: 6)',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `tb_deposit_to_address_txid` (`to_address`, `txid`) USING BTREE,
  KEY `tb_deposit_txid` (`txid` ) USING BTREE,
  KEY `tb_deposit_from_address` (`from_address`) USING BTREE,
  KEY `tb_deposit_type_from_address` (`coin_type`, `from_address`) USING BTREE,
  KEY `tb_deposit_to_address` (`to_address`) USING BTREE,
  KEY `tb_deposit_type_to_address` (`coin_type`, `to_address`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='充值记录表';


DROP TABLE IF EXISTS `tb_blockscan`;
CREATE TABLE `tb_blockscan` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `coin_type`varchar(20) COLLATE utf8mb4_bin NOT NULL COMMENT '地址类型,BTC,ETH,USDT',
  `block_number` bigint  NOT NULL COMMENT '区块高度',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uindex_cointype` (`coin_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='区块扫描记录表';

DROP TABLE IF EXISTS `tb_lock_order_blindbox`;
CREATE TABLE `tb_lock_order_blindbox` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `event_id` bigint NOT NULL COMMENT '活动id',
  `order_id` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '订单号',
  `blindbox_id` bigint NOT NULL COMMENT '盲盒id',
  `version` bigint NOT NULL DEFAULT '0' COMMENT '版本号',
  `deleted` tinyint(1) DEFAULT '0' COMMENT '逻辑删除',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `tb_blindbox_id_IDX` (`blindbox_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='锁库存表';


SET FOREIGN_KEY_CHECKS = 1;

