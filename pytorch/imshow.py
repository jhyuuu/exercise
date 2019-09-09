import os
import matplotlib.pyplot as plt

img = plt.imread('20171229212728386.png')
img = img[:32,:32,0:3]
plt.figure("Image") # 图像窗口名称
plt.imshow(img)
plt.axis('on') # 关掉坐标轴为 off
plt.title('image') # 图像题目
plt.show()
