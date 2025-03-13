# 购物车模块 API 文档

此类接口都需要用户登录，需要在请求头中带入 token

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| go-mall-token | 是 | string | 用户 access_token |

## 购物车管理

### 获取购物车列表

- 请求路径：`/cart/item/`
- 请求方式：GET
- 请求头：
  - go-mall-token: {access_token}
- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "4897136c2e8ada8c",
    "data": [
        {
            "cart_item_id": 1,
            "user_id": 1,
            "commodity_id": 1,
            "commodity_num": 6,
            "commodity_name": "Apple iPhone 11 (A2223)",
            "commodity_img": "https://static.toastmemo.com/img/go-mall/upload/4755f3e5-257c-424c-a5f4-63908061d6d9.jpg",
            "commodity_selling_price": 549900,
            "add_cart_at": "2025-01-21 16:46:23"
        },
        {
            "cart_item_id": 2,
            "user_id": 1,
            "commodity_id": 21,
            "commodity_num": 1,
            "commodity_name": "Apple iPhone XS Max",
            "commodity_img": "https://static.toastmemo.com/img/go-mall/upload/ec4af4a5-0a53-4246-bd88-919b0541a55c.jpg",
            "commodity_selling_price": 899900,
            "add_cart_at": "2025-01-21 17:20:33"
        },
        {
            "cart_item_id": 3,
            "user_id": 1,
            "commodity_id": 52,
            "commodity_num": 4,
            "commodity_name": "Apple 苹果 iPhone xr",
            "commodity_img": "https://static.toastmemo.com/img/go-mall/upload/41b10e86-857c-435c-b86d-d822e35450ab.jpg",
            "commodity_selling_price": 507900,
            "add_cart_at": "2025-01-22 11:26:37"
        }
    ]
}
```

### 添加商品到购物车

- 请求路径：`/cart/add-item`
- 请求方式：POST
- 请求头：
  - go-mall-token: {access_token}
- 请求参数：

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| commodity_id | 是 | int | 商品 ID |
| commodity_num | 是 | int | 商品数量，一个商品往购物车里一次性最多放 5 个 |

```json
{
    "commodity_id": 70,
    "commodity_num": 3
}
```

- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "6a539e804a369573",
    "data": ""
}
```

### 更新购物车商品数量

一次请求只允许修改一个购物项目，修改只限于商品数量

- 请求路径：`/cart/update-item`
- 请求方式：PATCH
- 请求头：
  - go-mall-token: {access_token}
- 请求参数：

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| item_id | 是 | int | 购物车 ID |
| commodity_num | 是 | int | 商品数量，一个商品往购物车里一次性最多放 5 个 |

```json
{
    "item_id": 4,
    "commodity_num": 4
}
```

- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "ce2a8a68456aa916",
    "data": ""
}
```

### 删除购物车商品

一次请求只允许删除一个购物项目

- 请求路径：`/cart/item/:item_id`
- 请求方式：DELETE
- 请求头：
  - go-mall-token: {access_token}
- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "cfc1c0784c4981fd",
    "data": ""
}
```

### 查看购物项账单

- 请求路径：`/cart/item/check-bill?item_id=1&item_id=2&...`
- 请求方式：GET
- 请求头：
  - go-mall-token: {access_token}
- 请求参数：

- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "ee22280868ddb1b5",
    "data": {
        "items": [
            {
                "cart_item_id": 1,
                "user_id": 1,
                "commodity_id": 1,
                "commodity_num": 6,
                "commodity_name": "Apple iPhone 11 (A2223)",
                "commodity_img": "https://static.toastmemo.com/img/go-mall/upload/4755f3e5-257c-424c-a5f4-63908061d6d9.jpg",
                "commodity_selling_price": 549900,
                "add_cart_at": "2025-01-21 16:46:23"
            },
            {
                "cart_item_id": 2,
                "user_id": 1,
                "commodity_id": 21,
                "commodity_num": 1,
                "commodity_name": "Apple iPhone XS Max",
                "commodity_img": "https://static.toastmemo.com/img/go-mall/upload/ec4af4a5-0a53-4246-bd88-919b0541a55c.jpg",
                "commodity_selling_price": 899900,
                "add_cart_at": "2025-01-21 17:20:33"
            }
        ],
        "bill_detail": {
            "coupon": {
                "coupon_id": 1,
                "coupon_name": "",
                "discount_money": 100
            },
            "discount": {
                "discount_id": 1,
                "discount_name": "",
                "discount_money": 100
            },
            "vip_discount_money": 0,
            "original_total_price": 4199300,
            "total_price": 4199100
        }
    }
}
```
