.text	

.p2align	6

K256:
.long	0x428a2f98,0x71374491,0xb5c0fbcf,0xe9b5dba5
.long	0x3956c25b,0x59f111f1,0x923f82a4,0xab1c5ed5
.long	0xd807aa98,0x12835b01,0x243185be,0x550c7dc3
.long	0x72be5d74,0x80deb1fe,0x9bdc06a7,0xc19bf174
.long	0xe49b69c1,0xefbe4786,0x0fc19dc6,0x240ca1cc
.long	0x2de92c6f,0x4a7484aa,0x5cb0a9dc,0x76f988da
.long	0x983e5152,0xa831c66d,0xb00327c8,0xbf597fc7
.long	0xc6e00bf3,0xd5a79147,0x06ca6351,0x14292967
.long	0x27b70a85,0x2e1b2138,0x4d2c6dfc,0x53380d13
.long	0x650a7354,0x766a0abb,0x81c2c92e,0x92722c85
.long	0xa2bfe8a1,0xa81a664b,0xc24b8b70,0xc76c51a3
.long	0xd192e819,0xd6990624,0xf40e3585,0x106aa070
.long	0x19a4c116,0x1e376c08,0x2748774c,0x34b0bcb5
.long	0x391c0cb3,0x4ed8aa4a,0x5b9cca4f,0x682e6ff3
.long	0x748f82ee,0x78a5636f,0x84c87814,0x8cc70208
.long	0x90befffa,0xa4506ceb,0xbef9a3f7,0xc67178f2

.long	0x00010203,0x04050607,0x08090a0b,0x0c0d0e0f
.long	0x03020100,0x0b0a0908,0xffffffff,0xffffffff
.long	0xffffffff,0xffffffff,0x03020100,0x0b0a0908
.byte	83,72,65,50,53,54,32,98,108,111,99,107,32,116,114,97,110,115,102,111,114,109,32,102,111,114,32,120,56,54,95,54,52,44,32,67,82,89,80,84,79,71,65,77,83,32,98,121,32,64,100,111,116,45,97,115,109,0
.globl	blst_sha256_block_data_order_shaext

.def	blst_sha256_block_data_order_shaext;	.scl 2;	.type 32;	.endef
.p2align	6
blst_sha256_block_data_order_shaext:
	.byte	0xf3,0x0f,0x1e,0xfa
	movq	%rdi,8(%rsp)
	movq	%rsi,16(%rsp)
	movq	%rsp,%r11
.LSEH_begin_blst_sha256_block_data_order_shaext:
	movq	%rcx,%rdi
	movq	%rdx,%rsi
	movq	%r8,%rdx


	subq	$0x58,%rsp

	movaps	%xmm6,-88(%r11)

	movaps	%xmm7,-72(%r11)

	movaps	%xmm8,-56(%r11)

	movaps	%xmm9,-40(%r11)

	movaps	%xmm10,-24(%r11)

.LSEH_body_blst_sha256_block_data_order_shaext:

	leaq	K256+128(%rip),%rcx
	movdqu	(%rdi),%xmm1
	movdqu	16(%rdi),%xmm2
	movdqa	256-128(%rcx),%xmm7

	pshufd	$0x1b,%xmm1,%xmm0
	pshufd	$0xb1,%xmm1,%xmm1
	pshufd	$0x1b,%xmm2,%xmm2
	movdqa	%xmm7,%xmm8
.byte	102,15,58,15,202,8
	punpcklqdq	%xmm0,%xmm2
	jmp	.Loop_shaext

.p2align	4
.Loop_shaext:
	movdqu	(%rsi),%xmm3
	movdqu	16(%rsi),%xmm4
	movdqu	32(%rsi),%xmm5
.byte	102,15,56,0,223
	movdqu	48(%rsi),%xmm6

	movdqa	0-128(%rcx),%xmm0
	paddd	%xmm3,%xmm0
.byte	102,15,56,0,231
	movdqa	%xmm2,%xmm10
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
	nop
	movdqa	%xmm1,%xmm9
.byte	15,56,203,202

	movdqa	16-128(%rcx),%xmm0
	paddd	%xmm4,%xmm0
.byte	102,15,56,0,239
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
	leaq	64(%rsi),%rsi
.byte	15,56,204,220
.byte	15,56,203,202

	movdqa	32-128(%rcx),%xmm0
	paddd	%xmm5,%xmm0
.byte	102,15,56,0,247
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
	movdqa	%xmm6,%xmm7
.byte	102,15,58,15,253,4
	nop
	paddd	%xmm7,%xmm3
.byte	15,56,204,229
.byte	15,56,203,202

	movdqa	48-128(%rcx),%xmm0
	paddd	%xmm6,%xmm0
.byte	15,56,205,222
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
	movdqa	%xmm3,%xmm7
.byte	102,15,58,15,254,4
	nop
	paddd	%xmm7,%xmm4
.byte	15,56,204,238
.byte	15,56,203,202
	movdqa	64-128(%rcx),%xmm0
	paddd	%xmm3,%xmm0
.byte	15,56,205,227
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
	movdqa	%xmm4,%xmm7
.byte	102,15,58,15,251,4
	nop
	paddd	%xmm7,%xmm5
.byte	15,56,204,243
.byte	15,56,203,202
	movdqa	80-128(%rcx),%xmm0
	paddd	%xmm4,%xmm0
.byte	15,56,205,236
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
	movdqa	%xmm5,%xmm7
.byte	102,15,58,15,252,4
	nop
	paddd	%xmm7,%xmm6
.byte	15,56,204,220
.byte	15,56,203,202
	movdqa	96-128(%rcx),%xmm0
	paddd	%xmm5,%xmm0
.byte	15,56,205,245
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
	movdqa	%xmm6,%xmm7
.byte	102,15,58,15,253,4
	nop
	paddd	%xmm7,%xmm3
.byte	15,56,204,229
.byte	15,56,203,202
	movdqa	112-128(%rcx),%xmm0
	paddd	%xmm6,%xmm0
.byte	15,56,205,222
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
	movdqa	%xmm3,%xmm7
.byte	102,15,58,15,254,4
	nop
	paddd	%xmm7,%xmm4
.byte	15,56,204,238
.byte	15,56,203,202
	movdqa	128-128(%rcx),%xmm0
	paddd	%xmm3,%xmm0
.byte	15,56,205,227
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
	movdqa	%xmm4,%xmm7
.byte	102,15,58,15,251,4
	nop
	paddd	%xmm7,%xmm5
.byte	15,56,204,243
.byte	15,56,203,202
	movdqa	144-128(%rcx),%xmm0
	paddd	%xmm4,%xmm0
.byte	15,56,205,236
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
	movdqa	%xmm5,%xmm7
.byte	102,15,58,15,252,4
	nop
	paddd	%xmm7,%xmm6
.byte	15,56,204,220
.byte	15,56,203,202
	movdqa	160-128(%rcx),%xmm0
	paddd	%xmm5,%xmm0
