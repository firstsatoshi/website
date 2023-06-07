#coding:utf8
from hashlib import sha256
from binascii import hexlify

import os
import glob

def getname(id):
    salt = "qiyihuo"
    id = id
    h =  hexlify( sha256((salt + id).encode('latin')).digest()[5:15]).decode('latin')
    return h

g_category = {}

def load_category():

    bald = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_720/bald/*.png')
    for f in bald:
        name = os.path.basename(f)
        no = name[ :name.find(".png") ]
        g_category[no] = "bald"


    rich = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_720/rich/*.png')
    for f in rich:
        name = os.path.basename(f)
        no = name[ :name.find(".png") ]
        g_category[no] = "rich"

    elite = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_720/elite/*.png')
    for f in elite:
        name = os.path.basename(f)
        no = name[ :name.find(".png") ]
        g_category[no] = "elite"

    punk = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_720/punk/*.png')
    for f in punk:
        name = os.path.basename(f)
        no = name[ :name.find(".png") ]
        g_category[no] = "punk"
    pass


def main():


    srcdir = "/home/yqq/firstsat/website/images/"
    with open( os.path.join(srcdir,  "2023-06-07/freemint2000.csv")) as infle, open('/home/yqq/firstsat/website/doc/20230607_insert_blindbox.sql','w') as outfile:
        lines = infle.readlines()
        for line in lines:
            line = str(line).strip()
            no = line[ :line.find(".png") ]

            id = int(no)
            name =  '#'+no
            description = 'Bit Eagle ' + name
            category = g_category[no]
            image_url = f"https://biteagle.io/images/{getname(str(id))}.png" # TODO:
            sql = f"INSERT INTO website.tb_blindbox (id, name, description, category, img_url, is_active, is_locked, status, commit_txid, reveal_txid, create_time, update_time) VALUES({id}, '{name}', '{description}', '{category}', '{image_url}', 1, 0, 'NOTMINT', NULL, NULL, '2023-06-06 12:03:13', '2023-06-06 12:03:13');"
            outfile.write( sql + "\n" )


    pass

if __name__ == '__main__':
    load_category()
    main()
