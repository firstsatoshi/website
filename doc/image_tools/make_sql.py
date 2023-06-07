#coding:utf8

import os



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
            category = 'bald' # TODO:
            image_url = "" # TODO:
            sql = f"INSERT INTO website.tb_blindbox (id, name, description, category, img_url, is_active, is_locked, status, commit_txid, reveal_txid, create_time, update_time) VALUES({id}, '{name}', '{description}', '{category}', '{image_url}', 1, 0, 'NOTMINT', NULL, NULL, '2023-06-06 12:03:13', '2023-06-06 12:03:13');"
            outfile.write( sql + "\n" )


    pass

if __name__ == '__main__':
    main()
