#!/bin/bash

shortname="test"

if [ "$1" == "clean" ]
then
    rm -rf ${shortname}.ll
    rm -rf ${shortname}.bc
    rm -rf ${shortname}.s
else
    #clang -emit-llvm -S ${shortname}.c -o ${shortname}.ll
    clang -cc1 -emit-llvm ${shortname}.c -o ${shortname}.ll
    
    llvm-as ${shortname}.ll -o ${shortname}.bc
    
    llc ${shortname}.bc -o ${shortname}.s
fi   

