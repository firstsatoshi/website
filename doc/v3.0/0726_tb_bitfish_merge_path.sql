DROP TABLE IF EXISTS `tb_bitfish_merge_path`;
CREATE TABLE `tb_bitfish_merge_path` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `merge_path`varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '合成路径',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uindex_merge_path` (`merge_path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='bitfish合成路径';

