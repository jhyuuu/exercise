#ifndef __KERNEK_H_
#define __KERNEK_H_

struct uchar4;
struct int2;

void kernelLauncher(uchar4* d_out, int w, int h, int2 pos);

#endif // __KERNEK_H_
