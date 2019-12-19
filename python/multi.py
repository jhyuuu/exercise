from pprint import pprint
#class base1():
#    def __init__(self,value):
#        print('base1')
#        self.value = value
#        
#class subbase1(base1):
#    def __init__(self,value1, value2):
#        print('subbase1', "value1=", value1, "value2=",  value2)
#        value2 = 2
#        super().__init__(value1,value2)
#        
#class subbase2(base1):
#    def __init__(self,value1,value2):
#        print('subbase2', "value1=", value1, "value2=", value2)
#        super().__init__(value2)
#        
#class subbase1base2(subbase1,subbase2):
#    def __init__(self,value1,value2):
#        print('subbase1base2')
#        super().__init__(value1, value2)
#        super(subbase1, self).__init__(value1)
#        super(subbase2, self).__init__(value1, value2)

class base1():
    def __init__(self,value):
        print('base1')
        self.value = value
        
class subbase1(base1):
    def __init__(self,value1):
        print('subbase1', "value1=", value1)
        value2 = 2
        super().__init__(value1,value2)
        
class subbase2(base1):
    def __init__(self,value1,value2):
        print('subbase2', "value1=", value1, "value2=", value2)
        super().__init__(value2)
        
class subbase1base2(subbase1,subbase2):
    def __init__(self,value1,value2):
        print('subbase1base2')
        super(subbase1base2, self).__init__(value1)


if __name__ == '__main__':

    s = subbase1base2(111, 222)
    print(s.value)
    
# 在python中，每个类都有一个mro的类方法,通过mro，python巧妙地将多继承的图结构，转变为list的顺序结构。
# super在继承体系中向上的查找过程，变成了在mro中向右的线性查找过程，任何类都只会被处理一次
