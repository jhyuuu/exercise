#include <stdio.h>

union A {
  struct {
    int a : 2;
    int b : 3;
  };
  int v : 8;
};

int main () {
    A a;
    a.a = 1;
    a.b = 2;
    printf("v: %x\n", a.v);
    printf("a: %x\n", a.a);
    printf("b: %x\n", a.b);
    printf("%lx\n", sizeof(a));
    return 0;
}
