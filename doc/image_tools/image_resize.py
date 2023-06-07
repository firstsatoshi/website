#!/usr/local/bin/python3

import glob
from PIL import Image
from multiprocessing import Pool

def thumbnail(params):
    filename, N = params
    try:
        # Load just once, then successively scale down
        im=Image.open(filename)
        im.thumbnail((1600,1600))
        im.save("t-%d-1600.jpg" % (N))
        im.thumbnail((720,720))
        im.save("t-%d-720.jpg"  % (N))
        im.thumbnail((120,120))
        im.save("t-%d-120.jpg"  % (N))
        return 'OK'
    except Exception as e:
        return e


files = glob.glob('image*.jpg')
pool = Pool(8)
results = pool.map(thumbnail, zip(files,range((len(files)))))