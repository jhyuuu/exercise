import torch
import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim
import sys


conv = nn.Conv2d(3, 16, 4, 2, 1, bias = True)
assert conv.in_channels == 3
assert conv.out_channels == 16
assert conv.kernel_size == (4, 4)
assert conv.stride == (2, 2)
assert conv.padding == (1, 1)
assert conv.groups == 1
assert conv.bias.shape == (16,)
assert conv.weight.shape == (16, 3, 4, 4) #out_channels, in_channels/groups, *kernel_size

x = torch.randn(1, 3, 416, 416)
x = x.cuda()

x = conv(x)

conv1 = nn.Conv2d(3, 16, 4, 2, 1, bias = False)
assert conv1.bias == None