.byte	15,56,205,245
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
	movdqa	%xmm6,%xmm7
.byte	102,15,58,15,253,4
	nop
	paddd	%xmm7,%xmm3
.byte	15,56,204,229
.byte	15,56,203,202
	movdqa	176-128(%rcx),%xmm0
	paddd	%xmm6,%xmm0
.byte	15,56,205,222
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
	movdqa	%xmm3,%xmm7
.byte	102,15,58,15,254,4
	nop
	paddd	%xmm7,%xmm4
.byte	15,56,204,238
.byte	15,56,203,202
	movdqa	192-128(%rcx),%xmm0
	paddd	%xmm3,%xmm0
.byte	15,56,205,227
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
	movdqa	%xmm4,%xmm7
.byte	102,15,58,15,251,4
	nop
	paddd	%xmm7,%xmm5
.byte	15,56,204,243
.byte	15,56,203,202
	movdqa	208-128(%rcx),%xmm0
	paddd	%xmm4,%xmm0
.byte	15,56,205,236
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
	movdqa	%xmm5,%xmm7
.byte	102,15,58,15,252,4
.byte	15,56,203,202
	paddd	%xmm7,%xmm6

	movdqa	224-128(%rcx),%xmm0
	paddd	%xmm5,%xmm0
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
.byte	15,56,205,245
	movdqa	%xmm8,%xmm7
.byte	15,56,203,202

	movdqa	240-128(%rcx),%xmm0
	paddd	%xmm6,%xmm0
	nop
.byte	15,56,203,209
	pshufd	$0x0e,%xmm0,%xmm0
	decq	%rdx
	nop
.byte	15,56,203,202

	paddd	%xmm10,%xmm2
	paddd	%xmm9,%xmm1
	jnz	.Loop_shaext

	pshufd	$0xb1,%xmm2,%xmm2
	pshufd	$0x1b,%xmm1,%xmm7
	pshufd	$0xb1,%xmm1,%xmm1
	punpckhqdq	%xmm2,%xmm1
.byte	102,15,58,15,215,8

	movdqu	%xmm1,(%rdi)
	movdqu	%xmm2,16(%rdi)
	movaps	-88(%r11),%xmm6
	movaps	-72(%r11),%xmm7
	movaps	-56(%r11),%xmm8
	movaps	-40(%r11),%xmm9
	movaps	-24(%r11),%xmm10
	movq	%r11,%rsp

.LSEH_epilogue_blst_sha256_block_data_order_shaext:
	mov	8(%r11),%rdi
	mov	16(%r11),%rsi

	.byte	0xf3,0xc3

.LSEH_end_blst_sha256_block_data_order_shaext:
.globl	blst_sha256_block_data_order

.def	blst_sha256_block_data_order;	.scl 2;	.type 32;	.endef
.p2align	6
blst_sha256_block_data_order:
	.byte	0xf3,0x0f,0x1e,0xfa
	movq	%rdi,8(%rsp)
	movq	%rsi,16(%rsp)
	movq	%rsp,%r11
.LSEH_begin_blst_sha256_block_data_order:
	movq	%rcx,%rdi
	movq	%rdx,%rsi
	movq	%r8,%rdx


	pushq	%rbp

	pushq	%rbx

	pushq	%r12

	pushq	%r13

	pushq	%r14

	pushq	%r15

	shlq	$4,%rdx
	subq	$104,%rsp

	leaq	(%rsi,%rdx,4),%rdx
	movq	%rdi,0(%rsp)

	movq	%rdx,16(%rsp)
	movaps	%xmm6,32(%rsp)

	movaps	%xmm7,48(%rsp)

	movaps	%xmm8,64(%rsp)

	movaps	%xmm9,80(%rsp)

	movq	%rsp,%rbp

.LSEH_body_blst_sha256_block_data_order:


	leaq	-64(%rsp),%rsp
	movl	0(%rdi),%eax
	andq	$-64,%rsp
	movl	4(%rdi),%ebx
	movl	8(%rdi),%ecx
	movl	12(%rdi),%edx
	movl	16(%rdi),%r8d
	movl	20(%rdi),%r9d
	movl	24(%rdi),%r10d
	movl	28(%rdi),%r11d


	jmp	.Lloop_ssse3
.p2align	4
.Lloop_ssse3:
	movdqa	K256+256(%rip),%xmm7
	movq	%rsi,8(%rbp)
	movdqu	0(%rsi),%xmm0
	movdqu	16(%rsi),%xmm1
	movdqu	32(%rsi),%xmm2
.byte	102,15,56,0,199
	movdqu	48(%rsi),%xmm3
	leaq	K256(%rip),%rsi
.byte	102,15,56,0,207
	movdqa	0(%rsi),%xmm4
	movdqa	16(%rsi),%xmm5
.byte	102,15,56,0,215
	paddd	%xmm0,%xmm4
	movdqa	32(%rsi),%xmm6
.byte	102,15,56,0,223
	movdqa	48(%rsi),%xmm7
	paddd	%xmm1,%xmm5
	paddd	%xmm2,%xmm6
	paddd	%xmm3,%xmm7
	movdqa	%xmm4,0(%rsp)
	movl	%eax,%r14d
	movdqa	%xmm5,16(%rsp)
	movl	%ebx,%edi
	movdqa	%xmm6,32(%rsp)
	xorl	%ecx,%edi
	movdqa	%xmm7,48(%rsp)
	movl	%r8d,%r13d
	jmp	.Lssse3_00_47

.p2align	4
.Lssse3_00_47:
	subq	$-64,%rsi
	rorl	$14,%r13d
	movdqa	%xmm1,%xmm4
	movl	%r14d,%eax
	movl	%r9d,%r12d
	movdqa	%xmm3,%xmm7
	rorl	$9,%r14d
	xorl	%r8d,%r13d
	xorl	%r10d,%r12d
	rorl	$5,%r13d
	xorl	%eax,%r14d
.byte	102,15,58,15,224,4
	andl	%r8d,%r12d
	xorl	%r8d,%r13d
