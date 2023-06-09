#coding:utf8

import os



# def main():

#     with open( "/home/yqq/firstsat/website/doc/image_tools/freemint_0609_2000.csv") as infle:
#         lines = infle.readlines()
#         for line in lines:
#             line = str(line).strip()
#             no = line[ :line.find(".png") ]
#             srcfile = f'/home/yqq/firstsat/website/doc/image_tools/compressed/{no}-crunch.png')
#             newfile = os.path.join(srcdir,  f'2023-06-07/720x720_compressed_2000/{no}.png')
#             os.system(f"cp {srcfile}  {newfile}")

#     pass


# def main():

#     srcdir = "/home/yqq/firstsat/website/images/"
#     with open( os.path.join(srcdir,  "2023-06-07/freemint2000.csv")) as infle:
#         lines = infle.readlines()
#         for line in lines:
#             line = str(line).strip()
#             no = line[ :line.find(".png") ]
#             srcfile = os.path.join(srcdir,  f'raw_300x300/{no}.png')
#             newfile = os.path.join(srcdir,  f'2023-06-07/raw_300x300_2000/{no}.png')
#             os.system(f"cp {srcfile}  {newfile}")

#     pass

if __name__ == '__main__':
    main()

