#include <stdio.h>
#include <iostream>

namespace {
  int anonymous = 100;
  //无名命名空间允许无限定的使用其成员函数，并且为它提供了内部连接（只有在定义的文件内可以使用）
  //类似static
}

namespace First 
{
	int i = 0;
}
 
namespace Second
{
	int i = 1;
 
	namespace Internal //嵌套命名空间
	{
		struct P  //嵌套命名空间
		{
			int x;
			int y;
		};
	}
}
 
int main()
{
	using namespace First; //使用整个命名空间
	using Second::Internal::P;  //使用嵌套的命名空间
 
	printf("First::i = %d\n", i);
	printf("Second::i = %d\n", Second::i);  //使用命名空间中的变量
 
	P p = { 2, 3 };
 
	printf("p.x = %d\n", p.x);
	printf("p.y = %d\n", p.y);
 
        printf("anonymous = %d\n", anonymous);
       
	return 0;
}
