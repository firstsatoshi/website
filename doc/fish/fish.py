#coding:utf8

import glob
import os
import hashlib
import requests
import json

# 编写一个函数读取一个目录下面的图片文件，返回图片字节，图片名称，图片的sha256哈希

def load_imgs_map():
    files = glob.glob('/home/yqq/下载/加密咸鱼/xmeta/*.png')
    ret = {}
    for file in files:
        # print(file)
        filename = os.path.basename(file)
        filename = filename.replace('.png', '')
        # print(filename)
        with open(file, 'rb') as img:
            d = img.read(100000)
            h = hashlib.sha256(d).hexdigest()
            ret[h] = [filename, d]
            # print(h)

    return ret

def main2():
    m = load_imgs_map()


    o = {}
    with open('./fish_inscription_fish.json', 'r') as jfile:
        iscriptions = json.load(jfile)
        l = iscriptions['list']
        skip = True
        for item in l:
            content_url = item['content']
            s = content_url
            inscription_id = s[s.rfind('/') + 1:].strip()

            if skip :
                if inscription_id == '565c7b4f658df51540dd51635de13121436b6a00caa776cde8ca6684a5d492f0i0':
                    skip = False
                continue

            r = requests.get(content_url, timeout=60)

            h = hashlib.sha256(r.content).hexdigest()
            if h in m:
                o[m[h][0]] = inscription_id
                print( "{}, {}".format(m[h][0] , inscription_id) )
            else:
                print("nonno========")
            pass
        # with open('items_dict.json', 'w') as output:
            # output.write( json.dumps(o) )


def main3():

    o = {}

    with open('./result.txt', 'r') as infile:
        lines = infile.readlines()
        for line in lines:
            ss = line.split(',')
            name = ss[0].strip()
            inscription_id = ss[1].strip()

            item_type = name[2:]
            if item_type not in o:
                o[item_type] = {}
            o[item_type][name] = inscription_id

        with open('bitfish_items_dict.json', 'w') as outfile:
            outfile.write( json.dumps(o, indent=4) )

    pass


def main4():
    files = glob.glob('/home/yqq/下载/fff/加密咸鱼/加密咸鱼/*.png')
    ret = []
    retjson = {}
    for file in files:
        filename = os.path.basename(file)
        filename = filename.replace('.png', '')

        merge_path = ''
        # print(filename)
        merge_path += filename[:4]  + '_'

        x = 0
        if filename[6:7] == 'f' :
            merge_path += filename[4:7] + '_'
            x = 7
        elif filename[6:8] == 'bb':
            merge_path += filename[4:8] + '_'
            merge_path += filename[8:11] + '_'
            x = 11


        while True:
            if x >= len(filename):
                break
            start_index = x
            end_index = start_index + 4
            merge_path += filename[  start_index :  end_index] + '_'

            x = end_index
            pass

        # remove the last '_'
        merge_path = merge_path[:-1]
        # print(merge_path)
        ret.append(merge_path)
        retjson[merge_path] = True

    with open('merge_path.sql', 'w') as merge_path:
        for path in ret:
            s = f"""INSERT INTO website.tb_bitfish_merge_path (merge_path) VALUES('{path}');"""
            merge_path.write(s + '\n')

    with open('bitfish_merged_paths.json', 'w') as outfile:
        outfile.write( json.dumps(retjson, indent=4) )

    pass

def fac(n):
    """ 计算阶乘 n! """
    if n == 0 or n == 1:
        return 1
    return n * fac(n-1)

def C(m, n):
    """从m个元素中取n个,组合公式 C(m, n) = m! / (n! * (m - n)!)"""
    return fac(m) / (fac(n) *fac(m - n))

def calc_bitfish_total():
    """计算bitfish总数"""
    items = [6 , 10 , 9 , 9 , 20 , 17 , 9]
    m = sum(items)
    total = 0
    for n in range(1, 8):
        total += C(m, n)
    print(int(total))

if __name__ == '__main__':
    calc_bitfish_total()




