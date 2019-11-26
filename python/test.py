import unittest

def add (a, b):
    return a + b

class TestCap(unittest.TestCase):
    def test_add(self):
        print(add(1,4))

if __name__ == '__main__':
    unittest.main()
