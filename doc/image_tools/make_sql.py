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

    bald = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_48/bald/*.png')
    for f in bald:
        name = os.path.basename(f)
        no = name[ :name.find(".png") ]
        g_category[no] = "bald"


    rich = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_48/rich/*.png')
    for f in rich:
        name = os.path.basename(f)
        no = name[ :name.find(".png") ]
        g_category[no] = "rich"

    elite = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_48/elite/*.png')
    for f in elite:
        name = os.path.basename(f)
        no = name[ :name.find(".png") ]
        g_category[no] = "elite"

    punk = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_48/punk/*.png')
    for f in punk:
        name = os.path.basename(f)
        no = name[ :name.find(".png") ]
        g_category[no] = "punk"
    pass


def main2():


    # infle,open('/home/yqq/firstsat/website/doc/image_tools/20230609_insert_blindbox.sql','w') as outfile:

    old_no_list = []

    with open(  "/home/yqq/firstsat/website/doc/image_tools/freemint_0609_2000.csv") as infile:
        lines = infile.readlines()
        for line in lines:
            line = str(line).strip()
            no = line[ :line.find(".png") ]
            id = int(no)
            old_no_list.append(id)



    # with open('/home/yqq/firstsat/website/doc/image_tools/20231121_sell_blindbox_url.sql','w') as outfile:
    with open('/home/yqq/firstsat/website/doc/image_tools/20231121_delete_elite_blindbox_url.sql','w') as outfile:
        for n in range(1, 10239):
            if n in old_no_list:
                continue

            name =  '#'+str(n)
            id = n
            no = str(n)
            description = 'Bit Eagle ' + name
            category = g_category[str(n)]

            if category != 'elite':
                continue


            image_url = f"https://d30f95b5opmmez.cloudfront.net/images/{id}.png" # TODO:
            # sql = f"INSERT INTO website.tb_blindbox (id, name, description, category, img_url, is_active, is_locked, status, commit_txid, reveal_txid, create_time, update_time) VALUES({id}, '{name}', '{description}', '{category}', '{image_url}', 1, 0, 'NOTMINT', NULL, NULL, '2023-06-06 12:03:13', '2023-06-06 12:03:13');"
            # sql = f"update website.tb_blindbox set img_url='{image_url}' where id={id};"
            sql = f"delete from website.tb_blindbox where id={id};"
            outfile.write( sql + "\n" )
    pass

# def main():
#     files1 = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_48/bald/*.png')
#     files2 = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_48/elite/*.png')
#     files3 = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_48/punk/*.png')
#     files4 = glob.glob('/home/yqq/firstsat/website/images/NFT_EAGLE_48/rich/*.png')

#     files = []
#     files.extend(files1)
#     files.extend(files2)
#     files.extend(files3)
#     files.extend(files4)


#     with open('./0609_blindbox.')
#     for f in files:
#         name = os.path.basename(f.strip())
#         no = name[ :name.find(".png") ]

#         id = int(no)
#         name =  '#'+no
#         description = 'Bit Eagle ' + name
#         category = g_category[no]
#         image_url = f"https://biteagle.io/images/{getname(str(id))}.png" # TODO:
#         sql = f"INSERT INTO website.tb_blindbox (id, name, description, category, img_url, is_active, is_locked, status, commit_txid, reveal_txid, create_time, update_time) VALUES({id}, '{name}', '{description}', '{category}', '{image_url}', 1, 0, 'NOTMINT', NULL, NULL, '2023-06-06 12:03:13', '2023-06-06 12:03:13');"
#         outfile.write( sql + "\n" )
#     pass


if __name__ == '__main__':
    load_category()
    # main()
    main2()
