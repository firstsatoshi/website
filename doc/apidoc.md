# 接口文档


- 服务地址： http://8.219.200.107:8888/api/v1/{API}
- 请求方法： 所有接口都使用 `POST`
- 示例：
    ```json
    curl -s --location 'http://8.219.200.107:8888/api/v1/joinwaitlist' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "email":"youngqqcn@163.com",
        "btcAddress":"bc1pv5d2mmx2v9cx9menxl5zlhacljqu9zqhltl4d303n3rjjcxfrgwqa80as0"
    }'
    ```



## `joinwaitlist`加入预约名单

- 请求方式: POST

- 请求参数:

    | 字段 | 说明| 类型 | 可选? | 示例 |
    |-----|------|------|----|----|
    | `email` | 邮箱 | string | 必填 | helloworld@163.com |
    | `btcAddress` | btc地址，主网地址以`bc1p`开头，测试网地址以`tb1p`开头，长度为`62`字符, | string | 必填 | bc1phjsyw73de6ap8nfjzg4erxmdw7lzlfgvm447v82fytn78nm0mwnsq654e7 |

- 请求示例：

    ```json
    {
        "email":"youngqqcn@163.com",
        "btcAddress":"bc1pv5d2mmx2v9cx9menxl5zlhacljqu9zqhltl4d303n3rjjcxfrgwqa80as0"
    }
    ```


- 响应示例
  ```json
    {
        "code": 0,
        "msg": "ok",
        "data": {
            "no": 1, // 预约序号
            "duplicated": true   // 是否重复预约，
        }
    }
  ```

每个邮箱只能预约一次


## `queryblindboxevent` 查询盲盒活动详情

- 请求方式: POST

- 请求参数: 无参数

- 请求示例：

- 响应示例:

    ```json
    {
        "code": 0,
        "msg": "ok",
        "data": [
            {
            "eventId": 1, // 活动id
            "name": "Bitcoin Eagle",  // 盲盒活动名
            "description": "This is Bitcoin Eagle NFT mint", // 盲盒活动描述
            "priceBtcSats": 123456, // 盲盒价格， 单位是聪(satoshi),如果要换算成BTC要除以10^8, 例如：123456 satoshi = 0.00123456BTC
            "priceUsd": 0, // 盲盒价格（美元），
            "paymentCoin": "BTC", // 收款币种，用户必须使用此币种进行支付
            "supply": 1000, // 总供应量（本次活动供应总量）
            "avail": 1000,  // 当前可用量(背刺活动当前可用库存)
            "enable": true, // 活动是否开启
            "onlyWhitelist": false, // 是否仅对白名单用户开放
            "startTime": "2023-05-27 16:28:39 +0800 CST", // 活动开始时间
            "endTime": "2024-05-27 16:28:39 +0800 CST" // 活动结束时间
            }
        ]
    }
    ```



## `createorder`创建订单

- 请求方式: POST

- 请求参数:

    | 字段 | 说明| 类型 | 可选? | 示例 |
    |-----|------|------|----|----|
    | `evntId` | 活动id | integer | 必填 | 1 |
    | `count` | 数量（批量） , 限制`0 < count <= 10` | integer | 必填 | 2 |
    | `receiveAddress` |btc NFT 接收地址，主网地址以`bc1p`开头，测试网地址以`tb1p`开头，长度为`62`字符,  | string | 必填 | bc1phjsyw73de6ap8nfjzg4erxmdw7lzlfgvm447v82fytn78nm0mwnsq654e7 |
    | `feeRate` | 费率 | integer | 必填 | 25 |

- 请求示例：

    ```json
    {
        "eventId": 1,
        "count": 2,
        "receiveAddress":"bc1phjsyw73de6ap8nfjzg4erxmdw7lzlfgvm447v82fytn78nm0mwnsq654e7",
        "feeRate":25
    }
    ```

- 响应示例

    ```json
    {
        "code": 0,
        "msg": "ok",
        "data": {
            "orderId": "BX2023052718471354726281", // 订单id
            "eventId": 1, // 活动id
            "count": 2, // 数量
            "depositAddress": "bc1p2yzcv24v9tpw6ffhkqcq994y8p4ps2xfv65wx7nsmg4meuvzd0fqyesxg7", // 充值地址，用户需要支付BTC到这个地址
            "receiveAddress": "bc1phjsyw73de6ap8nfjzg4erxmdw7lzlfgvm447v82fytn78nm0mwnsq654e7", //  用户提供的 BTC NFT 接收地址
            "feeRate": 25, // 费率,例如 25  表示每个字节需要25sat(聪)
            "bytes": 12345, // 盲盒字节数
            "inscribeFee": 123456, // 铭刻费用（单位是sat聪）
            "serviceFee": 123456, // 服务费 （单位是sat聪）
            "price": 123456, // 盲盒价格  （单位是sat聪）
            "total": 1123456, // 总价格 （ 单位是sat聪） , 总价格 = price + inscribeFee + serviceFee
            "createTime": "2023-05-27 19:01:15 +0800 CST" // 订单生成时间
        }
    }
    ```



