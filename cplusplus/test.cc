#include <stdio.h>
#include <map>

typedef struct A {
  int var = 11;
} a;


int main() {
    A a0;
    a a1;
    a1.var = 22;
    printf("%d\n", a0.var);
    printf("%d\n", a1.var);

    static std::map<char, int> BinopPrecedence;
    int& a = BinopPrecedence['c'];
    printf("a: %d\n", a);
    a = 1;
    for (auto p : BinopPrecedence) {
	printf("%c:%d\n", p.first, p.second);
    }

    return 0;
}
