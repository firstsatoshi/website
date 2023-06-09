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



# def thumbnail():
#     # filename = "/home/yqq/firstsat/website/images/2023-06-07/50x50_resize_compressed_2000/569.png"
#     # filename = "/home/yqq/下载/ceshi/ceshi/10003.png"
#     # filename = "/home/yqq/firstsat/website/images/NFT_EAGLE_48/elite/10003.png"
#     try:
#         # Load just once, then successively scale down
#         im=Image.open(filename)
#         im.thumbnail((50,50))
#         im.save("./{}".format(os.path.basename(filename)))
#         return 'OK'
#     except Exception as e:
#         return e
# files = glob.glob('/home/yqq/firstsat/website/images/2023-06-07/720x720_compressed_2000/*.png')
# print(len(files))
# pool = Pool(16)
# results = pool.map(thumbnail, zip(files,range((len(files)))))


def thumbnail(params):
    filename, N = params

    try:
        # Load just once, then successively scale down
        im=Image.open(filename)
        # nim = im.resize((138, 138))
        # nim.save("/home/yqq/firstsat/website/images/138x138_gallery/{}".format(os.path.basename(filename)))
        # nim = im.resize((276, 276))
        # nim.save("/home/yqq/firstsat/website/images/276x276_gallery/{}".format(os.path.basename(filename)))
        nim = im.resize((1104, 1104))
        nim.save("/home/yqq/firstsat/website/images/1104x1104_gallery/{}".format(os.path.basename(filename)))
        im.close()
        return 'OK'
    except Exception as e:
        return e


files1 = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_48/bald/*.png')
files2 = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_48/elite/*.png')
files3 = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_48/punk/*.png')
files4 = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_48/rich/*.png')

files = []
files.extend(files1)
files.extend(files2)
files.extend(files3)
files.extend(files4)

print(len(files))
pool = Pool(16)
results = pool.map(thumbnail, zip(files,range((len(files)))))