# 用户模块 API 文档

## 认证相关

### 用户注册

- 请求路径：`/user/register`
- 请求方式：POST
- 请求参数：

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| login_name | 是 | string | 登录名，必须是手机号或邮箱 |
| password | 是 | string | 密码，最少是八位 |
| password_confirm | 是 | string | 确认密码 |
| nickname | 否 | string | 昵称 |
| slogan | 否 | string | 个性签名 |
| avatar | 否 | string | 头像 |

```json
{
    "login_name": "hd2yao@gmail.com",
    "password": "1231@QQ.com",
    "password_confirm": "1231@QQ.com",
    "nickname": "Dysania",
    "slogan": "你在看孤独的风景",
    "avatar": ""
}
```

- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "73f6c1ac0b892039",
    "data": ""
}
```

### 用户登录

- 请求路径：`/user/login`
- 请求方式：POST
- 请求头：

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| platform | 是 | string | 登录平台，必须是 APP 或 H5 |

- 请求参数：

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| login_name | 是 | string | 登录名，必须是手机号或邮箱 |
| password | 是 | string | 密码，最少是八位 |

```json
{
    "login_name": "string",
    "password": "string"
}
```

- 响应数据：

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| access_token | 是 | string | 用于用户认证的 token |
| refresh_token | 是 | string | 更新 access_token 的 token |

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "2789e0447327b302",
    "data": {
        "access_token": "2501056762ec33862f428a45cea765d9edf62d3f",
        "refresh_token": "2501056762ec33862f428a45cea765d9edf62d3f",
        "duration": 7200,
        "srv_create_time": "2025-03-13 13:28:30"
    }
}
```

### 刷新 Token

- 请求路径：`/user/token/refresh?refresh_token={refresh_token}`
- 请求方式：GET
- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "29c5edd878757a28",
    "data": {
        "access_token": "beb685b08f1c2d2e5ea6367aba45de38bcd73ada",
        "refresh_token": "beb685b08f1c2d2e5ea6367aba45de38bcd73ada",
        "duration": 7200,
        "srv_create_time": "2025-03-13 13:31:30"
    }
}
```

### 用户登出

- 请求路径：`/user/logout`
- 请求方式：DELETE
- 请求头：
  - go-mall-token: {access_token}

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| go-mall-token | 是 | string | 用户 access_token |

- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "d64b5e5c54103de6",
    "data": ""
}
```

## 密码管理

### 申请重置密码

- 请求路径：`/user/password/apply-reset`
- 请求方式：POST
- 请求参数：

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| login_name | 是 | string | 登录名，必须是手机号或邮箱 |

```json
{
    "login_name": "hd2yao@gmail.com"
}
```

- 响应数据：

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| password_reset_token | 是 | string | 用于密码重置的 token |

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "1c4088f440b68443",
    "data": {
        "password_reset_token": "18fb7bf4d6eba03a7dde938a779eb2fe82c0eb99"
    }
}
```

### 重置密码

- 请求路径：`/user/password/reset`
- 请求方式：POST
- 请求参数：

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| password | 是 | string | 新密码，最少是八位 |
| password_confirm | 是 | string | 新密码确认 |
| password_reset_token | 是 | string | 为 `/user/password/apply-reset` 接口返回的 token 值 |
| password_reset_code | 是 | string | 需要从 redis 中获取 |

```json
{
    "password": "12312@QQ.com",
    "password_confirm": "12312@QQ.com",
    "password_reset_token": "18fb7bf4d6eba03a7dde938a779eb2fe82c0eb99",
    "password_reset_code": "403886"
}
```

- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "79b5dd8803669623",
    "data": ""
}
```

## 用户信息管理（需要登录）

此类接口都需要在请求头中带入 token

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| go-mall-token | 是 | string | 用户 access_token |

### 获取用户信息

- 请求路径：`/user/info`
- 请求方式：GET
- 请求头：
  - go-mall-token: {access_token}
