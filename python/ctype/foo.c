// filename: foo.c

// gcc -fPIC -shared foo.c -o foo.so

#include "stdio.h"

char* myprint(char *str)
{
    puts(str);
    return str;
}

float add(float a, float b)
{
    return a + b;
}

typedef union {
  long long int v_int64;
  double v_float64;
  void* v_handle;
  const char* v_str;
} TVMValue;

int TVMFuncCall(TVMValue* args, int num_args) {
    printf("%u\n", args[0].v_int64);
    printf("%s\n", args[1].v_str);
}
