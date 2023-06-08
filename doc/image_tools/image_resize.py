#!/usr/local/bin/python3

import glob
from PIL import Image
from multiprocessing import Pool
import os

# def thumbnail(params):
#     filename, N = params

#     try:
#         # Load just once, then successively scale down
#         im=Image.open(filename)
#         im.thumbnail((50,50))
#         im.save("/home/yqq/firstsat/website/images/2023-06-07/50x50_resize_2000/{}".format(os.path.basename(filename)))
#         return 'OK'
#     except Exception as e:
#         return e
# files = glob.glob('/home/yqq/firstsat/website/images/2023-06-07/720x720_compressed_2000/*.png')
# print(len(files))
# pool = Pool(16)
# results = pool.map(thumbnail, zip(files,range((len(files)))))



def thumbnail():
    # filename = "/home/yqq/firstsat/website/images/2023-06-07/50x50_resize_compressed_2000/569.png"
    # filename = "/home/yqq/下载/ceshi/ceshi/10003.png"
    filename = "/home/yqq/firstsat/website/images/NFT_EAGLE_720/elite/10003.png"
    try:
        # Load just once, then successively scale down
        im=Image.open(filename)
        im.thumbnail((50,50))
        im.save("./{}".format(os.path.basename(filename)))
        return 'OK'
    except Exception as e:
        return e
# files = glob.glob('/home/yqq/firstsat/website/images/2023-06-07/720x720_compressed_2000/*.png')
# print(len(files))
# pool = Pool(16)
# results = pool.map(thumbnail, zip(files,range((len(files)))))


# def thumbnail(params):
#     filename, N = params

#     try:
#         # Load just once, then successively scale down
#         im=Image.open(filename)
#         im.thumbnail((200,200))
#         im.save("/home/yqq/firstsat/website/images/2023-06-07/200x200_resize_compressed_2000/{}".format(os.path.basename(filename)))
#         return 'OK'
#     except Exception as e:
#         return e


# files = glob.glob('/home/yqq/firstsat/website/images/2023-06-07/720x720_compressed_2000/*.png')
# print(len(files))
# pool = Pool(16)
# results = pool.map(thumbnail, zip(files,range((len(files)))))