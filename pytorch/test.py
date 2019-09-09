import numpy as np
import torch
import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim
import sys

print(torch.cuda.is_available())
print(torch.cuda.set_device(1)) #RuntimeError: cuda runtime error (10) : invalid device ordinal
print(torch.cuda.get_device_name())
print(torch.cuda.get_device_capability())
print(torch.cuda.device_count())
print(torch.cuda.current_device())

