#include "kernel.h"

#define TX 32
#define TY 32

__device__ unsigned char clip(int n) {
  return n > 255 ? 255 : ((n < 0) ? 0 : n);
}

__global__ void distanceKernel(uchar4* d_out, int w, int h, int2 pos) {
  const int c = blockIdx.x * blockDim.x + threadIdx.x;
  const int r = blockIdx.y * blockDim.y + threadIdx.y;

  if (c >= w || r >= h) return;

  const int i = r * w + c;
  const int dist = sqrtf((c-pos.x)*(c-pos.x) + (r-pos.y)*(r-pos.y));

  unsigned char intensity = clip(255 - dist);

  d_out[i].x = intensity;
  d_out[i].y = intensity;
  d_out[i].z = 0;
  d_out[i].w = 255;
}

void kernelLauncher(uchar4* d_out, int w, int h, int2 pos) {
  const dim3 blockSize(TX, TY);
  const dim3 gridSize((w + TX - 1)/TX, (h + TY - 1)/TY);
  distanceKernel<<<gridSize, blockSize>>>(d_out, w, h, pos);
}
