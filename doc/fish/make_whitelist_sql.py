#coding:utf8



def main():
    addrs = {}

    address_map = {}
    with open('./whitelist1.txt', 'r') as infile,  open('whitelist.sql', 'w') as ofile:
        lines = infile.readlines()
        for line in lines:
            l = line.strip()
            ls = l.split('\t')
            address = ls[0]
            count = ls[1]

            if address not in address_map:
                address_map[address] = [line ]
            else:
                address_map[address].append(line)

            if address not in addrs:
                addrs[address] = 0
                # print(f"=====duplicated address : {address}")
            addrs[address] += 1

            email = address[5:11] + "@whitelist.com"
            s = f"""INSERT INTO website.tb_waitlist (event_id, email, btc_address, referee_id, mint_limit) VALUES(1, '{email}', '{address}', 0,  {count});"""
            # print(s)
            ofile.write(s)

            # print(s)
        for k, v in addrs.items():
            if v > 1:
                print(f'{k} 重复. 一共 {v}  次')
                # print(address_map[address])


    pass

if __name__ == '__main__':
    main()
