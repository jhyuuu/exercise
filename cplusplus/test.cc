#include <stdio.h>
#include <map>


int main() {
    char A[2] = {'a', 'b'};
    printf("%d\n", sizeof(char));
    printf("%p %p\n", A, A+1);
    return 0;
}
