#coding:utf8


import glob
import json
import os



def main():
    inscriptions = []
    with open('./tb_blindbox_202307241650.json', 'r') as jfile:
        jd = json.load(jfile)
        boxs = jd['tb_blindbox']
        for box in boxs:
            if 'reveal_txid' not in box :
                continue
            if box['reveal_txid'] is None:
                continue

            id = box['reveal_txid'] + 'i0'
            name = 'BitEagle ' + box['name']
            category = box['category'].upper()[:1] + box['category'][1:]
            item = {
                "id": f"{id}",
                "meta": {
                    "name": f"{name}",
                    "attributes": [
                        {
                            "value": f"{category}",
                            "trait_type": "type"
                        }
                    ]
                }
            }

            inscriptions.append(item)
            pass
        with open('inscription.json', 'w') as outfile:
            outfile.write( json.dumps(inscriptions, indent=2) )


    pass

if __name__ == '__main__':
    main()
