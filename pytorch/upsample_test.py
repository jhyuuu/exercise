import torch
import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim
import sys

input = torch.arange(1, 5, dtype=torch.float32).view(1, 1, 2, 2)
m = nn.Upsample(scale_factor=3, mode='nearest')
print(input)
result = m(input)
print(result)