## `queryorder`查询订单

- 请求方式: POST

- 请求参数:

    3个参数，**至少填1个**。 按照`orderId`，`receiveAddress`,`depositAddress`优先级查找

    | 字段 | 说明| 类型 | 可选? | 示例 |
    |-----|------|------|----|----|
    | `orderId` | 订单id | string | 可选 | `BXHJSY54E7P836W4KU01252023052917044049093727` |
    | `receiveAddress` |btc NFT 接收地址，主网地址以`bc1p`开头，测试网地址以`tb1p`开头，长度为`62`字符,  | string | 可选 | bc1phjsyw73de6ap8nfjzg4erxmdw7lzlfgvm447v82fytn78nm0mwnsq654e7 |
    | `depositAddress` | 充值地址 | string | 可选 | bc1phjsyw73de6ap8nfjzg4erxmdw7lzlfgvm447v82fytn78nm0mwnsq654e7 |

- 请求示例：

    ```json
    {
        "receiveAddress":"bc1phjsyw73de6ap8nfjzg4erxmdw7lzlfgvm447v82fytn78nm0mwnsq654e8"
    }
    ```

- 响应示例

    暂时不考虑分页

    ```json
    {
        "code": 0,
        "msg": "ok",
        "data": [
            {
                "orderId": "BXHJSY54E7P836W4KU01252023052917044049093727", // 订单id
                "eventId": 1,
                "depositAddress": "xxxxxxxxxxxxxxxxx", // 充值地址
                "total": 1123456, // 总金额（单位，聪），
                "receiveAddress": "bc1pp836v6am6mf2pr8dvm6tsx8atjqqfe4w7nrpzgxzd2t5ytm25erqu0w4ku", // 用户的nft接收地址
                "orderStatus": "NOTPAID", // 订单状态, NOTPAID:未支付;PAYPENDING:支付确认中;PAYSUCCESS:支付成功;PAYTIMEOUT:超时未支付;INSCRIBING:铭刻交易等待确认中;ALLSUCCESS:订单成功
                "paytime": "", // 支付交易发起时间
                "payConfirmedTime": "", // 支付交易确认时间
                "revealTxid": "", // 铭文交易id
                "createTime": "2023-05-27 17:25:47 +0800 CST" // 订单生成时间
            },
            {
                "orderId": "BXHJSY54E7VLPTMTSL01252023052917044134049267",
                "eventId": 1,
                "depositAddress": "bc1pvlpt3h7ydtnmraw49wgvme2pcgyrp6ha379uj6xwdsz058kek0vq53mtsl",
                "total": 1123456,
                "receiveAddress": "bc1phjsyw73de6ap8nfjzg4erxmdw7lzlfgvm447v82fytn78nm0mwnsq654e7",
                "orderStatus": "NOTPAID",
                "paytime": "",
                "payConfirmedTime": "",
                "revealTxid": "",
                "createTime": "2023-05-27 17:27:23 +0800 CST"
            }
        ]
    }
    ```


##  `coinprice` 获取`BTC`价格 (每小时更新一次)

- 请求方式: `POST`

- 请求参数: 无

- 响应示例

    ```json
    {
        "code": 0,
        "msg": "ok",
        "data": {
            "btcPriceUsd": 27848   // BTC的价格（美元）
        }
    }
    ```


## **querygallerylist**查询图鉴列表(分页)

- 请求方式: `POST`

- 请求参数：

    | 字段 | 说明| 类型 | 可选? | 示例 |
    |-----|------|------|----|----|
    | `curPage` | 页号 | integer | 必填 | 0 |
    | `pageSize` |  页大小 | integer | 必填 | 100 |
    | `category` | 分类 bald,punk,rich,elite  | string | 必填 | bald |


- 请求示例:

    ```json
    {
        "curPage":0,
        "pageSize":100,
        "category":"bald"
    }
    ```

- 响应示例


    ```json
    {
        "code": 0,
        "msg": "ok",
        "data": {
            "category": "bald",
            "curPage": 0,
            "totalPage": 1,
            "pageSize": 100,
            "nfts": [
                {
                    "id": 1,
                    "name": "#1",
                    "description": "bitegale no1",
                    "imageUrl": "https://c-ssl.dtstatic.com/uploads/item/201504/16/20150416H4223_vG4eY.thumb.1000_0.jpeg"
                }
            ]
        }
    }
    ```