.byte	102,15,58,15,250,4
	addl	0(%rsp),%r11d
	movl	%eax,%r15d
	xorl	%r10d,%r12d
	rorl	$11,%r14d
	movdqa	%xmm4,%xmm5
	xorl	%ebx,%r15d
	addl	%r12d,%r11d
	movdqa	%xmm4,%xmm6
	rorl	$6,%r13d
	andl	%r15d,%edi
	psrld	$3,%xmm4
	xorl	%eax,%r14d
	addl	%r13d,%r11d
	xorl	%ebx,%edi
	paddd	%xmm7,%xmm0
	rorl	$2,%r14d
	addl	%r11d,%edx
	psrld	$7,%xmm6
	addl	%edi,%r11d
	movl	%edx,%r13d
	pshufd	$250,%xmm3,%xmm7
	addl	%r11d,%r14d
	rorl	$14,%r13d
	pslld	$14,%xmm5
	movl	%r14d,%r11d
	movl	%r8d,%r12d
	pxor	%xmm6,%xmm4
	rorl	$9,%r14d
	xorl	%edx,%r13d
	xorl	%r9d,%r12d
	rorl	$5,%r13d
	psrld	$11,%xmm6
	xorl	%r11d,%r14d
	pxor	%xmm5,%xmm4
	andl	%edx,%r12d
	xorl	%edx,%r13d
	pslld	$11,%xmm5
	addl	4(%rsp),%r10d
	movl	%r11d,%edi
	pxor	%xmm6,%xmm4
	xorl	%r9d,%r12d
	rorl	$11,%r14d
	movdqa	%xmm7,%xmm6
	xorl	%eax,%edi
	addl	%r12d,%r10d
	pxor	%xmm5,%xmm4
	rorl	$6,%r13d
	andl	%edi,%r15d
	xorl	%r11d,%r14d
	psrld	$10,%xmm7
	addl	%r13d,%r10d
	xorl	%eax,%r15d
	paddd	%xmm4,%xmm0
	rorl	$2,%r14d
	addl	%r10d,%ecx
	psrlq	$17,%xmm6
	addl	%r15d,%r10d
	movl	%ecx,%r13d
	addl	%r10d,%r14d
	pxor	%xmm6,%xmm7
	rorl	$14,%r13d
	movl	%r14d,%r10d
	movl	%edx,%r12d
	rorl	$9,%r14d
	psrlq	$2,%xmm6
	xorl	%ecx,%r13d
	xorl	%r8d,%r12d
	pxor	%xmm6,%xmm7
	rorl	$5,%r13d
	xorl	%r10d,%r14d
	andl	%ecx,%r12d
	pshufd	$128,%xmm7,%xmm7
	xorl	%ecx,%r13d
	addl	8(%rsp),%r9d
	movl	%r10d,%r15d
	psrldq	$8,%xmm7
	xorl	%r8d,%r12d
	rorl	$11,%r14d
	xorl	%r11d,%r15d
	addl	%r12d,%r9d
	rorl	$6,%r13d
	paddd	%xmm7,%xmm0
	andl	%r15d,%edi
	xorl	%r10d,%r14d
	addl	%r13d,%r9d
	pshufd	$80,%xmm0,%xmm7
	xorl	%r11d,%edi
	rorl	$2,%r14d
	addl	%r9d,%ebx
	movdqa	%xmm7,%xmm6
	addl	%edi,%r9d
	movl	%ebx,%r13d
	psrld	$10,%xmm7
	addl	%r9d,%r14d
	rorl	$14,%r13d
	psrlq	$17,%xmm6
	movl	%r14d,%r9d
	movl	%ecx,%r12d
	pxor	%xmm6,%xmm7
	rorl	$9,%r14d
	xorl	%ebx,%r13d
	xorl	%edx,%r12d
	rorl	$5,%r13d
	xorl	%r9d,%r14d
	psrlq	$2,%xmm6
	andl	%ebx,%r12d
	xorl	%ebx,%r13d
	addl	12(%rsp),%r8d
	pxor	%xmm6,%xmm7
	movl	%r9d,%edi
	xorl	%edx,%r12d
	rorl	$11,%r14d
	pshufd	$8,%xmm7,%xmm7
	xorl	%r10d,%edi
	addl	%r12d,%r8d
	movdqa	0(%rsi),%xmm6
	rorl	$6,%r13d
	andl	%edi,%r15d
	pslldq	$8,%xmm7
	xorl	%r9d,%r14d
	addl	%r13d,%r8d
	xorl	%r10d,%r15d
	paddd	%xmm7,%xmm0
	rorl	$2,%r14d
	addl	%r8d,%eax
	addl	%r15d,%r8d
	paddd	%xmm0,%xmm6
	movl	%eax,%r13d
	addl	%r8d,%r14d
	movdqa	%xmm6,0(%rsp)
	rorl	$14,%r13d
	movdqa	%xmm2,%xmm4
	movl	%r14d,%r8d
	movl	%ebx,%r12d
	movdqa	%xmm0,%xmm7
	rorl	$9,%r14d
	xorl	%eax,%r13d
	xorl	%ecx,%r12d
	rorl	$5,%r13d
	xorl	%r8d,%r14d
.byte	102,15,58,15,225,4
	andl	%eax,%r12d
	xorl	%eax,%r13d
.byte	102,15,58,15,251,4
	addl	16(%rsp),%edx
	movl	%r8d,%r15d
	xorl	%ecx,%r12d
	rorl	$11,%r14d
	movdqa	%xmm4,%xmm5
	xorl	%r9d,%r15d
	addl	%r12d,%edx
	movdqa	%xmm4,%xmm6
	rorl	$6,%r13d
	andl	%r15d,%edi
	psrld	$3,%xmm4
	xorl	%r8d,%r14d
	addl	%r13d,%edx
	xorl	%r9d,%edi
	paddd	%xmm7,%xmm1
	rorl	$2,%r14d
	addl	%edx,%r11d
	psrld	$7,%xmm6
	addl	%edi,%edx
	movl	%r11d,%r13d
	pshufd	$250,%xmm0,%xmm7
	addl	%edx,%r14d
	rorl	$14,%r13d
	pslld	$14,%xmm5
	movl	%r14d,%edx
	movl	%eax,%r12d
	pxor	%xmm6,%xmm4
	rorl	$9,%r14d
	xorl	%r11d,%r13d
	xorl	%ebx,%r12d
	rorl	$5,%r13d
	psrld	$11,%xmm6
	xorl	%edx,%r14d
	pxor	%xmm5,%xmm4
	andl	%r11d,%r12d
	xorl	%r11d,%r13d
	pslld	$11,%xmm5
	addl	20(%rsp),%ecx
	movl	%edx,%edi
	pxor	%xmm6,%xmm4
	xorl	%ebx,%r12d
	rorl	$11,%r14d
	movdqa	%xmm7,%xmm6
	xorl	%r8d,%edi
	addl	%r12d,%ecx
	pxor	%xmm5,%xmm4
	rorl	$6,%r13d
	andl	%edi,%r15d
	xorl	%edx,%r14d
	psrld	$10,%xmm7
	addl	%r13d,%ecx
	xorl	%r8d,%r15d
	paddd	%xmm4,%xmm1
	rorl	$2,%r14d
	addl	%ecx,%r10d
	psrlq	$17,%xmm6
	addl	%r15d,%ecx
	movl	%r10d,%r13d
	addl	%ecx,%r14d
	pxor	%xmm6,%xmm7
	rorl	$14,%r13d
	movl	%r14d,%ecx
	movl	%r11d,%r12d
	rorl	$9,%r14d
	psrlq	$2,%xmm6
	xorl	%r10d,%r13d
	xorl	%eax,%r12d
	pxor	%xmm6,%xmm7
	rorl	$5,%r13d
	xorl	%ecx,%r14d
	andl	%r10d,%r12d
	pshufd	$128,%xmm7,%xmm7
	xorl	%r10d,%r13d
	addl	24(%rsp),%ebx
	movl	%ecx,%r15d
	psrldq	$8,%xmm7
	xorl	%eax,%r12d
	rorl	$11,%r14d
	xorl	%edx,%r15d
	addl	%r12d,%ebx
	rorl	$6,%r13d
	paddd	%xmm7,%xmm1
	andl	%r15d,%edi
	xorl	%ecx,%r14d
	addl	%r13d,%ebx
	pshufd	$80,%xmm1,%xmm7
	xorl	%edx,%edi
	rorl	$2,%r14d
	addl	%ebx,%r9d
	movdqa	%xmm7,%xmm6
	addl	%edi,%ebx
	movl	%r9d,%r13d
	psrld	$10,%xmm7
	addl	%ebx,%r14d
	rorl	$14,%r13d
	psrlq	$17,%xmm6
	movl	%r14d,%ebx
	movl	%r10d,%r12d
	pxor	%xmm6,%xmm7
	rorl	$9,%r14d
	xorl	%r9d,%r13d
	xorl	%r11d,%r12d
	rorl	$5,%r13d
	xorl	%ebx,%r14d
	psrlq	$2,%xmm6
	andl	%r9d,%r12d
	xorl	%r9d,%r13d
	addl	28(%rsp),%eax
	pxor	%xmm6,%xmm7
	movl	%ebx,%edi
	xorl	%r11d,%r12d
	rorl	$11,%r14d
	pshufd	$8,%xmm7,%xmm7
	xorl	%ecx,%edi
	addl	%r12d,%eax
	movdqa	16(%rsi),%xmm6
	rorl	$6,%r13d
	andl	%edi,%r15d
	pslldq	$8,%xmm7
	xorl	%ebx,%r14d
	addl	%r13d,%eax
	xorl	%ecx,%r15d
	paddd	%xmm7,%xmm1
	rorl	$2,%r14d
	addl	%eax,%r8d
	addl	%r15d,%eax
	paddd	%xmm1,%xmm6
	movl	%r8d,%r13d
	addl	%eax,%r14d
	movdqa	%xmm6,16(%rsp)
	rorl	$14,%r13d
	movdqa	%xmm3,%xmm4
	movl	%r14d,%eax
	movl	%r9d,%r12d
	movdqa	%xmm1,%xmm7
	rorl	$9,%r14d
	xorl	%r8d,%r13d
	xorl	%r10d,%r12d
	rorl	$5,%r13d
	xorl	%eax,%r14d
