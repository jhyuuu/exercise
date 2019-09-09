
##### manipulate #########
ELF = run.elf

CFILES = main.c \

CPLUSPLUSFILES = \

CUFILES =  \

HFILE = \

##### end manipulate #####
CC = g++ 

CFLAGS := -std=c++11
LDFLAGS := -lpthread

DEBUG = 1
export DEBUG
ifeq ($(DEBUG), 1)
	CFLAGS += -O0 -g -DDEBUG
endif

OPENCV = 1
export OPENCV
ifeq ($(OPENCV), 1)
COMMON+= -DOPENCV
CFLAGS+= -DOPENCV
LDFLAGS+= `pkg-config --libs opencv`
COMMON+= `pkg-config --cflags opencv`
endif

NVCC = /usr/local/cuda/bin/nvcc
NVCC_FLAGS = -g -G -Xcompiler -Wall

OBJS = $(patsubst %.c, %.o, $(CFILES))
OBJS_CPP = $(patsubst %.cpp, %.o, $(CPLUSPLUSFILES))
OBJS_CU = $(patsubst %.cu, %.o, $(CUFILES))

DEPS = $(HFILES) Makefile

$(ELF): $(OBJS) $(OBJS_CPP) $(OBJS_CU)
	$(CC) $^ -o $@ $(LDFLAGS)

$(OBJS): %.o:./%.c $(DEPS)
	$(CC) $(COMMON) $(CFLAGS) -c $< -o $@

$(OBJS_CPP): %.o:./%.cpp $(DEPS)
	$(CC) $(CFLAGS) -c $< -o $@

$(OBJS_CU): %.o:./%.cu $(DEPS)
	$(CC) $(CFLAGS) -c $< -o $@

.PHONY: clean
clean:
	@rm -rf $(ELF)
	@rm -rf *.o
