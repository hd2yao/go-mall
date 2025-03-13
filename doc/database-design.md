# 数据库设计文档

## 概述

本文档描述了 go-mall 电商系统的数据库设计。系统使用 MySQL 作为主要数据库，采用 InnoDB 引擎，字符集使用 utf8mb4。

## 表设计

### 用户相关表

#### users 用户表

```sql
CREATE TABLE `users` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` varchar(50) NOT NULL COMMENT '用户名',
    `password` varchar(100) NOT NULL COMMENT '密码（加密存储）',
    `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
    `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
    `avatar` varchar(200) DEFAULT NULL COMMENT '头像URL',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-正常',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`),
    KEY `idx_email` (`email`),
    KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

#### user_addresses 用户地址表

```sql
CREATE TABLE `user_addresses` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '地址ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `receiver` varchar(50) NOT NULL COMMENT '收货人',
    `phone` varchar(20) NOT NULL COMMENT '联系电话',
    `province` varchar(50) NOT NULL COMMENT '省份',
    `city` varchar(50) NOT NULL COMMENT '城市',
    `district` varchar(50) NOT NULL COMMENT '区县',
    `detail` varchar(200) NOT NULL COMMENT '详细地址',
    `is_default` tinyint NOT NULL DEFAULT '0' COMMENT '是否默认地址：0-否，1-是',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户地址表';
```

### 商品相关表

#### categories 商品分类表

```sql
CREATE TABLE `categories` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '分类ID',
    `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '父分类ID',
    `name` varchar(50) NOT NULL COMMENT '分类名称',
    `level` tinyint NOT NULL DEFAULT '1' COMMENT '层级：1-一级，2-二级，3-三级',
    `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-正常',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品分类表';
```

#### products 商品表

```sql
CREATE TABLE `products` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '商品ID',
    `category_id` bigint unsigned NOT NULL COMMENT '分类ID',
    `name` varchar(100) NOT NULL COMMENT '商品名称',
    `brief` varchar(200) DEFAULT NULL COMMENT '商品简介',
    `description` text COMMENT '商品描述',
    `price` decimal(10,2) NOT NULL COMMENT '销售价',
    `market_price` decimal(10,2) NOT NULL COMMENT '市场价',
    `stock` int NOT NULL DEFAULT '0' COMMENT '库存',
    `sales` int NOT NULL DEFAULT '0' COMMENT '销量',
    `images` json DEFAULT NULL COMMENT '商品图片JSON',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-下架，1-上架',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';
```

#### product_specs 商品规格表

```sql
CREATE TABLE `product_specs` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '规格ID',
    `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
    `name` varchar(50) NOT NULL COMMENT '规格名称',
    `value` varchar(50) NOT NULL COMMENT '规格值',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品规格表';
```

### 购物车相关表

#### carts 购物车表

```sql
CREATE TABLE `carts` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '购物车ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
    `quantity` int NOT NULL DEFAULT '1' COMMENT '数量',
    `selected` tinyint NOT NULL DEFAULT '1' COMMENT '是否选中：0-否，1-是',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_product` (`user_id`,`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='购物车表';
```

### 订单相关表

#### orders 订单表

```sql
CREATE TABLE `orders` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '订单ID',
    `order_no` varchar(50) NOT NULL COMMENT '订单编号',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `total_amount` decimal(10,2) NOT NULL COMMENT '订单总金额',
    `pay_amount` decimal(10,2) NOT NULL COMMENT '实付金额',
    `freight_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '运费',
    `status` tinyint NOT NULL DEFAULT '0' COMMENT '订单状态：0-待支付，1-已支付，2-已发货，3-已完成，4-已取消',
    `payment_type` tinyint DEFAULT NULL COMMENT '支付方式：1-支付宝，2-微信',
    `payment_time` timestamp NULL DEFAULT NULL COMMENT '支付时间',
    `delivery_time` timestamp NULL DEFAULT NULL COMMENT '发货时间',
    `complete_time` timestamp NULL DEFAULT NULL COMMENT '完成时间',
    `address_snapshot` json NOT NULL COMMENT '收货地址快照',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_order_no` (`order_no`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';
```

#### order_items 订单项表

```sql
CREATE TABLE `order_items` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '订单项ID',
    `order_id` bigint unsigned NOT NULL COMMENT '订单ID',
    `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
    `product_snapshot` json NOT NULL COMMENT '商品快照',
    `quantity` int NOT NULL COMMENT '数量',
    `price` decimal(10,2) NOT NULL COMMENT '单价',
    `total_amount` decimal(10,2) NOT NULL COMMENT '总金额',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单项表';
```

### 支付相关表

#### payments 支付记录表

```sql
CREATE TABLE `payments` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '支付ID',
    `order_id` bigint unsigned NOT NULL COMMENT '订单ID',
    `payment_no` varchar(50) NOT NULL COMMENT '支付流水号',
    `transaction_id` varchar(100) DEFAULT NULL COMMENT '第三方交易号',
    `payment_type` tinyint NOT NULL COMMENT '支付方式：1-支付宝，2-微信',
    `amount` decimal(10,2) NOT NULL COMMENT '支付金额',
    `status` tinyint NOT NULL DEFAULT '0' COMMENT '支付状态：0-待支付，1-支付成功，2-支付失败',
    `callback_content` text COMMENT '回调内容',
    `callback_time` timestamp NULL DEFAULT NULL COMMENT '回调时间',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_payment_no` (`payment_no`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_transaction_id` (`transaction_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付记录表';
```

## 索引设计说明

1. 主键索引
   - 所有表都使用自增 ID 作为主键
   - 使用 unsigned bigint 类型，避免负数

2. 唯一索引
   - users 表的 username 字段
   - orders 表的 order_no 字段
   - payments 表的 payment_no 字段

3. 普通索引
   - 外键关联字段（如 user_id, product_id 等）
   - 常用查询条件字段（如 status, category_id 等）
   - 排序字段（如 sort）

## 字段设计规范

1. 时间字段
   - created_at: 创建时间
   - updated_at: 更新时间
   - deleted_at: 删除时间（软删除）

2. 状态字段
   - 使用 tinyint 类型
   - 默认值为正常状态
   - 注释中说明状态含义

3. 金额字段
   - 使用 decimal(10,2) 类型
   - 精确到分

4. JSON 类型
   - 用于存储结构化数据
   - 如商品图片、地址快照等

## 数据库优化建议

1. 分表策略
   - orders 表按用户 ID 范围分表
   - order_items 表按订单 ID 范围分表

2. 缓存策略
   - 商品信息缓存
   - 用户信息缓存
   - 购物车数据缓存

3. 性能优化
   - 合理使用索引
   - 避免大事务
   - 定期维护索引和统计信息
