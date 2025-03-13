# 商品模块 API 文档

## 商品分类

### 按层级划分的查询所有商品

- 请求路径：`/commodity/category-hierarchy/`
- 请求方式：GET
- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "2d3fb3e4586d4816",
    "data": [
        {
            "id": 1,
            "level": 1,
            "parent_id": 0,
            "name": "家电 数码 手机",
            "icon_img": "",
            "rank": 100,
            "sub_categories": [
                {
                    "id": 3,
                    "level": 2,
                    "parent_id": 1,
                    "name": "家电",
                    "icon_img": "",
                    "rank": 10,
                    "sub_categories": [
                        {
                            "id": 6,
                            "level": 3,
                            "parent_id": 3,
                            "name": "生活电器",
                            "icon_img": "",
                            "rank": 0,
                            "sub_categories": null
                        },
                        // ...
                        {
                            "id": 17,
                            "level": 3,
                            "parent_id": 3,
                            "name": "空气净化器",
                            "icon_img": "",
                            "rank": 0,
                            "sub_categories": null
                        }
                    ]
                }
                // ...
            ]
        },
        // ...
        {
            "id": 94,
            "level": 1,
            "parent_id": 0,
            "name": "测试分类2",
            "icon_img": "",
            "rank": 0,
            "sub_categories": [
                {
                    "id": 95,
                    "level": 2,
                    "parent_id": 94,
                    "name": "测试分类2-1",
                    "icon_img": "",
                    "rank": 0,
                    "sub_categories": null
                }
            ]
        }
    ]
}
```

### 获取商品子类列表

- 请求路径：`/commodity/category/?parent_id=1`
- 请求方式：GET
- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "a52962c66dda6d2",
    "data": [
        {
            "id": 3,
            "level": 2,
            "parent_id": 1,
            "name": "家电",
            "icon_img": "",
            "rank": 10
        },
        {
            "id": 4,
            "level": 2,
            "parent_id": 1,
            "name": "数码",
            "icon_img": "",
            "rank": 9
        },
        {
            "id": 5,
            "level": 2,
            "parent_id": 1,
            "name": "手机",
            "icon_img": "",
            "rank": 8
        }
    ]
}
```

## 商品管理

### 按分类查询商品列表

- 请求路径：`/commodity/commodity-in-cate/?category_id=48&page=1&page_size=20`
- 请求方式：GET
- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "6c20e24026d8f1ce",
    "data": [
        {
            "id": 64,
            "name": "迪奥（Dior）烈艳蓝金唇膏-哑光999# 3.5g 传奇红（口红",
            "intro": "雾面质地 显色持久 显白 正红色 李佳琦推荐）",
            "category_id": 72,
            "cover_img": "https://static.toastmemo.com/img/go-mall/upload/d8d4ac7e-7189-459a-aef2-7116f723cb0b.jpg",
            "original_price": 40000,
            "selling_price": 31500,
            "tag": "",
            "sell_status": 1,
            "created_at": "2025-01-21 12:30:37"
        },
        // ...
        {
            "id": 83,
            "name": "海囤全球 魅可（MAC)经典唇膏 子弹头口红3g",
            "intro": "Chili 秀智色/小辣椒色",
            "category_id": 72,
            "cover_img": "https://static.toastmemo.com/img/go-mall/upload/2b678c5d-820c-4174-bc0c-5a65ff9501b6.jpg",
            "original_price": 17000,
            "selling_price": 15500,
            "tag": "",
            "sell_status": 1,
            "created_at": "2025-01-21 12:30:37"
        }
    ],
    "Pagination": {
        "page": 1,
        "page_size": 20,
        "total_rows": 38
    }
}
```

### 商品搜索

- 请求路径：`/commodity/search?keyword=iPhone&page=1&page_size=3`
- 请求方式：GET
- 请求参数：

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| keyword | 是 | string | 关键词 |
| page | 是 | int | 页码，最小为1 |
| page_size | 否 | int | 大小 |

- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "ef9f014c29547033",
    "data": [
        {
            "id": 1,
            "name": "Apple iPhone 11 (A2223)",
            "intro": "64GB 黑色 移动联通电信4G手机 双卡双待",
            "category_id": 33,
            "cover_img": "https://static.toastmemo.com/img/go-mall/upload/4755f3e5-257c-424c-a5f4-63908061d6d9.jpg",
            "original_price": 549900,
            "selling_price": 549900,
            "tag": "2019 新品",
            "sell_status": 1,
            "created_at": "2025-01-21 12:30:37"
        },
        {
            "id": 2,
            "name": "Apple iPhone 11 (A2223)",
            "intro": "128GB 白色 移动联通电信4G手机 双卡双待",
            "category_id": 33,
            "cover_img": "https://static.toastmemo.com/img/go-mall/upload/a0d09f94-9c46-4ee1-aaef-dfd132e7543e.jpg",
            "original_price": 599900,
            "selling_price": 599900,
            "tag": "2019 新品",
            "sell_status": 1,
            "created_at": "2025-01-21 12:30:37"
        },
        {
            "id": 3,
            "name": "Apple iPhone 11 (A2223)",
            "intro": "128GB 紫色 移动联通电信4G手机 双卡双待",
            "category_id": 33,
            "cover_img": "https://static.toastmemo.com/img/go-mall/upload/8dfe8ea9-2279-4132-a72b-4f8a52d002a4.jpg",
            "original_price": 599900,
            "selling_price": 599900,
            "tag": "",
            "sell_status": 1,
            "created_at": "2025-01-21 12:30:37"
        }
    ],
    "Pagination": {
        "page": 1,
        "page_size": 3,
        "total_rows": 55
    }
}
```

### 商品详情

- 请求路径：`/commodity/:commodity_id/info`
- 请求方式：GET
- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "a12418300b2b94c2",
    "data": {
        "id": 1,
        "name": "Apple iPhone 11 (A2223)",
        "intro": "64GB 黑色 移动联通电信4G手机 双卡双待",
        "category_id": 33,
        "cover_img": "https://static.toastmemo.com/img/go-mall/upload/4755f3e5-257c-424c-a5f4-63908061d6d9.jpg",
        "images": "",
        "detail_content": "",
        "original_price": 549900,
        "selling_price": 549900,
        "stock_num": 994,
        "tag": "2019 新品",
        "sell_status": 1,
        "created_at": "2025-01-21T12:30:37+08:00"
    }
}
```