.byte	102,15,58,15,226,4
	andl	%r8d,%r12d
	xorl	%r8d,%r13d
.byte	102,15,58,15,248,4
	addl	32(%rsp),%r11d
	movl	%eax,%r15d
	xorl	%r10d,%r12d
	rorl	$11,%r14d
	movdqa	%xmm4,%xmm5
	xorl	%ebx,%r15d
	addl	%r12d,%r11d
	movdqa	%xmm4,%xmm6
	rorl	$6,%r13d
	andl	%r15d,%edi
	psrld	$3,%xmm4
	xorl	%eax,%r14d
	addl	%r13d,%r11d
	xorl	%ebx,%edi
	paddd	%xmm7,%xmm2
	rorl	$2,%r14d
	addl	%r11d,%edx
	psrld	$7,%xmm6
	addl	%edi,%r11d
	movl	%edx,%r13d
	pshufd	$250,%xmm1,%xmm7
	addl	%r11d,%r14d
	rorl	$14,%r13d
	pslld	$14,%xmm5
	movl	%r14d,%r11d
	movl	%r8d,%r12d
	pxor	%xmm6,%xmm4
	rorl	$9,%r14d
	xorl	%edx,%r13d
	xorl	%r9d,%r12d
	rorl	$5,%r13d
	psrld	$11,%xmm6
	xorl	%r11d,%r14d
	pxor	%xmm5,%xmm4
	andl	%edx,%r12d
	xorl	%edx,%r13d
	pslld	$11,%xmm5
	addl	36(%rsp),%r10d
	movl	%r11d,%edi
	pxor	%xmm6,%xmm4
	xorl	%r9d,%r12d
	rorl	$11,%r14d
	movdqa	%xmm7,%xmm6
	xorl	%eax,%edi
	addl	%r12d,%r10d
	pxor	%xmm5,%xmm4
	rorl	$6,%r13d
	andl	%edi,%r15d
	xorl	%r11d,%r14d
	psrld	$10,%xmm7
	addl	%r13d,%r10d
	xorl	%eax,%r15d
	paddd	%xmm4,%xmm2
	rorl	$2,%r14d
	addl	%r10d,%ecx
	psrlq	$17,%xmm6
	addl	%r15d,%r10d
	movl	%ecx,%r13d
	addl	%r10d,%r14d
	pxor	%xmm6,%xmm7
	rorl	$14,%r13d
	movl	%r14d,%r10d
	movl	%edx,%r12d
	rorl	$9,%r14d
	psrlq	$2,%xmm6
	xorl	%ecx,%r13d
	xorl	%r8d,%r12d
	pxor	%xmm6,%xmm7
	rorl	$5,%r13d
	xorl	%r10d,%r14d
	andl	%ecx,%r12d
	pshufd	$128,%xmm7,%xmm7
	xorl	%ecx,%r13d
	addl	40(%rsp),%r9d
	movl	%r10d,%r15d
	psrldq	$8,%xmm7
	xorl	%r8d,%r12d
	rorl	$11,%r14d
	xorl	%r11d,%r15d
	addl	%r12d,%r9d
	rorl	$6,%r13d
	paddd	%xmm7,%xmm2
	andl	%r15d,%edi
	xorl	%r10d,%r14d
	addl	%r13d,%r9d
	pshufd	$80,%xmm2,%xmm7
	xorl	%r11d,%edi
	rorl	$2,%r14d
	addl	%r9d,%ebx
	movdqa	%xmm7,%xmm6
	addl	%edi,%r9d
	movl	%ebx,%r13d
	psrld	$10,%xmm7
	addl	%r9d,%r14d
	rorl	$14,%r13d
	psrlq	$17,%xmm6
	movl	%r14d,%r9d
	movl	%ecx,%r12d
	pxor	%xmm6,%xmm7
	rorl	$9,%r14d
	xorl	%ebx,%r13d
	xorl	%edx,%r12d
	rorl	$5,%r13d
	xorl	%r9d,%r14d
	psrlq	$2,%xmm6
	andl	%ebx,%r12d
	xorl	%ebx,%r13d
	addl	44(%rsp),%r8d
	pxor	%xmm6,%xmm7
	movl	%r9d,%edi
	xorl	%edx,%r12d
	rorl	$11,%r14d
	pshufd	$8,%xmm7,%xmm7
	xorl	%r10d,%edi
	addl	%r12d,%r8d
	movdqa	32(%rsi),%xmm6
	rorl	$6,%r13d
	andl	%edi,%r15d
	pslldq	$8,%xmm7
	xorl	%r9d,%r14d
	addl	%r13d,%r8d
	xorl	%r10d,%r15d
	paddd	%xmm7,%xmm2
	rorl	$2,%r14d
	addl	%r8d,%eax
	addl	%r15d,%r8d
	paddd	%xmm2,%xmm6
	movl	%eax,%r13d
	addl	%r8d,%r14d
	movdqa	%xmm6,32(%rsp)
	rorl	$14,%r13d
	movdqa	%xmm0,%xmm4
	movl	%r14d,%r8d
	movl	%ebx,%r12d
	movdqa	%xmm2,%xmm7
	rorl	$9,%r14d
	xorl	%eax,%r13d
	xorl	%ecx,%r12d
	rorl	$5,%r13d
	xorl	%r8d,%r14d
