
# 解决brc20和域名查重的问题


## 问题描述

在brc20和域名的注册过程中，需要对brc20和域名进行查重，如果重复则不能注册。

但是这种方式有以下问题：

- 目前没有提供基础数据的服务商（API），需要我们自行对ordinals的数据进行索引index,并建立相应的数据库，这个工作量比较大，需要的服务器较多。
- 由于brc20和域名的注册是异步的，所以在查重的时候，可能会出现重复的情况。


## 解决办法


- 方案一： 通过调用unisat的api实现查重功能（虽然unisat没有对外提供api）
- 方案二： 自建数据服务


## 决定：

暂时用 方案一快速实现我们的功能。

看后续btc生态的发展情况，再决定是否自建数据服务。最好的方式是用第三方数据服务商的api,降低我们的成本。


- 查询brc20：
  - https://unisat.io/brc20-api-v2/brc20/ordi/info
  - https://unisat.io/brc20-api-v2/inscriptions/category/sats/search/v2?name=hello&limit=32&start=0

- sats: https://unisat.io/brc20-api-v2/inscriptions/category/sats/existence



