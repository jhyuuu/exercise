	.text
	.file	"test.c"
	.globl	func                    # -- Begin function func
	.p2align	4, 0x90
	.type	func,@function
func:                                   # @func
# %bb.0:                                # %entry
	movl	%edi, -4(%rsp)
	movl	-4(%rsp), %eax
	shll	$1, %eax
	movl	%eax, -4(%rsp)
	movl	-4(%rsp), %eax
	retq
.Lfunc_end0:
	.size	func, .Lfunc_end0-func
                                        # -- End function
	.ident	"clang version 10.0.0 "
	.section	".note.GNU-stack","",@progbits
