#coding:utf8

import hashlib

import json
import base64

fish_items = json.load( open('fish_items.json', 'r') )

items_type = {
    "bj": "Background",
    "bb": "Back",
    "f": "Fish",
    "yf": "Clothes",
    "tb": "Head",
    "lb": "Face",
    "yj": "Eyes",
    "ps": "Accessories",
    "zc":"Mouth"
}

def parse_merge_path(merge_path: str):

    print(merge_path)
    ps = merge_path.split('_')

    attributes = []

    for p in ps:

        item_type = items_type[ p.strip()[2:] ]

        value = fish_items[p].strip()
        attributes.append({
            "trait_type": f"{item_type}",
            "value": f"{value}",
        })


    return attributes


def main():

    # s = base64.b64decode('OTlial8wNmZfMTJwc18wM3pj').decode('latin')
    # r = parse_merge_path(s)
    # print(r)

    inscriptions = []

    name_id = 1
    with open('fish_969.txt', 'r') as infile:
        lines = infile.readlines()
        for line in lines:
            l = line.strip()
            if len(l) == 0: continue
            ls = l.split('\t')
            inscribe_id = ls[0].strip()

            fish_path = ls[1].strip()
            fish_path = fish_path.replace('bitcoinfish_', '')
            fish_path = fish_path.replace('.html', '')
            merge_path = base64.b64decode(fish_path).decode('latin')
            attributes = parse_merge_path(merge_path=merge_path)


            name = 'BicoinFish #' + str(name_id)
            item = {
                "id": inscribe_id + "i0",
                "meta": {
                    "name": name,
                    "attributes": attributes
                }
            }

            inscriptions.append(item)
            name_id += 1


    with open('inscription.json', 'w') as outfile:
            outfile.write( json.dumps(inscriptions, indent=2) )


    pass

if __name__ == '__main__':
    main()
