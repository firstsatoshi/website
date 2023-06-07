#coding:utf8

import os



def main():

    srcdir = "/home/yqq/firstsat/website/images/"
    with open( os.path.join(srcdir,  "2023-06-07/freemint2000.csv")) as infle:
        lines = infle.readlines()
        for line in lines:
            line = str(line).strip()
            no = line[ :line.find(".png") ]
            srcfile = os.path.join(srcdir,  f'compressed/{no}-crunch.png')
            newfile = os.path.join(srcdir,  f'2023-06-07/imgs/{no}.png')
            os.system(f"cp {srcfile}  {newfile}")

    pass

if __name__ == '__main__':
    main()

