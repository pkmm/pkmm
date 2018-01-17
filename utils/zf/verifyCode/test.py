# author:  ZhangChuanCheng

import os
import urllib
import urllib2
from PIL import Image
from cStringIO import StringIO
from os import path
import sys

from libsvm.python.svmutil import *

base_folder = path.dirname(path.abspath(__file__))


def scale(x):
    return x / 255.0


def load_data():
    kv_dict = {}
    with open(base_folder + os.sep + "datasets/answer.txt") as f:
        for pre, answer in enumerate(f):
            answer = answer.strip()
            answer = map(lambda x: x - 48 if x <= 57 else x -
                                                          87 if x <= 110 else x - 88, map(ord, answer))
            for i, v in enumerate(answer):
                kv_dict['%s-%d.png' % (pre, i)] = v
    floder = base_folder + os.sep + "datasets/sample_single"
    imgs = kv_dict.keys()
    data = []
    labels = []
    for index, img_name in enumerate(imgs):
        img = Image.open("%s/%s" % (floder, img_name))
        labels.append(kv_dict[img_name])
        data.append(map(scale, list(img.getdata())))
    with open('tmp.txt', 'w') as f:
        for index, x in enumerate(data):
            f.write(str(labels[index]) + ' ')
            for i, xx in enumerate(x):
                f.write(str(i + 1) + ':' + str(x[i]) + ' ')
            f.write('\n')
    return labels, data


def train():
    y, x = load_data()
    prob = svm_problem(y, x)
    param = svm_parameter('-g 0.0001220703125 -c 512 -b 0')
    m = svm_train(prob, param)
    svm_save_model(base_folder + os.sep + "t.model", m)


def verify(url, save=False):
    picBaseUrl = base_folder + os.sep + 'cache' + os.sep
    if save:
        pic_file = picBaseUrl + 'todo.png'
        urllib.urlretrieve(url, pic_file)
        image = Image.open(pic_file).convert("L")
    else:
        image = Image.open(StringIO(urllib2.urlopen(url).read()))
    x_size, y_size = image.size
    y_size -= 5

    # y from 1 to y_size-5
    # x from 4 to x_size-18
    piece = (x_size - 22) / 8
    centers = [4 + piece * (2 * i + 1) for i in range(4)]
    data = []
    answers = []
    m = svm_load_model(base_folder + os.sep + 't.model')
    for i, center in enumerate(centers):
        single_pic = image.crop(
            (center - (piece + 2), 1, center + (piece + 2), y_size))
        data.append(map(scale, list(single_pic.getdata())))
        if save:
            single_pic.save(picBaseUrl + 'todo-%s.png' % i)
    p_label, p_acc, p_val = svm_predict([0, 0, 0, 0], data, m, "-q")

    answers = p_label

    answers = map(chr, map(lambda x: x + 48 if x <= 9 else x +
                                                           87 if x <= 23 else x + 88, map(int, answers)))
    return ''.join(answers)

def verify_local_img(img_name):
    image = Image.open(img_name).convert("L")
    x_size, y_size = image.size
    y_size -= 5
    # y from 1 to y_size-5
    # x from 4 to x_size-18
    piece = (x_size - 22) / 8
    centers = [4 + piece * (2 * i + 1) for i in range(4)]
    data = []
    answers = []
    m = svm_load_model(base_folder + os.sep + 't.model')
    for i, center in enumerate(centers):
        single_pic = image.crop(
            (center - (piece + 2), 1, center + (piece + 2), y_size))
        data.append(map(scale, list(single_pic.getdata())))
    p_label, p_acc, p_val = svm_predict([1, 2, 6, 3], data, m)

    answers = p_label

    answers = map(chr, map(lambda x: x + 48 if x <= 9 else x + 87 if x <= 23 else x + 88, map(int, answers)))
    os.remove(img_name)
    return ''.join(answers)

if __name__ == "__main__":
    # train()
    # Url = 'http://zfxk.zjtcm.net/(ys0fjh555nqen445o53dza45)/CheckCode.aspx'
    # print verify(Url, True)
    picName = sys.argv[1]
    # print picName
    # print base_folder + os.sep + picName
    print verify_local_img(base_folder + os.sep + picName)