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

if __name__ == '__main__':
    # main()
    # main2()
    main3()