.byte	102,15,58,15,227,4
	andl	%eax,%r12d
	xorl	%eax,%r13d
.byte	102,15,58,15,249,4
	addl	48(%rsp),%edx
	movl	%r8d,%r15d
	xorl	%ecx,%r12d
	rorl	$11,%r14d
	movdqa	%xmm4,%xmm5
	xorl	%r9d,%r15d
	addl	%r12d,%edx
	movdqa	%xmm4,%xmm6
	rorl	$6,%r13d
	andl	%r15d,%edi
	psrld	$3,%xmm4
	xorl	%r8d,%r14d
	addl	%r13d,%edx
	xorl	%r9d,%edi
	paddd	%xmm7,%xmm3
	rorl	$2,%r14d
	addl	%edx,%r11d
	psrld	$7,%xmm6
	addl	%edi,%edx
	movl	%r11d,%r13d
	pshufd	$250,%xmm2,%xmm7
	addl	%edx,%r14d
	rorl	$14,%r13d
	pslld	$14,%xmm5
	movl	%r14d,%edx
	movl	%eax,%r12d
	pxor	%xmm6,%xmm4
	rorl	$9,%r14d
	xorl	%r11d,%r13d
	xorl	%ebx,%r12d
	rorl	$5,%r13d
	psrld	$11,%xmm6
	xorl	%edx,%r14d
	pxor	%xmm5,%xmm4
	andl	%r11d,%r12d
	xorl	%r11d,%r13d
	pslld	$11,%xmm5
	addl	52(%rsp),%ecx
	movl	%edx,%edi
	pxor	%xmm6,%xmm4
	xorl	%ebx,%r12d
	rorl	$11,%r14d
	movdqa	%xmm7,%xmm6
	xorl	%r8d,%edi
	addl	%r12d,%ecx
	pxor	%xmm5,%xmm4
	rorl	$6,%r13d
	andl	%edi,%r15d
	xorl	%edx,%r14d
	psrld	$10,%xmm7
	addl	%r13d,%ecx
	xorl	%r8d,%r15d
	paddd	%xmm4,%xmm3
	rorl	$2,%r14d
	addl	%ecx,%r10d
	psrlq	$17,%xmm6
	addl	%r15d,%ecx
	movl	%r10d,%r13d
	addl	%ecx,%r14d
	pxor	%xmm6,%xmm7
	rorl	$14,%r13d
	movl	%r14d,%ecx
	movl	%r11d,%r12d
	rorl	$9,%r14d
	psrlq	$2,%xmm6
	xorl	%r10d,%r13d
	xorl	%eax,%r12d
	pxor	%xmm6,%xmm7
	rorl	$5,%r13d
	xorl	%ecx,%r14d
	andl	%r10d,%r12d
	pshufd	$128,%xmm7,%xmm7
	xorl	%r10d,%r13d
	addl	56(%rsp),%ebx
	movl	%ecx,%r15d
	psrldq	$8,%xmm7
	xorl	%eax,%r12d
	rorl	$11,%r14d
	xorl	%edx,%r15d
	addl	%r12d,%ebx
	rorl	$6,%r13d
	paddd	%xmm7,%xmm3
	andl	%r15d,%edi
	xorl	%ecx,%r14d
	addl	%r13d,%ebx
	pshufd	$80,%xmm3,%xmm7
	xorl	%edx,%edi
	rorl	$2,%r14d
	addl	%ebx,%r9d
	movdqa	%xmm7,%xmm6
	addl	%edi,%ebx
	movl	%r9d,%r13d
	psrld	$10,%xmm7
	addl	%ebx,%r14d
	rorl	$14,%r13d
	psrlq	$17,%xmm6
	movl	%r14d,%ebx
	movl	%r10d,%r12d
	pxor	%xmm6,%xmm7
	rorl	$9,%r14d
	xorl	%r9d,%r13d
	xorl	%r11d,%r12d
	rorl	$5,%r13d
	xorl	%ebx,%r14d
	psrlq	$2,%xmm6
	andl	%r9d,%r12d
	xorl	%r9d,%r13d
	addl	60(%rsp),%eax
	pxor	%xmm6,%xmm7
	movl	%ebx,%edi
	xorl	%r11d,%r12d
	rorl	$11,%r14d
	pshufd	$8,%xmm7,%xmm7
	xorl	%ecx,%edi
	addl	%r12d,%eax
	movdqa	48(%rsi),%xmm6
	rorl	$6,%r13d
	andl	%edi,%r15d
	pslldq	$8,%xmm7
	xorl	%ebx,%r14d
	addl	%r13d,%eax
	xorl	%ecx,%r15d
	paddd	%xmm7,%xmm3
	rorl	$2,%r14d
	addl	%eax,%r8d
	addl	%r15d,%eax
	paddd	%xmm3,%xmm6
	movl	%r8d,%r13d
	addl	%eax,%r14d
	movdqa	%xmm6,48(%rsp)
	cmpb	$0,67(%rsi)
	jne	.Lssse3_00_47
	rorl	$14,%r13d
	movl	%r14d,%eax
	movl	%r9d,%r12d
	rorl	$9,%r14d
	xorl	%r8d,%r13d
	xorl	%r10d,%r12d
	rorl	$5,%r13d
	xorl	%eax,%r14d
	andl	%r8d,%r12d
	xorl	%r8d,%r13d
	addl	0(%rsp),%r11d
	movl	%eax,%r15d
	xorl	%r10d,%r12d
	rorl	$11,%r14d
	xorl	%ebx,%r15d
	addl	%r12d,%r11d
	rorl	$6,%r13d
	andl	%r15d,%edi
	xorl	%eax,%r14d
	addl	%r13d,%r11d
	xorl	%ebx,%edi
	rorl	$2,%r14d
	addl	%r11d,%edx
	addl	%edi,%r11d
	movl	%edx,%r13d
	addl	%r11d,%r14d
	rorl	$14,%r13d
	movl	%r14d,%r11d
	movl	%r8d,%r12d
	rorl	$9,%r14d
	xorl	%edx,%r13d
	xorl	%r9d,%r12d
	rorl	$5,%r13d
	xorl	%r11d,%r14d
	andl	%edx,%r12d
	xorl	%edx,%r13d
	addl	4(%rsp),%r10d
	movl	%r11d,%edi
	xorl	%r9d,%r12d
	rorl	$11,%r14d
	xorl	%eax,%edi
	addl	%r12d,%r10d
	rorl	$6,%r13d
	andl	%edi,%r15d
	xorl	%r11d,%r14d
	addl	%r13d,%r10d
	xorl	%eax,%r15d
	rorl	$2,%r14d
	addl	%r10d,%ecx
	addl	%r15d,%r10d
	movl	%ecx,%r13d
	addl	%r10d,%r14d
	rorl	$14,%r13d
	movl	%r14d,%r10d
	movl	%edx,%r12d
	rorl	$9,%r14d
	xorl	%ecx,%r13d
	xorl	%r8d,%r12d
	rorl	$5,%r13d
	xorl	%r10d,%r14d
	andl	%ecx,%r12d
	xorl	%ecx,%r13d
	addl	8(%rsp),%r9d
	movl	%r10d,%r15d
	xorl	%r8d,%r12d
	rorl	$11,%r14d
	xorl	%r11d,%r15d
	addl	%r12d,%r9d
	rorl	$6,%r13d
	andl	%r15d,%edi
	xorl	%r10d,%r14d
	addl	%r13d,%r9d
	xorl	%r11d,%edi
	rorl	$2,%r14d
	addl	%r9d,%ebx
	addl	%edi,%r9d
	movl	%ebx,%r13d
	addl	%r9d,%r14d
	rorl	$14,%r13d
	movl	%r14d,%r9d
	movl	%ecx,%r12d
	rorl	$9,%r14d
	xorl	%ebx,%r13d
	xorl	%edx,%r12d
	rorl	$5,%r13d
	xorl	%r9d,%r14d
	andl	%ebx,%r12d
	xorl	%ebx,%r13d
	addl	12(%rsp),%r8d
	movl	%r9d,%edi
	xorl	%edx,%r12d
	rorl	$11,%r14d
	xorl	%r10d,%edi
	addl	%r12d,%r8d
	rorl	$6,%r13d
	andl	%edi,%r15d
	xorl	%r9d,%r14d
	addl	%r13d,%r8d
	xorl	%r10d,%r15d
	rorl	$2,%r14d
	addl	%r8d,%eax
	addl	%r15d,%r8d
	movl	%eax,%r13d
	addl	%r8d,%r14d
	rorl	$14,%r13d
	movl	%r14d,%r8d
	movl	%ebx,%r12d
	rorl	$9,%r14d
	xorl	%eax,%r13d
	xorl	%ecx,%r12d
	rorl	$5,%r13d
	xorl	%r8d,%r14d
	andl	%eax,%r12d
	xorl	%eax,%r13d
	addl	16(%rsp),%edx
	movl	%r8d,%r15d
	xorl	%ecx,%r12d
	rorl	$11,%r14d
	xorl	%r9d,%r15d
	addl	%r12d,%edx
	rorl	$6,%r13d
	andl	%r15d,%edi
	xorl	%r8d,%r14d
	addl	%r13d,%edx
	xorl	%r9d,%edi
	rorl	$2,%r14d
	addl	%edx,%r11d
	addl	%edi,%edx
	movl	%r11d,%r13d
	addl	%edx,%r14d
	rorl	$14,%r13d
	movl	%r14d,%edx
	movl	%eax,%r12d
	rorl	$9,%r14d
	xorl	%r11d,%r13d
	xorl	%ebx,%r12d
	rorl	$5,%r13d
	xorl	%edx,%r14d
	andl	%r11d,%r12d
	xorl	%r11d,%r13d
	addl	20(%rsp),%ecx
	movl	%edx,%edi
	xorl	%ebx,%r12d
	rorl	$11,%r14d
	xorl	%r8d,%edi
	addl	%r12d,%ecx
	rorl	$6,%r13d
	andl	%edi,%r15d
	xorl	%edx,%r14d
	addl	%r13d,%ecx
	xorl	%r8d,%r15d
	rorl	$2,%r14d
	addl	%ecx,%r10d
	addl	%r15d,%ecx
	movl	%r10d,%r13d
	addl	%ecx,%r14d
	rorl	$14,%r13d
	movl	%r14d,%ecx
	movl	%r11d,%r12d
	rorl	$9,%r14d
	xorl	%r10d,%r13d
	xorl	%eax,%r12d
	rorl	$5,%r13d
	xorl	%ecx,%r14d
	andl	%r10d,%r12d
	xorl	%r10d,%r13d
	addl	24(%rsp),%ebx
	movl	%ecx,%r15d
	xorl	%eax,%r12d
	rorl	$11,%r14d
	xorl	%edx,%r15d
	addl	%r12d,%ebx
	rorl	$6,%r13d
	andl	%r15d,%edi
	xorl	%ecx,%r14d
	addl	%r13d,%ebx
	xorl	%edx,%edi
	rorl	$2,%r14d
	addl	%ebx,%r9d
	addl	%edi,%ebx
	movl	%r9d,%r13d
	addl	%ebx,%r14d
	rorl	$14,%r13d
	movl	%r14d,%ebx
	movl	%r10d,%r12d
	rorl	$9,%r14d
	xorl	%r9d,%r13d
	xorl	%r11d,%r12d
	rorl	$5,%r13d
	xorl	%ebx,%r14d
	andl	%r9d,%r12d
	xorl	%r9d,%r13d
	addl	28(%rsp),%eax
	movl	%ebx,%edi
	xorl	%r11d,%r12d
	rorl	$11,%r14d
	xorl	%ecx,%edi
	addl	%r12d,%eax
	rorl	$6,%r13d
	andl	%edi,%r15d
	xorl	%ebx,%r14d
	addl	%r13d,%eax
	xorl	%ecx,%r15d
	rorl	$2,%r14d
	addl	%eax,%r8d
	addl	%r15d,%eax
	movl	%r8d,%r13d
	addl	%eax,%r14d
	rorl	$14,%r13d
	movl	%r14d,%eax
	movl	%r9d,%r12d
	rorl	$9,%r14d
	xorl	%r8d,%r13d
	xorl	%r10d,%r12d
	rorl	$5,%r13d
	xorl	%eax,%r14d
	andl	%r8d,%r12d
	xorl	%r8d,%r13d
	addl	32(%rsp),%r11d
	movl	%eax,%r15d
	xorl	%r10d,%r12d
	rorl	$11,%r14d
	xorl	%ebx,%r15d
	addl	%r12d,%r11d
	rorl	$6,%r13d
	andl	%r15d,%edi
	xorl	%eax,%r14d
	addl	%r13d,%r11d
	xorl	%ebx,%edi
	rorl	$2,%r14d
	addl	%r11d,%edx
	addl	%edi,%r11d
	movl	%edx,%r13d
	addl	%r11d,%r14d
	rorl	$14,%r13d
	movl	%r14d,%r11d
	movl	%r8d,%r12d
	rorl	$9,%r14d
	xorl	%edx,%r13d
	xorl	%r9d,%r12d
	rorl	$5,%r13d
	xorl	%r11d,%r14d
	andl	%edx,%r12d
	xorl	%edx,%r13d
	addl	36(%rsp),%r10d
	movl	%r11d,%edi
	xorl	%r9d,%r12d
	rorl	$11,%r14d
	xorl	%eax,%edi
	addl	%r12d,%r10d
	rorl	$6,%r13d
	andl	%edi,%r15d
	xorl	%r11d,%r14d
	addl	%r13d,%r10d
	xorl	%eax,%r15d
	rorl	$2,%r14d
	addl	%r10d,%ecx
	addl	%r15d,%r10d
	movl	%ecx,%r13d
	addl	%r10d,%r14d
	rorl	$14,%r13d
	movl	%r14d,%r10d
	movl	%edx,%r12d
	rorl	$9,%r14d
	xorl	%ecx,%r13d
	xorl	%r8d,%r12d
	rorl	$5,%r13d
	xorl	%r10d,%r14d
	andl	%ecx,%r12d
	xorl	%ecx,%r13d
	addl	40(%rsp),%r9d
	movl	%r10d,%r15d
	xorl	%r8d,%r12d
	rorl	$11,%r14d
	xorl	%r11d,%r15d
	addl	%r12d,%r9d
	rorl	$6,%r13d
	andl	%r15d,%edi
	xorl	%r10d,%r14d
	addl	%r13d,%r9d
	xorl	%r11d,%edi
	rorl	$2,%r14d
	addl	%r9d,%ebx
	addl	%edi,%r9d
	movl	%ebx,%r13d
	addl	%r9d,%r14d
	rorl	$14,%r13d
	movl	%r14d,%r9d
	movl	%ecx,%r12d
	rorl	$9,%r14d
	xorl	%ebx,%r13d
	xorl	%edx,%r12d
	rorl	$5,%r13d
	xorl	%r9d,%r14d
	andl	%ebx,%r12d
	xorl	%ebx,%r13d
	addl	44(%rsp),%r8d
	movl	%r9d,%edi
	xorl	%edx,%r12d
	rorl	$11,%r14d
	xorl	%r10d,%edi
	addl	%r12d,%r8d
	rorl	$6,%r13d
	andl	%edi,%r15d
	xorl	%r9d,%r14d
	addl	%r13d,%r8d
	xorl	%r10d,%r15d
	rorl	$2,%r14d
	addl	%r8d,%eax
	addl	%r15d,%r8d
	movl	%eax,%r13d
	addl	%r8d,%r14d
	rorl	$14,%r13d
	movl	%r14d,%r8d
	movl	%ebx,%r12d
	rorl	$9,%r14d
	xorl	%eax,%r13d
	xorl	%ecx,%r12d
	rorl	$5,%r13d
	xorl	%r8d,%r14d
	andl	%eax,%r12d
	xorl	%eax,%r13d
	addl	48(%rsp),%edx
	movl	%r8d,%r15d
	xorl	%ecx,%r12d
	rorl	$11,%r14d
	xorl	%r9d,%r15d
	addl	%r12d,%edx
	rorl	$6,%r13d
	andl	%r15d,%edi
	xorl	%r8d,%r14d
	addl	%r13d,%edx
	xorl	%r9d,%edi
	rorl	$2,%r14d
	addl	%edx,%r11d
	addl	%edi,%edx
	movl	%r11d,%r13d
	addl	%edx,%r14d
	rorl	$14,%r13d
	movl	%r14d,%edx
	movl	%eax,%r12d
	rorl	$9,%r14d
	xorl	%r11d,%r13d
	xorl	%ebx,%r12d
	rorl	$5,%r13d
	xorl	%edx,%r14d
	andl	%r11d,%r12d
	xorl	%r11d,%r13d
	addl	52(%rsp),%ecx
	movl	%edx,%edi
	xorl	%ebx,%r12d
	rorl	$11,%r14d
	xorl	%r8d,%edi
	addl	%r12d,%ecx
	rorl	$6,%r13d
	andl	%edi,%r15d
	xorl	%edx,%r14d
	addl	%r13d,%ecx
	xorl	%r8d,%r15d
	rorl	$2,%r14d
	addl	%ecx,%r10d
	addl	%r15d,%ecx
	movl	%r10d,%r13d
	addl	%ecx,%r14d
	rorl	$14,%r13d
	movl	%r14d,%ecx
	movl	%r11d,%r12d
	rorl	$9,%r14d
	xorl	%r10d,%r13d
	xorl	%eax,%r12d
	rorl	$5,%r13d
	xorl	%ecx,%r14d
	andl	%r10d,%r12d
	xorl	%r10d,%r13d
	addl	56(%rsp),%ebx
	movl	%ecx,%r15d
	xorl	%eax,%r12d
	rorl	$11,%r14d
	xorl	%edx,%r15d
	addl	%r12d,%ebx
	rorl	$6,%r13d
	andl	%r15d,%edi
	xorl	%ecx,%r14d
	addl	%r13d,%ebx
	xorl	%edx,%edi
	rorl	$2,%r14d
	addl	%ebx,%r9d
	addl	%edi,%ebx
	movl	%r9d,%r13d
	addl	%ebx,%r14d
	rorl	$14,%r13d
	movl	%r14d,%ebx
	movl	%r10d,%r12d
	rorl	$9,%r14d
	xorl	%r9d,%r13d
	xorl	%r11d,%r12d
	rorl	$5,%r13d
	xorl	%ebx,%r14d
	andl	%r9d,%r12d
	xorl	%r9d,%r13d
	addl	60(%rsp),%eax
	movl	%ebx,%edi
	xorl	%r11d,%r12d
	rorl	$11,%r14d
	xorl	%ecx,%edi
	addl	%r12d,%eax
	rorl	$6,%r13d
	andl	%edi,%r15d
	xorl	%ebx,%r14d
	addl	%r13d,%eax
	xorl	%ecx,%r15d
	rorl	$2,%r14d
	addl	%eax,%r8d
	addl	%r15d,%eax
	movl	%r8d,%r13d
	addl	%eax,%r14d
	movq	0(%rbp),%rdi
	movl	%r14d,%eax
	movq	8(%rbp),%rsi

	addl	0(%rdi),%eax
	addl	4(%rdi),%ebx
	addl	8(%rdi),%ecx
	addl	12(%rdi),%edx
	addl	16(%rdi),%r8d
	addl	20(%rdi),%r9d
	addl	24(%rdi),%r10d
	addl	28(%rdi),%r11d

	leaq	64(%rsi),%rsi
	cmpq	16(%rbp),%rsi

	movl	%eax,0(%rdi)
	movl	%ebx,4(%rdi)
	movl	%ecx,8(%rdi)
	movl	%edx,12(%rdi)
	movl	%r8d,16(%rdi)
	movl	%r9d,20(%rdi)
	movl	%r10d,24(%rdi)
	movl	%r11d,28(%rdi)
	jb	.Lloop_ssse3

	xorps	%xmm0,%xmm0
	leaq	104+48(%rbp),%r11

	movaps	%xmm0,0(%rsp)
	movaps	%xmm0,16(%rsp)
	movaps	%xmm0,32(%rsp)
	movaps	%xmm0,48(%rsp)
	movaps	32(%rbp),%xmm6
	movaps	48(%rbp),%xmm7
	movaps	64(%rbp),%xmm8
	movaps	80(%rbp),%xmm9
	movq	104(%rbp),%r15

	movq	-40(%r11),%r14

	movq	-32(%r11),%r13

	movq	-24(%r11),%r12

	movq	-16(%r11),%rbx

	movq	-8(%r11),%rbp

