CC=g++ -std=c++14
#CC=clang++

TARGET=toy

COMPILE_FLAGS=$(shell llvm-config --cxxflags)
LINK_FLAGS=$(shell llvm-config --ldflags --libs)

a.out: toy.cpp
	@$(CC) -g -c $(COMPILE_FLAGS) $(TARGET).cpp -o $(TARGET).o
	@$(CC) $(TARGET).o $(LINK_FLAGS) -lpthread
	# @$(CC) -g -c $(TARGET).cpp -o $(TARGET).o
	# @$(CC) $(TARGET).o -lpthread

.PHONY:clean
clean:
	@rm -rf *.out
	@rm -rf *.o

