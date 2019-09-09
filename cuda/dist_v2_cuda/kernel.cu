#include "kernel.h"

#define TPB 32

__device__ float distance(float x1, float x2) {
  return sqrt((x2-x1)*(x2-x1));
}

__global__ void distanceKernel(float* d_out, float* d_in, float ref) {
  const int i = blockIdx.x * blockDim.x + threadIdx.x;
  const float x = d_in[i];
  d_out[i] = distance(x, ref);
  printf("i = %2d: dist from %f to %f is %f.\n", i, ref, x, d_out[i]);
}

void distanceArray(float* out, float* in, float ref, int n) {
  float* d_in = NULL;
  float* d_out = NULL;

  cudaMalloc(&d_in, n * sizeof(float));
  cudaMalloc(&d_out, n * sizeof(float));

  cudaMemcpy(d_in, in, n * sizeof(float), cudaMemcpyHostToDevice);

  distanceKernel<<<n/TPB, TPB>>>(d_out, d_in, ref);

  cudaMemcpy(out, d_out, n * sizeof(float), cudaMemcpyDeviceToHost);

  cudaFree(d_in);
  cudaFree(d_out);
}
