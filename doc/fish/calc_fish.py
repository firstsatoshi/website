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

def calc_bitfish_total():
    """计算bitfish总数"""
    items = [10, 5, 6 , 10 , 9 , 9 , 20 , 17 , 9]
    m = sum(items)
    total = 0
    for n in range(1, 8):
        total += C(m, n)
    print(int(total))

if __name__ == '__main__':
    calc_bitfish_total()

