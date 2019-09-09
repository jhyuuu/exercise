import torch
import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim
import sys

activn = nn.LeakyReLU(0.1, inplace = True)
x = torch.randn(1, 3, 416, 416)
#print(x[0, 0, :10, :10])
#activn(x)
#print(x[0, 0, :10, :10])
#
#activn = nn.LeakyReLU(0.1)
#x = torch.randn(1, 3, 416, 416)
#print(x[0, 0, :10, :10])
#y = activn(x)
#print(x[0, 0, :10, :10])
#print(y[0, 0, :10, :10])

x = x.cuda()
assert x.is_cuda
print(x[0, 0, :10, :10])
activn(x)
print(x[0, 0, :10, :10])

