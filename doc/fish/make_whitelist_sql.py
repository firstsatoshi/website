#coding:utf8



def main():
    addrs = {}
    with open('./whitelist1.txt', 'r') as infile,  open('whitelist1.sql', 'w') as ofile:
        lines = infile.readlines()
        for line in lines:
            l = line.strip()
            ls = l.split('\t')
            address = ls[0]
            count = int(ls[1])
            if address not in addrs:
                addrs[address] = 0
            addrs[address] += int(count)

        for addr, count in addrs.items():

            email = addr[5:11] + "@whitelist.com"
            s = f"""INSERT INTO website.tb_waitlist (event_id, email, btc_address, referee_id, mint_limit) VALUES(1, '{email}', '{addr}', 0,  {count});"""
            # print(s)
            ofile.write(s  + '\n')

    with open('whitelist.csv', 'w') as ofile:
        for a, c in addrs.items():
            ofile.write(f'{a},{c}\n')

    pass

if __name__ == '__main__':
    main()
