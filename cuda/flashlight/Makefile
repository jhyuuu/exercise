##### manipulate #########
ELF = run.elf

CFILES = \

CPLUSPLUSFILES = main.cpp \

CUFILES = kernel.cu \

HFILE = kernel.h \
        interactions.h \

##### end manipulate #####

LDFLAGS += -L/usr/local/cuda/samples/common/lib/linux/x86_64
LDFLAGS += -lglut -lGL -lGLU -lGLEW

NVCC = /usr/local/cuda/bin/nvcc
NVCC_FLAGS = -g -G -Xcompiler "-Wall -Wno-deprecated-declarations"
INC = -I/usr/local/cuda/samples/common/inc

OBJS = $(patsubst %.c, %.o, $(CFILES))
OBJS_CPP = $(patsubst %.cpp, %.o, $(CPLUSPLUSFILES))
OBJS_CU = $(patsubst %.cu, %.o, $(CUFILES))

DEPS = $(HFILES) Makefile

$(ELF): $(OBJS) $(OBJS_CPP) $(OBJS_CU) $(HFILES)
	$(NVCC) $(LDFLAGS) $^ -o $@ 

$(OBJS): %.o:./%.c $(DEPS)
	$(NVCC) $(NVCC_FLAGS) -c $< -o $@

$(OBJS_CPP): %.o:./%.cpp $(DEPS)
	$(NVCC) $(NVCC_FLAGS) -c $< -o $@

$(OBJS_CU): %.o:./%.cu $(DEPS)
	$(NVCC) $(NVCC_FLAGS) -c $< -o $@

.PHONY: clean
clean:
	@rm -rf $(ELF)
	@rm -rf *.o


