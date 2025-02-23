CREATE TABLE `user_info` (
    `user_id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '用户id,自增',
    `user_name` VARCHAR(64) NOT NULL COMMENT '用户昵称',
    `password` VARCHAR(64) NOT NULL COMMENT '密码',
    `telephone` VARCHAR(64) NOT NULL COMMENT '用户电话',
    `description` VARCHAR(512) COMMENT '用户描述',
    `avatar` VARCHAR(512) COMMENT '用户头像',
    `address` VARCHAR(256) NULL COMMENT '地址',
    `role` VARCHAR(10) NULL COMMENT '用户角色,管理员;客户;供应商',
    `balance` DECIMAL(10, 2) NULL COMMENT '余额',
    `block_balance` DECIMAL(10, 2) NULL COMMENT '冻结余额',
    `status` TINYINT NOT NULL DEFAULT '1',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间'
) COMMENT '用户信息表';

CREATE TABLE `drug_info` (
    `drug_id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '药品ID, 主键',
    `drug_name` VARCHAR(100) NOT NULL COMMENT '药品名称',
    `manufacturer` VARCHAR(100) NULL COMMENT '生产厂家',
    `unit` VARCHAR(50) NULL COMMENT '单位',
    `specification` VARCHAR(50) NULL COMMENT '规格',
    `stock_lower_limit` INT NOT NULL COMMENT '库存下限',
    `stock_upper_limit` INT NOT NULL COMMENT '库存上限',
    `sale_price` DECIMAL(10, 2) unsigned NOT NULL COMMENT '售价',
    `drug_description` VARCHAR(500) NULL COMMENT '药品说明',
    `img` VARCHAR(1000) NULL COMMENT '药品图片',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间'
) COMMENT '药品信息表';

CREATE TABLE `stock_info` (
    `stock_id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '库存ID, 主键',
    `drug_id` INT NOT NULL COMMENT '药品ID',
    `batch_number` VARCHAR(50) NULL COMMENT '批号',
    `production_date` DATE NULL COMMENT '生产日期',
    `supply_price` DECIMAL(10, 2) NULL COMMENT '进货单价',
    `remaining_quantity` INT NULL COMMENT '剩余数量',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间'
) COMMENT '库存信息表';

CREATE TABLE `order_info` (
    `order_id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '订单ID, 主键',
    `drug_id` INT NOT NULL COMMENT '药品ID',
    `user_id` INT NOT NULL COMMENT '客户ID',
    `sale_quantity` INT NOT NULL COMMENT '销售数量',
    `sale_amount` DECIMAL(10, 2) NOT NULL COMMENT '销售金额',
    `supply_amount` DECIMAL(10, 2) NOT NULL COMMENT '进货金额',
    `stock_info` VARCHAR(1024) COMMENT '扣除库存',
    `order_status` INT NOT NULL COMMENT '订单状态,0处理中,1已完成,2已撤销,3待确认,4退款审核中,5已退款',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间'
) COMMENT '销售信息表';

CREATE TABLE `supply_order` (
    `supply_id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '进货单ID, 主键',
    `drug_id` INT NOT NULL COMMENT '药品ID',
    `user_id` INT NOT NULL COMMENT '供应商ID',
    `batch_number` VARCHAR(64) NULL COMMENT '批号',
    `production_date` DATE NULL COMMENT '生产日期',
    `supply_quantity` INT NOT NULL COMMENT '进货数量',
    `supply_price` DECIMAL(10, 2) NULL COMMENT '进货单价',
    `note` VARCHAR(1024) NULL COMMENT '备注',
    `supply_status` INT NOT NULL COMMENT '进货单状态,0审批中,1已完成,2已拒绝',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间'
) COMMENT '进货单表';

CREATE TABLE `approval_info` (
    `approval_id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '审批单ID, 主键',
    `user_id` INT NOT NULL COMMENT '用户ID',
    `approval_type` INT NOT NULL COMMENT '审批类型,0退货审批,1进货审批',
    `approval_content` VARCHAR(10240) NOT NULL COMMENT '审批内容',
    `reason` VARCHAR(1024) NULL COMMENT '申请理由',
    `approval_user_id` INT NOT NULL COMMENT '审批人ID',
    `approval_opinion` VARCHAR(1024) NULL COMMENT '申请意见',
    `approval_status` INT NOT NULL COMMENT '审批状态,0审核中,1已通过,2已拒绝',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间'
) COMMENT '审批单表';