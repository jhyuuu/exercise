from __future__ import absolute_import, print_function

import os
import tvm
import vta
import numpy as np

tvm.my_api_test("aa", 11)

@tvm.register_func
def my_packed_func(*args):
    return sum(args)

f = tvm.get_global_func("my_packed_func")
print(f(1, 2, 3, 4))
