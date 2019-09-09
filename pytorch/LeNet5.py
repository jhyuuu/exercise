import torch
import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim


class Net(nn.Module):
  
  def __init__(self):
      super(Net, self).__init__()
      self.conv1 = nn.Conv2d(1, 6, 5)
      self.conv2 = nn.Conv2d(6, 16, 5)

      self.fc1 = nn.Linear(16*5*5, 120)
      self.fc2 = nn.Linear(120, 84)
      self.fc3 = nn.Linear(84, 10)

  def forward(self, x):
      x = F.max_pool2d(F.relu(self.conv1(x)), (2, 2))
      x = F.max_pool2d(F.relu(self.conv2(x)), (2, 2))

      x = x.view(-1, self.num_flat_features(x))
      x = F.relu(self.fc1(x))
      x = F.relu(self.fc2(x))
      x= self.fc3(x)
      return x

  def num_flat_features(self, x):
      size = x.size()[1:]
      num_features = 1
      for s in size:
          num_features *= s
      return num_features

net = Net()
print(net)

learning_rate = 0.01

input = torch.randn(1, 1, 32, 32)
target = torch.randn(10)
target = target.view(1, -1)
print('target', target)
criterion = nn.MSELoss()
optimizer = optim.Adam(net.parameters(), lr = learning_rate)

cycles = 1001
segments = 100
for cycle in range(cycles):
    net.zero_grad()
    output = net(input)
    loss = criterion(output, target)
    loss.backward()
    
    optimizer.step()

    if cycle % (segments - 1) == 0:
        print(cycle, loss.item(), output)

    if cycle == cycles - 1:
        print('output', output)
    #for f in net.parameters():
    #    f.data.sub_(learning_rate * f.grad.data)
