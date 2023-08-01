#coding:utf8

def fac(n):
    """ 计算阶乘 n! """
    if n == 0 or n == 1:
        return 1
    return n * fac(n-1)

def C(m, n):
    """
    计算从m个元素中取n个所有可能的组合数, 根据组合公式 C(m, n) = m! / (n! * (m - n)!)
    """
    return fac(m) / (fac(n) *fac(m - n))

def counter_xiwu_n_elements(n):
    """计算稀物n元素鱼的个数"""
    count = 0
    with open('bitfish_merged_paths.json', 'r') as infile:
        lines = infile.readlines()
        for line in lines:
            if line.count('_') == n - 1:
                count += 1
    return count

def calc_bitcoinfish_total(n_elements):
    """计算n元素FitcoinFish总数, 除去稀物的鱼"""

    items = [
        10, # 背景
         6, # 鱼
         5, # 背部
        10, # 衣服
         9, # 头部
         9, # 面部
        20, # 眼睛
        17, # 配饰
         9, # 嘴型
    ]

    bg_fish_count = items[0] * items[1] # 背景 和 鱼 的组合数
    items = items[2:] # 去掉背景 和 鱼
    xiwu_fish_count = counter_xiwu_n_elements(n_elements) # 计算稀物已有的 n元素的鱼的个数，去重

    # 从剩下的7种元素（69个）元素中，选2种不同类的元素进行组合
    two_element_total = 0
    for i in range(len(items) - 1):
        for j in range(i + 1, len(items)):
            two_element_total += items[i] * items[j]

    # 计算总数
    total = bg_fish_count * two_element_total  - xiwu_fish_count
    print(int(total))

if __name__ == '__main__':
    # 计算BitcoinFish 4元素的总个数
    calc_bitcoinfish_total(4)

