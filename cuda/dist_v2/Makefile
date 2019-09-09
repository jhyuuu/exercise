##### manipulate #########
ELF = run.elf

CFILES = \

CPLUSPLUSFILES = main.cpp \
                 aux_functions.cpp \

HFILE = aux_functions.hh \
##### end manipulate #####

NVCC = /usr/local/cuda/bin/nvcc
NVCC_FLAGS = -g -G -Xcompiler -Wall

OBJS = $(patsubst %.c, %.o, $(CFILES))
OBJS_CPP = $(patsubst %.cpp, %.o, $(CPLUSPLUSFILES))

DEPS = $(HFILES) Makefile

$(ELF): $(OBJS) $(OBJS_CPP)
	$(NVCC) $^ -o $@

$(OBJS): %.o:./%.c $(DEPS)
	$(NVCC) $(NVCC_FLAGS) -c $< -o $@

$(OBJS_CPP): %.o:./%.cpp $(DEPS)
	$(NVCC) $(NVCC_FLAGS) -c $< -o $@


.PHONY: clean
clean:
	@rm -rf $(ELF)
	@rm -rf *.o

