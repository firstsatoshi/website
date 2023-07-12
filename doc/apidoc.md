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
    | `referalCode` | 邮箱 | string | **可选** | HXsg |

- 请求示例：

    ```json
    {
        "email":"youngqqcn@163.com",
        "btcAddress":"bc1pv5d2mmx2v9cx9menxl5zlhacljqu9zqhltl4d303n3rjjcxfrgwqa80as0"
        "referalCode":"2342"
    }
    ```


- 响应示例
  ```json
    {
        "code": 0,
        "msg": "ok",
        "data": {
            "referalCode": "HXsg", // 该用户的邀请码
            "duplicated": true   // 是否重复预约，
        }
    }
  ```

每个邮箱只能预约一次


## `queryblindboxevent` 查询盲盒活动详情

- 请求方式: POST

- 请求参数:

    | 字段 | 说明| 类型 | 可选? | 示例 |
    |-----|------|------|----|----|
    | `eventId` | 活动id | int | 可选(如不填，则查询所有) | 1 |
    | `status` | 活动状态, `active`, `inactive` | string | 可选(如不填，则查询所有) | active |
    | `eventEndpoint` | 活动名(唯一), 例如: fsat.io/collection/biteagle  | string | 可选(如不填, 根据其他字段查询) | biteagle |


- 请求示例：

```json
{
    "status":"active"
}
```

