# 订单模块 API 文档

此类接口都需要用户登录，需要在请求头中带入 token

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| go-mall-token | 是 | string | 用户 access_token |

## 订单管理

### 创建订单

会删除购物车中相应的购物项

- 请求路径：`/order/create`
- 请求方式：POST
- 请求头：
  - go-mall-token: {access_token}
- 请求参数：

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| cart_item_id_list | 是 | array | 购物项 ID 列表 |
| user_address_id | 是 | int | 用户地址 ID |

```json
{
    "cart_item_id_list": [1,2,3],
    "user_address_id": 1
}
```

- 响应数据：

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| order_no | 是 | string | 订单编号 |

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "b353d9287ab3061c",
    "data": {
        "order_no": "20250313036149887007930002"
    }
}
```

### 获取用户订单列表

- 请求路径：`/order/user-order/?page=1&page_size=3`
- 请求方式：GET
- 请求头：
  - go-mall-token: {access_token}
- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "19eed83c3063921f",
    "data": [
        {
            "order_no": "20250313036149887007930002",
            "pay_trans_id": "",
            "pay_type": 0,
            "bill_money": 199000,
            "pay_money": 198800,
            "pay_state": 1,
            "status": "待付款",
            "address": {
                "user_name": "七�****�弟",
                "user_phone": "130****1234",
                "province_name": "上海",
                "city_name": "上海",
                "region_name": "松江区",
                "detail_address": "九亭U天地"
            },
            "items": [
                {
                    "commodity_id": 70,
                    "commodity_name": "圣罗兰（YSL）纯口红13#（正橘色）3.8g",
                    "commodity_img": "https://static.toastmemo.com/img/go-mall/upload/53a4a428-8ca2-4d19-937d-15d18f324237.jpg",
                    "commodity_selling_price": 32000,
                    "commodity_num": 4
                },
                {
                    "commodity_id": 68,
                    "commodity_name": "纪梵希高定香榭天鹅绒唇膏306#(小羊皮口红 法式红 雾面哑光",
                    "commodity_img": "https://static.toastmemo.com/img/go-mall/upload/f30bd8cb-aadd-43aa-8615-2c4795ee7f5f.jpg",
                    "commodity_selling_price": 35500,
                    "commodity_num": 2
                }
            ],
            "created_at": "2025-03-13 16:25:16"
        }
    ]
}
```

### 获取订单详情

- 请求路径：`/order/:order_no/info`
- 请求方式：GET
- 请求头：
  - go-mall-token: {access_token}
- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "f810c2ac6915138b",
    "data": {
        "order_no": "20250313036149887007930002",
        "pay_trans_id": "",
        "pay_type": 0,
        "bill_money": 199000,
        "pay_money": 198800,
        "pay_state": 1,
        "status": "待付款",
        "address": {
            "user_name": "七**弟",
            "user_phone": "130****1234",
            "province_name": "上海",
            "city_name": "上海",
            "region_name": "松江区",
            "detail_address": "九亭U天地"
        },
        "items": [
            {
                "commodity_id": 70,
                "commodity_name": "圣罗兰（YSL）纯口红13#（正橘色）3.8g",
                "commodity_img": "https://static.toastmemo.com/img/go-mall/upload/53a4a428-8ca2-4d19-937d-15d18f324237.jpg",
                "commodity_selling_price": 32000,
                "commodity_num": 4
            },
            {
                "commodity_id": 68,
                "commodity_name": "纪梵希高定香榭天鹅绒唇膏306#(小羊皮口红 法式红 雾面哑光",
                "commodity_img": "https://static.toastmemo.com/img/go-mall/upload/f30bd8cb-aadd-43aa-8615-2c4795ee7f5f.jpg",
                "commodity_selling_price": 35500,
                "commodity_num": 2
            }
        ],
        "created_at": "2025-03-13 16:25:16"
    }
}
```

### 取消订单

- 请求路径：`/order/:order_no/cancel`
- 请求方式：PATCH
- 请求头：
  - go-mall-token: {access_token}
- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "fa36d2cc424be45b",
    "data": ""
}
```

### 发起订单支付

- 请求路径：`/order/create-pay`
- 请求方式：POST
- 请求头：
  - go-mall-token: {access_token}
- 请求参数：

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| order_no | 是 | string | 订单编号 |
| pay_type | 是 | int | 支付类型 1：微信；2：暂无 |

```json
{
    "order_no": "string",
    "pay_type": 1
}
```

- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "c1715f5053dc8fb0",
    "data": {
        "appId": "123456",
        "timeStamp": "1741854866",
        "nonceStr": "e61463f8efa94090b1f366cccfbbb444",
        "package": "prepay_id=wx21201855730335ac86f8c43d1889123400",
        "signType": "RSA",
        "paySign": "oR9d8PuhnIc+YZ8cBHFCwfgpaK9gd7vaRvkYD7rthRAZ/X+QBhcCYL21N7cHCTUxbQ+EAt6Uy+lwSN22f5YZvI45MLko8Pfso0jm46v5hqcVwrk6uddkGuT+Cdvu4WBqDzaDjnNa5UK3GfE1Wfl2gHxIIY5lLdUgWFts17D4WuolLLkiFZV+JSHMvH7eaLdT9N5GBovBwu5yYKUR7skR8Fu+LozcSqQixnlEZUfyE55feLOQTUYzLmR9pNtPbPsu6WVhbNHMS3Ss2+AehHvz+n64GDmXxbX++IOBvm2olHu3PsOUGRwhudhVf7UcGcunXt8cqNjKNqZLhLw4jq/xDg=="
    }
}
```
