#include <stdio.h>
#include <stdlib.h>

#include "aux_functions.hh"

#define N 64

int main() {
  float* in = (float*)calloc(N, sizeof(float));
  float* out = (float*)calloc(N, sizeof(float));
  float ref = 0.5f;

  for (int i = 0; i < N; ++i) {
    in[i] = scale(i, N);
  }

  distanceArray(out, in, ref, N);
  for (int i = 0; i < N; ++i) {
    printf("%.4f ", out[i]);
  }
  printf("\n");

  free(in);
  free(out);
  return 0;
}
