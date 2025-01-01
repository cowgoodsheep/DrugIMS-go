CREATE TABLE `user_info`(  
    `user_id` int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '用户id,自增',
    `user_name` varchar(64) NOT NULL COMMENT '用户昵称',
    `password` varchar(64) NOT NULL COMMENT '密码',
    `telephone` varchar(64) NOT NULL COMMENT '用户电话',
    `description` varchar(512) COMMENT '用户描述',
    `avatar` varchar(512) COMMENT '用户头像',
    `address` varchar(256) null comment '地址',
    `role` char(1)  null comment '用户角色,1:管理员;2:客户;3:供应商',
    `status` tinyint NOT NULL DEFAULT '1',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间'
) COMMENT '用户信息表';