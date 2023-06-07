#coding:utf8

from hashlib import sha256
from binascii import hexlify
import glob
import os

def rename(id):
    salt = "qiyihuo"
    id = id
    h =  hexlify( sha256((salt + id).encode('latin')).digest()[5:15]).decode('latin')
    return h


def main():
    files = glob.glob("/home/yqq/firstsat/website/images/2023-06-07/rename_50x50_resize_compressed_2000/*.png")

    # print(rename(1))
    for file in files:
        name = os.path.basename(file)
        id = name.replace('.png', '')
        newname = rename(id) + '.png'
        newpath = os.path.join(os.path.dirname(file), newname)
        os.system(f"cp {file} {newpath}")



    pass


if __name__ == '__main__':
    main()