- 响应数据：

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| verified | 是 | int | 用户验证状态 0：未验证；1：已验证 |
| is_blocked | 是 | int | 用户禁用状态 0：正常；1：已禁用 |

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "7f1ecf840df875b3",
    "data": {
        "id": 2,
        "nickname": "Dysania",
        "login_name": "hd**ao@gmail.com",
        "verified": 0,
        "avatar": "",
        "slogan": "你在看孤独的风景",
        "is_blocked": 0,
        "created_at": "2025-03-13 13:23:17"
    }
}
```

### 更新用户信息

只能修改用户 昵称、个性签名、头像

- 请求路径：`/user/info`
- 请求方式：PATCH
- 请求头：
  - go-mall-token: {access_token}
- 请求参数：

```json
{
    "nickname": "Dysania-2",
    "slogan": "你在看孤独的风景！",
    "avatar": ""
}
```

- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "a44ec1d4583a2962",
    "data": ""
}
```

## 收货地址管理（需要登录）

此类接口都需要在请求头中带入 token

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| go-mall-token | 是 | string | 用户 access_token |

### 新增收货地址

- 请求路径：`/user/address`
- 请求方式：POST
- 请求头：
  - go-mall-token: {access_token}
- 请求参数：

| 参数名 | 必选 | 类型 | 描述 |
|-------|------|------|-----|
| user_name | 是 | string | 收件人名称 |
| user_phone | 是 | string | 收件人联系方式 |
| default | 是 | string | 是否默认地址 0：否；1：是 |
| province_name | 是 | string | 省 |
| city_name | 是 | string | 市 |
| region_name | 是 | string | 区 |
| detail_address | 是 | string | 详细地址 |

```json
{
    "user_name": "七七",
    "user_phone": "13012341234",
    "default": 1,
    "province_name": "上海",
    "city_name": "上海",
    "region_name": "松江区",
    "detail_address": "金地广场"
}
```

- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "d276643c26aaba04",
    "data": ""
}
```

### 获取所有收货地址

- 请求路径：`/user/address/`
- 请求方式：GET
- 请求头：
  - go-mall-token: {access_token}
- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "917707807bf37e6a",
    "data": [
        {
            "id": 3,
            "user_name": "七七",
            "user_phone": "13012341234",
            "masked_user_name": "七*",
            "masked_user_phone": "130****1234",
            "default": 1,
            "province_name": "上海",
            "city_name": "上海",
            "region_name": "松江区",
            "detail_address": "金地广场",
            "created_at": "2025-03-13 13:55:23"
        },
        {
            "id": 4,
            "user_name": "七七",
            "user_phone": "13012341234",
            "masked_user_name": "七*",
            "masked_user_phone": "130****1234",
            "default": 0,
            "province_name": "上海",
            "city_name": "上海",
            "region_name": "松江区",
            "detail_address": "九亭U天地",
            "created_at": "2025-03-13 13:56:15"
        }
    ]
}
```

### 获取单个收货地址

- 请求路径：`/user/address/:address_id`
- 请求方式：GET
- 请求头：
  - go-mall-token: {access_token}
- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "2b5e344c07fb5e1b",
    "data": {
        "id": 3,
        "user_name": "七七",
        "user_phone": "13012341234",
        "masked_user_name": "七*",
        "masked_user_phone": "130****1234",
        "default": 1,
        "province_name": "上海",
        "city_name": "上海",
        "region_name": "松江区",
        "detail_address": "金地广场",
        "created_at": "2025-03-13 13:55:23"
    }
}
```

### 更新收货地址

- 请求路径：`/user/address/:address_id`
- 请求方式：PATCH
- 请求头：
  - go-mall-token: {access_token}
- 请求参数：

```json
{
    "user_name": "七七弟弟",
    "user_phone": "13012341234",
    "default": 1,
    "province_name": "上海",
    "city_name": "上海",
    "region_name": "松江区",
    "detail_address": "金地广场"
}
```

- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "f33f95dc6ec6edbf",
    "data": ""
}
```

### 删除收货地址

- 请求路径：`/user/address/:address_id`
- 请求方式：DELETE
- 请求头：
  - go-mall-token: {access_token}
- 响应数据：

```json
{
    "code": 0,
    "msg": "success",
    "request_id": "581c856053a46a91",
    "data": ""
}
```