.LSEH_epilogue_blst_sha256_block_data_order:
	mov	8(%r11),%rdi
	mov	16(%r11),%rsi

	leaq	(%r11),%rsp
	.byte	0xf3,0xc3

.LSEH_end_blst_sha256_block_data_order:
.globl	blst_sha256_emit

.def	blst_sha256_emit;	.scl 2;	.type 32;	.endef
.p2align	4
blst_sha256_emit:
	.byte	0xf3,0x0f,0x1e,0xfa

	movq	0(%rdx),%r8
	movq	8(%rdx),%r9
	movq	16(%rdx),%r10
	bswapq	%r8
	movq	24(%rdx),%r11
	bswapq	%r9
	movl	%r8d,4(%rcx)
	bswapq	%r10
	movl	%r9d,12(%rcx)
	bswapq	%r11
	movl	%r10d,20(%rcx)
	shrq	$32,%r8
	movl	%r11d,28(%rcx)
	shrq	$32,%r9
	movl	%r8d,0(%rcx)
	shrq	$32,%r10
	movl	%r9d,8(%rcx)
	shrq	$32,%r11
	movl	%r10d,16(%rcx)
	movl	%r11d,24(%rcx)
	.byte	0xf3,0xc3


.globl	blst_sha256_bcopy

.def	blst_sha256_bcopy;	.scl 2;	.type 32;	.endef
.p2align	4
blst_sha256_bcopy:
	.byte	0xf3,0x0f,0x1e,0xfa

	subq	%rdx,%rcx
