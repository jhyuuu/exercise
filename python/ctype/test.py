import ctypes

class TVMValue(ctypes.Union):
    """TVMValue in C API"""
    _fields_ = [("v_int64", ctypes.c_int64),
                ("v_float64", ctypes.c_double),
                ("v_handle", ctypes.c_void_p),
                ("v_str", ctypes.c_char_p)]
num_args = 2
b = (TVMValue * num_args)()
b[0].v_int64 = 100
b[1].v_str = ctypes.c_char_p("aaa".encode('utf-8'))

foo = ctypes.CDLL('./foo.so')
foo.TVMFuncCall(b, num_args)

