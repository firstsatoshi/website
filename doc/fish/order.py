#coding:utf8



def main():

    mint_order_map = {}
    with open('./order.txt', 'r') as infile:
        lines = infile.readlines()
        for line in lines:
            l = line.strip()
            if len(l) == 0: continue
            ls = l.split(',')
            address = ls[0].strip()
            count = int(ls[1].strip())

            if address not in mint_order_map:
                mint_order_map[address] = 0
            mint_order_map[address] += int(count)

    whitelist_count_map = {}
    with open('./whitelist1.txt', 'r') as infile:
        lines = infile.readlines()
        for line in lines:
            l = line.strip()
            if len(l) == 0: continue
            ls = l.split(',')
            address = ls[0].strip()
            count = int(ls[1].strip())
            if address not in whitelist_count_map:
                whitelist_count_map[address] = 0
            whitelist_count_map[address] += int(count)

    addr_name_map = {}
    with open('address_name.txt', 'r') as infile:
        lines = infile.readlines()
        for line in lines:
            l = line.strip()
            if len(l) == 0: continue
            ls = l.split(',')
            address = ls[0].strip()
            name = ls[1].strip()
            addr_name_map[address] = name

    name_telno = {}
    with open('name_telno.txt', 'r') as infile:
        lines = infile.readlines()
        for line in lines:
            l = line.strip()
            if len(l) == 0: continue
            ls = l.split(',')
            telno = ls[0].strip()
            name = ls[1].strip()
            name_telno[name] = telno



    with open('bitcoinfish_data_0809.xlsx', 'w') as ofile:
        ofile.write('用户姓名\t手机号\t地址\t白名单数量\t已mint数量\t完成?\n')
        for user_addr, count in whitelist_count_map.items():
            s = ''

            sep = '\t'
            # 用户姓名
            if user_addr in addr_name_map:
                s += addr_name_map[user_addr] + sep

                # 手机号
                name = addr_name_map[user_addr]
                if name in name_telno:
                    s += name_telno[name] +  sep
                else:
                    s += '    ' +  sep
            else:
                s += 'xxx' +  sep
                s += 'xxx' +  sep




            s += user_addr +  sep
            # 白名单登记的数量
            s += str(count) +  sep

            # 实际mint的数量
            mint_count = 0
            if user_addr in mint_order_map:
                s += str(mint_order_map[user_addr])  +  sep
                mint_count = mint_order_map[user_addr]
            else:
                s += '0' + sep
                mint_count = 0

            if mint_count != count:
                s += 'NO'
            else:
                s += '  '

            s += '\n'

            ofile.write(s)


    pass

if __name__ == '__main__':
    main()