- 响应示例:

    ```json
    {
        "code": 0,
        "msg": "ok",
        "data": [
            {
                "eventId": 1, // 活动id
                "name": "Azuki", // 盲盒活动名
                "eventEndpoint": "azuki",
                "description": "Azuki is Azuki", // 盲盒活动描述(支持富文本)
                "detail": "Azuki Azuki Azuki Azuki Azuki Azuki Azuki Azuki Azuki Azuki Azuki Azuki Azuki Azuki Azuki Azuki Azuki Azuki Azuki Azuki", // 盲盒活动详情(支持富文本)
                "avatarImageUrl": "https://azk.imgix.net/images/166a4190-1377-4cc2-9b93-e1cd32a72ac1.png", // 头像
                "backgroundImageUrl": "https://azk.imgix.net/images/166a4190-1377-4cc2-9b93-e1cd32a72ac1.png", // 盲盒活动背景图
                "roadmapDescription": "Azuki is Azuki", // 路线图描述(支持富文本)
                "roadmapList": [ // 路线图列表
                    "freemint",
                    "vip sale",
                    "auction sale",
                    "public sale"
                ],
                "websiteUrl": "https://www.azuki.com/zh",
                "twitterUrl": "https://twitter.com/Azuki",
                "discordUrl": "http://discord.gg/azuki",
                "imagesList": [ // banner页的图片
                    "https://azk.imgix.net/images/166a4190-1377-4cc2-9b93-e1cd32a72ac1.png",
                    "https://azk.imgix.net/images/57451d86-b9d2-4199-a91f-fcdfba68d2cf.png",
                    "https://azk.imgix.net/images/57451d86-b9d2-4199-a91f-fcdfba68d2cf.png",
                    "https://azk.imgix.net/images/57451d86-b9d2-4199-a91f-fcdfba68d2cf.png"
                ],
                "currentMintPlanIndex":0, // 当前盲盒计划的索引, 即mintPlanList的索引
                "mintPlanList": [
                    {
                        "title": "Freemint",
                        "plan": "June 29 at 11:00 PM UTC+0"
                    },
                    {
                        "title": "Public Sale",
                        "plan": "2023-07-11 20:00 UTC+0"
                    }
                ],
                "priceBtcSats": 1000, // 盲盒价格， 单位是聪(satoshi),如果要换算成BTC要除以10^8, 例如：123456 satoshi = 0.00123456BTC
                "priceUsd": 0, // 盲盒价格（美元），
                "paymentToken": "BTC", // 收款币种，用户必须使用此币种进行支付
                "averageImageBytes": 550, // NFT图片的平均大小（字节数）
                "supply": 1000, // 总供应量（本次活动供应总量）
                "avail": 500,  // 当前可用量(背刺活动当前可用库存)
                "mintLimit": 10,  // 限购数
                "isActive": true, // 活动是否开启
                "isDisplay": true, // 是否显示
                "onlyWhitelist": true, // 是否仅对白名单用户开放
                "startTime": 1688376539,  // 活动开始时间
                "endTime": 1691054939 // 活动结束时间
            },
            {
                "eventId": 2,
                "name": "Punks",
                "eventEndpoint": "cryptopunks",
                "description": "Punks is CryptoPunks",
                "detail": "Punks Punks Punks Punks Punks Punks Punks Punks Punks Punks Punks Punks Punks Punks Punks Punks Punks Punks Punks Punks",
                "avatarImageUrl": "https://bitcoinpunks.com/punks-bg/punk0070.png",
                "backgroundImageUrl": "https://bitcoinpunks.com/punks-bg/punk9085.png",
                "roadmapDescription": "Punks is CryptoPunks",
                "roadmapList": [
                    "freemint",
                    "vip sale",
                    "auction sale",
                    "public sale"
                ],
                "websiteUrl": "https://cryptopunks.app/",
                "twitterUrl": "https://twitter.com/cryptopunksnfts",
                "discordUrl": "https://t.co/zFvNWfBy6C",
                "imagesList": [
                    "https://bitcoinpunks.com/punks-bg/punk9085.png",
                    "https://bitcoinpunks.com/punks-bg/punk0070.png",
                    "https://bitcoinpunks.com/punks-bg/punk7223.png",
                    "https://bitcoinpunks.com/punks-bg/punk0482.png"
                ],
                "currentMintPlanIndex":0, // 当前盲盒计划的索引, 即mintPlanList的索引
                "mintPlanList": [
                    {
                        "title": "Freemint",
                        "plan": "June 29 at 11:00 PM UTC+0"
                    },
                    {
                        "title": "Public Sale",
                        "plan": "2023-07-11 20:00 UTC+0"
                    }
                ],
                "priceBtcSats": 1024,
                "priceUsd": 0,
                "paymentToken": "BTC",
                "averageImageBytes": 550,
                "supply": 500,
                "avail": 100,
                "mintLimit": 10,
                "isActive": true, // 活动是否开启
                "isDisplay": true, // 是否显示
                "onlyWhitelist": false,
                "startTime": 1688377409,
                "endTime": 1693734209
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

## `createinscribeorder`创建铭刻订单

- 请求方式: POST

- 请求参数:

    | 字段 | 说明| 类型 | 可选? | 示例 |
    |-----|------|------|----|----|
    | `fileUploads` | 上传文件列表 | [string] | 必填 |  |
    | `receiveAddress` | 接收地址 | string | 必填 |  |
    | `feeRate` | 费率 | integer | 必填 | |
    | `checksum` | 校验和 | string| 必填 | |
    | `token` | 验证器的token | string| 必填 | |


- 请求示例：

    ```json
    {
        "fileUploads":[
            {
                "dataUrl":"data:text/plain;charset=utf-8;base64,aGV5YQ==",  // 文件内容的 dataUrl 编码 https://www.rfc-editor.org/rfc/rfc2397
                "fileName":"1111"
            },
            {
                "dataUrl":"data:text/plain;charset=utf-8;base64,eyJwIjoiYnJjLTIwIiwib3AiOiJkZXBsb3kiLCJ0aWNrIjoieXFxcSIsIm1heCI6IjIxMDAwMDAwIiwibGltIjoiMTAwMCJ9Cg==",
                "fileName":"yqqq"
            }
        ],
        "receiveAddress":"tb1phjsyw73de6ap8nfjzg4erxmdw7lzlfgvm447v82fytn78nm0mwnshjz6r3",
        "feeRate":3,
        "checksum":"xx",
        "token": "xxx"
    }
    ```

- 响应示例:


```json
{
    "code": 0,
    "msg": "ok",
    "data": {
        "orderId": "ITHJSYZ6R3M729FY5S02032023070319122727388763",
        "count": 2,
        "filenames": [
            "1111",
            "yqqq"
        ],
        "depositAddress": "tb1pm729vkl73z0679wr8pztg09eq9ed0ynyr0g3qcz8ylxu885nzhcqrkfy5s",
        "receiveAddress": "tb1phjsyw73de6ap8nfjzg4erxmdw7lzlfgvm447v82fytn78nm0mwnshjz6r3",
        "feeRate": 3,
        "bytes": 0,
        "inscribeFee": 7092,
        "serviceFee": 0,
        "total": 7092,
        "createTime": 1688382747
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
        "receiveAddress":"tb1pv5d2mmx2v9cx9menxl5zlhacljqu9zqhltl4d303n3rjjcxfrgwq20ej2q"
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
                "depositAddress": "tb1pvnphdaypksz6allfdm0e7638uwt2a0tjvtssgser1234324324234", // 充值地址
                "total": 1123456, // 总金额（单位，聪），
                "receiveAddress": "tb1pv5d2mmx2v9cx9menxl5zlhacljqu9zqhltl4d303n3rjjcxfrgwq20ej2q", // 用户的nft接收地址
                "orderStatus": "NOTPAID", // 订单状态, NOTPAID:未支付;PAYPENDING:支付确认中;PAYSUCCESS:支付成功;PAYTIMEOUT:超时未支付;MINTING:铭刻交易等待确认中;ALLSUCCESS:订单成功
                "paytime": "", // 支付交易发起时间
                "payConfirmedTime": "", // 支付交易确认时间
                "revealTxids": [], // 铭文交易id数组
                "createTime": 1688382023 // 订单生成时间
            },
            {
                "orderId": "BEMV5D2EJ2Q0SX3WGP701042023060614215940101518",
                "eventId": 1,
                "depositAddress": "",
                "total": 3546,
                "receiveAddress": "tb1pv5d2mmx2v9cx9menxl5zlhacljqu9zqhltl4d303n3rjjcxfrgwq20ej2q",
                "orderStatus": "MINTING",
                "payTxid": "d9b91d743e2729db532360ad5bfde8b065b97a897cc1040945cec10b834645e1",
                "paytime": "2023-06-06 14:28:48 +0800 CST",
                "payConfirmedTime": "2023-06-06 14:28:48 +0800 CST",
                "nftDetails": [
                    {
                        "txid": "5c4fef5b8fd8d716edd482ca42195edf43788ad5598c9b9cd63cf257407202c6",
                        "confirmed": false,
                        "name": "#6",
                        "category": "bald",
                        "description": "no6",
                        "imageUrl": "https://c-ssl.dtstatic.com/uploads/item/201504/16/20150416H4223_vG4eY.thumb.1000_0.jpeg",
                        "contentType": "text/plain",  // 铭文类型
                        "fileName": "hello,world,world",  // 文件名
                        "inscription": "aGVsbG8sd29ybGQsd29ybGQ"  // 铭文内容
                    }
                ],
                "createTime": 1688382023
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

##  **checkwhitelist**检查白名单

- 请求方式: `POST`

- 请求参数：

    | 字段 | 说明| 类型 | 可选? | 示例 |
    |-----|------|------|----|----|
    | eventId | 活动id  | integer | 必填 | 1 |
    | `receiveAddress` | 接收地址  | string | 必填 | tb1pv5d2mmx2v9cx9menxl5zlhacljqu9zqhltl4d303n3rjjcxfrgwq20ej2q |


- 请求示例:

    ```json
    {
        "eventId":1,
        "receiveAddress":"xxxxxxxxxxxxxx"
    }
    ```

- 响应示例

    ```json
    {
        "code": 0,
        "msg": "ok",
        "data": {
            "eventId": 1,
            "isWhitelist": false    // true： 是白名单  ; false 不是白名单
        }
    }
    ```

##  **queryaddress**检查当前地址订单数量


- 请求方式: `POST`

- 请求参数：

    | 字段 | 说明| 类型 | 可选? | 示例 |
    |-----|------|------|----|----|
    | `eventId` | 活动id | integer | 必填 |  1 |
    | `receiveAddress` | 接收地址  | string | 必填 | tb1pv5d2mmx2v9cx9menxl5zlhacljqu9zqhltl4d303n3rjjcxfrgwq20ej2q |


- 请求示例:

    ```json
    {
        "eventId":1,
        "receiveAddress":"tb1pj2rglj70gwfydjed9w8l4n5vcurxdpahpcj02076wdykz376376q8l8w3h"
    }
    ```

- 响应示例


    ```json
    {
        "code": 0,
        "msg": "ok",
        "data": {
            "eventId": 1,
            "isWhitelist": true,
            "eventMintLimit": 10,
            "currentOrdersTotal": 3
        }
    }
    ```

## **querybrc20**查询BRC20信息


- 请求方式: `POST`

- 请求参数:
    | 字段 | 说明| 类型 | 可选? | 示例 |
    |-----|------|------|----|----|
    | `ticker` | BRC20代币符号 | string | 必填 |  `ordi` |

- 响应参数：

    | 字段 | 说明| 类型 | 可选? | 示例 |
    |-----|------|------|----|----|
    | `isExists` | BRC20是否已经存在（部署） | bool | 必填 |   |
    | `ticker` | BRC20代币符号  | string | 必填 | `ordi`|
    | `limit` | 每次mint最大限制量  | integer | 必填 | `10`|
    | `max` | BRC20总供应量 | integer | 必填 | `2100000`|
    | `minted` | 当前已经mint的数量| integer | 必填 | `2100000`|
    | `decimal` | 精度| integer | 可选 | `18`|

- 响应示例：

  - 当BRC20已经部署时:
    ```json
    {
        "code": 0,
        "msg": "ok",
        "data": {
            "isExists": true,
            "ticker": "ordi",
            "limit": 1000,
            "max": 21000000,
            "minted": 21000000,
            "decimal": 18,
            "inscriptionId": "b61b0172d95e266c18aea0c624db987e971a5d6d4ebc2aaed85da4642d635735i0"
        }
    }
    ```

  - 当BRC20尚未部署时:
    ```json
    {
        "code": 0,
        "msg": "ok",
        "data": {
            "isExists": false,
            "ticker": "xfkz",
            "limit": 0,
            "max": 0,
            "minted": 0,
            "decimal": 0,
            "inscriptionId": ""
        }
    }
    ```


## `checknames` 检查域名是否已经注册过

**!!注意：前端必须检验`.sats`域名规则**： [点击这里](./tech/sats域名.md)

- 请求方式: `POST`

- 请求参数:
    | 字段 | 说明| 类型 | 可选? | 示例 |
    |-----|------|------|----|----|
    | `type` | 域名类型, 目前只有: `sats` | string | 必填 |  `sats` |
    | `names` | 域名列表 | [string] | 必填 |  `["aaa.sats", "bbb.sats"]` |

- 请求示例:

```json
{
    "type":"sats",
    "names":["aaa.sats","bbb.sats","88888888.sats", "mars.sats"]
}
```


- 响应示例：


```json
{
    "code": 0,
    "msg": "ok",
    "data": [
        {
            "name": "88888888.sats",
            "isExists": true   // true: 已被注册; false: 尚未被注册
        },
        {
            "name": "aaa.sats",
            "isExists": true
        },
        {
            "name": "bbb.sats",
            "isExists": true
        },
        {
            "name": "mars.sats",
            "isExists": true
        }
    ]
}
```



## `estimatefee` 估算铭文手续费

- 请求方式: `POST`

- 请求参数:

    | 字段 | 说明| 类型 | 可选? | 示例 |
    |-----|------|------|----|----|
    | `fileUploads` | 上传文件列表 | [string] | 必填 |  |
    | `feeRate` | 费率 | integer | 必填 | |


- 请求示例:

```json
{
    "fileUploads":[
        {
            "dataUrl":"data:text/plain;charset=utf-8;base64,aGV5YQ==",
            "fileName":"1111"
        },
        {
            "dataUrl":"data:text/plain;charset=utf-8;base64,eyJwIjoiYnJjLTIwIiwib3AiOiJkZXBsb3kiLCJ0aWNrIjoieXFxcSIsIm1heCI6IjIxMDAwMDAwIiwibGltIjoiMTAwMCJ9Cg==",
            "fileName":"yqqq"
        }
    ],
    "feeRate":3
}
```


- 响应示例：

```json
{
    "code": 0,
    "msg": "ok",
    "data": {
        "totalFee": 2364
    }
}
```