.Loop_bcopy:
	movzbl	(%rdx),%eax
	leaq	1(%rdx),%rdx
	movb	%al,-1(%rcx,%rdx,1)
	decq	%r8
	jnz	.Loop_bcopy
	.byte	0xf3,0xc3


.globl	blst_sha256_hcopy

.def	blst_sha256_hcopy;	.scl 2;	.type 32;	.endef
.p2align	4
blst_sha256_hcopy:
	.byte	0xf3,0x0f,0x1e,0xfa

	movq	0(%rdx),%r8
	movq	8(%rdx),%r9
	movq	16(%rdx),%r10
	movq	24(%rdx),%r11
	movq	%r8,0(%rcx)
	movq	%r9,8(%rcx)
	movq	%r10,16(%rcx)
	movq	%r11,24(%rcx)
	.byte	0xf3,0xc3

.section	.pdata
.p2align	2
.rva	.LSEH_begin_blst_sha256_block_data_order_shaext
.rva	.LSEH_body_blst_sha256_block_data_order_shaext
.rva	.LSEH_info_blst_sha256_block_data_order_shaext_prologue

.rva	.LSEH_body_blst_sha256_block_data_order_shaext
.rva	.LSEH_epilogue_blst_sha256_block_data_order_shaext
.rva	.LSEH_info_blst_sha256_block_data_order_shaext_body

