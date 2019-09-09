#include "kernel.h"

#define TPB 32

__device__ float distance(float x1, float x2) {
  return sqrt((x2-x1)*(x2-x1));
}

__global__ void distanceKernel(float* d_out, float* d_in, float ref) {
  const int i = blockIdx.x * blockDim.x + threadIdx.x;
  const float x = d_in[i];
  d_out[i] = distance(x, ref);
  //printf("i = %2d: dist from %f to %f is %f.\n", i, ref, x, d_out[i]);
}

void distanceArray(float* out, float* in, float ref, int n) {
  cudaEvent_t startMemcpy, stopMemcpy;
  cudaEvent_t startKernel, stopKernel;

  cudaEventCreate(&startMemcpy);
  cudaEventCreate(&stopMemcpy);
  cudaEventCreate(&startKernel);
  cudaEventCreate(&stopKernel);

  float* d_in = NULL;
  float* d_out = NULL;

  cudaMalloc(&d_in, n * sizeof(float));
  cudaMalloc(&d_out, n * sizeof(float));

  cudaEventRecord(startMemcpy);
  for (int i = 0; i < 200; ++i) {
    cudaMemcpy(d_in, in, n * sizeof(float), cudaMemcpyHostToDevice);
  }
  cudaEventRecord(stopMemcpy);

  cudaEventRecord(startKernel);
  distanceKernel<<<n/TPB, TPB>>>(d_out, d_in, ref);
  cudaEventRecord(stopKernel);

  cudaMemcpy(out, d_out, n * sizeof(float), cudaMemcpyDeviceToHost);

  cudaEventSynchronize(stopMemcpy);
  cudaEventSynchronize(stopKernel);

  float memcpyTimeInMs = 0;
  float kernelTimeInMs = 0;
  cudaEventElapsedTime(&memcpyTimeInMs, startMemcpy, stopMemcpy);
  cudaEventElapsedTime(&kernelTimeInMs, startKernel, stopKernel);
  printf("Data transfer time (ms): %f\n", memcpyTimeInMs);
  printf("Kernel time (ms): %f\n", kernelTimeInMs);

  cudaFree(d_in);
  cudaFree(d_out);
}
