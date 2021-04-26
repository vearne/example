# -*- coding: UTF-8 -*-
import copy


SEPARATOR = '-'


# item -> trans
data = {}
# 支持度
# min_support = 0.2
min_support = 0.002
# 置信度
min_confidence = 0.8

size = 0
# 结果集
freq_itemsets = []


def default_judge(prefix_list):
    global freq_itemsets
    if len(prefix_list) == 0:
        return True

    key = SEPARATOR.join(prefix_list)
    if len(data[key]) >= size * min_support:
        freq_itemsets.append(key)
        return True
    else:
        # print key, False
        return False


def combination(prefix_list, item_list, func=default_judge):
    if not func(prefix_list):
        return

    for i in range(len(item_list)):
        temp = copy.copy(item_list)

        # print "length", len(item_list), i
        ch = temp.pop(i)
        if len(prefix_list) == 0 or (len(prefix_list) > 0 and prefix_list[-1] < ch):
            if len(prefix_list) > 0:
                key1 = SEPARATOR.join(prefix_list)
                key2 = SEPARATOR.join(prefix_list + [ch])
                data[key2] = data[key1] & data[ch]
            combination(prefix_list + [ch], temp)

def main(file_path):
    global size
    # 事务的数量
    with open(file_path, 'r') as fp:
        for line in fp:
            size += 1
            line = line.strip()
            item_list = line.split(":")

            items = item_list[1].split(",")
            for item in items:
                if item not in data:
                    data[item] = set()
                data[item].add(int(item_list[0]))

    print "---1. 寻找频繁项集------"
    combination([], data.keys())
    print "len(data)", len(data)
    # 结果数据集
    print "频繁项集", freq_itemsets
    print "---2. 寻找关联规则------"
    extract_rule(freq_itemsets)



def extract_rule(freq_itemsets):
    length = len(freq_itemsets)
    x = 0
    y = 0
    while x < length:
        y = x + 1
        while y < length:
            jude_rule(freq_itemsets[x], freq_itemsets[y])
            y += 1
        x += 1

def jude_rule(x, xy):
    # print '-' * 20
    # 判断 是否  x -> y 和判断    x -> xy 可以认为是等价的
    global size
    if len(x) > len(xy):
        x, xy = xy, x

    # A -> AB
    # B -> AB
    # 只有itemset2 包含 itemset1 才有可能有关联规则
    set_x = set(x.split(SEPARATOR))
    # print "set_x", set_x
    set_xy = set(xy.split(SEPARATOR))
    # print "set_xy", set_xy
    if set_x.issubset(set_xy):

        # 计算置信度
        # print 'len(data[itemset1])', len(data[x])
        # print 'len(data[itemset2])', len(data[xy])

        suport_xy = len(data[xy]) / float(size)
        suport_x = len(data[x]) / float(size)
        confidence = suport_xy / suport_x

        y = list(set_xy - set_x)
        y = SEPARATOR.join(sorted(y))
        suport_y = len(data[y]) / float(size)
        # 计算提升度
        lift = suport_xy / (suport_x * suport_y)
        # print 'suport_xy', suport_xy, 'suport_x', suport_x, 'suport_y', suport_y
        # print x, "-->", y, "confidence", confidence, "lift", lift
        # 提升度必须大于1, 规则才是正相关的
        if confidence > min_confidence and lift > 1:
            print "满足提升度要求的", x, "-->", y, "confidence", confidence, "lift", lift


if __name__ == "__main__":
    f = "./blog.dat"
    # f = "./eclat2.dat"
    main(f)