.rva	.LSEH_epilogue_blst_sha256_block_data_order_shaext
.rva	.LSEH_end_blst_sha256_block_data_order_shaext
.rva	.LSEH_info_blst_sha256_block_data_order_shaext_epilogue

.rva	.LSEH_begin_blst_sha256_block_data_order
.rva	.LSEH_body_blst_sha256_block_data_order
.rva	.LSEH_info_blst_sha256_block_data_order_prologue

.rva	.LSEH_body_blst_sha256_block_data_order
.rva	.LSEH_epilogue_blst_sha256_block_data_order
.rva	.LSEH_info_blst_sha256_block_data_order_body

.rva	.LSEH_epilogue_blst_sha256_block_data_order
.rva	.LSEH_end_blst_sha256_block_data_order
.rva	.LSEH_info_blst_sha256_block_data_order_epilogue

.section	.xdata
.p2align	3
.LSEH_info_blst_sha256_block_data_order_shaext_prologue:
.byte	1,0,5,0x0b
.byte	0,0x74,1,0
.byte	0,0x64,2,0
.byte	0,0x03
.byte	0,0
.LSEH_info_blst_sha256_block_data_order_shaext_body:
.byte	1,0,15,0
.byte	0x00,0x68,0x00,0x00
.byte	0x00,0x78,0x01,0x00
.byte	0x00,0x88,0x02,0x00
.byte	0x00,0x98,0x03,0x00
.byte	0x00,0xa8,0x04,0x00
.byte	0x00,0x74,0x0c,0x00
.byte	0x00,0x64,0x0d,0x00
.byte	0x00,0xa2
.byte	0x00,0x00,0x00,0x00,0x00,0x00
.LSEH_info_blst_sha256_block_data_order_shaext_epilogue:
.byte	1,0,5,11
.byte	0x00,0x74,0x01,0x00
.byte	0x00,0x64,0x02,0x00
.byte	0x00,0x03
.byte	0x00,0x00

.LSEH_info_blst_sha256_block_data_order_prologue:
.byte	1,0,5,0x0b
.byte	0,0x74,1,0
.byte	0,0x64,2,0
.byte	0,0x03
.byte	0,0
.LSEH_info_blst_sha256_block_data_order_body:
.byte	1,0,26,5
.byte	0x00,0x68,0x02,0x00
.byte	0x00,0x78,0x03,0x00
.byte	0x00,0x88,0x04,0x00
.byte	0x00,0x98,0x05,0x00
.byte	0x00,0xf4,0x0d,0x00
.byte	0x00,0xe4,0x0e,0x00
.byte	0x00,0xd4,0x0f,0x00
.byte	0x00,0xc4,0x10,0x00
.byte	0x00,0x34,0x11,0x00
.byte	0x00,0x74,0x14,0x00
.byte	0x00,0x64,0x15,0x00
.byte	0x00,0x03
.byte	0x00,0x01,0x12,0x00
.byte	0x00,0x50
.LSEH_info_blst_sha256_block_data_order_epilogue:
.byte	1,0,5,11
.byte	0x00,0x74,0x01,0x00
.byte	0x00,0x64,0x02,0x00
.byte	0x00,0x03
.byte	0x00,0x00

