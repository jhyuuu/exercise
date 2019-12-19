#include <stdio.h>

#define TVM_REGISTER_API(OpName) TVM_REGISTER_GLOBAL(OpName)
#define TVM_REGISTER_GLOBAL(OpName) printf("ppp "#OpName"")

#define TVM_STR_CONCAT_(__x, __y) __x##__y
#define TVM_STR_CONCAT(__x, __y) TVM_STR_CONCAT_(__x, __y)

int main() {
    TVM_REGISTER_API(aaaa);
    // equal to 
    // printf("ppp aaaa");

//    printf(TVM_STR_CONCAT(mmm, nnn));

    int i = 10;
    float* j = (float*)&i;
    printf("%f\n", *j);

    return 0;
}
