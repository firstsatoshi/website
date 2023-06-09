
from PIL import Image

im = Image.open("/home/yqq/firstsat/website/images/NFT_EAGLE_48/bald/1.png")
# im.thumbnail((138,138))
nim = im.resize((138, 138))
# nim.save('new_test.png')
rgb_im = nim.convert('RGB')
rgb_im.save('new_test.jpg')
