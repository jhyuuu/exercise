import sys
from vcd import VCDWriter

file_name = "wave.vcd"
fout = open(file_name, "w+")

with VCDWriter(fout, timescale='1 ns', date='today') as writer:
    counter_var0 = writer.register_var('a.b.c0', 'counter', 'integer', size=8)
    counter_var1 = writer.register_var('a.b.c1', 'counter', 'integer', size=8)
    for timestamp, value in enumerate(range(10, 20, 2)):
        writer.change(counter_var0, timestamp, value)
        writer.change(counter_var1, timestamp, value)
