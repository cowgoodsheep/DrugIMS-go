CREATE TABLE `user_info`(  
    `user_id` int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '用户id,自增',
    `user_name` varchar(64) NOT NULL COMMENT '用户昵称',
    `password` varchar(64) NOT NULL COMMENT '密码',
    `telephone` varchar(64) NOT NULL COMMENT '用户电话',
    `description` varchar(512) COMMENT '用户描述',
    `avatar` varchar(512) COMMENT '用户头像',
    `address` varchar(256) NULL comment '地址',
    `role` varchar(10)  NULL comment '用户角色,管理员;客户;供应商',
    `status` tinyint NOT NULL DEFAULT '1',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间'
) COMMENT '用户信息表';

create TABLE `drug_info`(
    `drug_id`           int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '药品ID, 主键',
    `drug_name`         varchar(100)  NOT NULL COMMENT '药品名称',
    `manufacturer`      varchar(100)  NULL COMMENT '生产厂家',
    `unit`              varchar(50)   NULL COMMENT '单位',
    `specification`     varchar(50)   NULL COMMENT '规格',
    `stock_lower_limit` int           NOT NULL COMMENT '库存下限',
    `stock_upper_limit` int           NOT NULL COMMENT '库存上限',
    `price`             int unsigned  NOT NULL COMMENT '售价',
    `drug_description`  varchar(500)  NULL COMMENT '药品说明',
    `img`               varchar(1000) NULL COMMENT '药品图片',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间'
) COMMENT '药品信息表';

create TABLE `stock_info`(
    `stock_id`            int NOT NULL PRIMARY KEY AUTO_INCREMENT comment '库存ID, 主键' ,
    `drug_id`             int           NOT NULL comment '药品ID',
    `batch_number`        varchar(50)   NULL comment '批号',
    `production_date`     date          NULL comment '生产日期',
    `purchase_date`       date          NULL comment '进货日期',
    `purchase_unit_price` decimal(10, 2) NULL comment '进货单价',
    `remaining_quantity`  int            NULL comment '剩余数量',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间'
) comment '库存信息表';