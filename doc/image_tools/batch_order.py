#coding:utf8

import json
import requests

token = "youngqqcn@gmail.com"

# 实现一个post接口请求，打印请求结果，请求参数为json格式,接口url为 https://www.biteagle.io/api/v1/createorder
#  请求参数示例 如下：
# {
#     "eventId":1,
#     "count": 2,
#     "receiveAddress":"tb1pv5d2mmx2v9cx9menxl5zlhacljqu9zqhltl4d303n3rjjcxfrgwq20ej2q",
#     "feeRate":3,
#     "token":"dffffffffffffffffff"
# }
def createorder(recv_addr):
    url = "http://18.181.187.89/api/v1/createorder"
    data = {
        "eventId": 1,
        "count": 1,
        "receiveAddress": recv_addr,
        "feeRate": 11,
        "token": token
    }
    headers = {'content-type': 'application/json'}
    response = requests.post(url, json=data, headers=headers)
    print(response.text)
    if response.status_code != 200:
        return None, None

    resp = json.loads(response.text)
    return resp['data']['depositAddress'], resp['data']['total']



# 实现一个函数获取一个文本文件的内容，并将每一行的内容存储到一个列表中，返回该列表
def read_file():
    ret = []
    with open("morin_b_100.txt", "r") as f:
        lines = f.readlines()
        for line in lines:
            line = line.strip()
            ret.append(line)
    return ret


if __name__ == '__main__':
    addrs = read_file()
    with open('b_100.csv', 'a') as outfile:
        for addr in addrs:
            a, t = createorder(addr)
            if a is None or t is None:
                continue
            outfile.write('%s,%s\n' % (a, t))
