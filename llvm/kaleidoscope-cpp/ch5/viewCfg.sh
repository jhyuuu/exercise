llvm-as < t.ll | opt -analyze -view-cfg

#sudo apt-get install graphviz
dot -Tpng -o cfg.png /tmp/cfgbaz-88bb5f.dot
eog cfg.png
