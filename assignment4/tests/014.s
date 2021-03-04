global _main
extern _sub_ints
extern _sub_floats
extern _has_size
extern _sepia
extern _blur
extern _resize
extern _crop
extern _read_image
extern _print
extern _write_image
extern _show
extern _fail_assertion
extern _get_time

section .data
const0: db `This should not be printed`, 0
const1: dq 0

section .text

_main:
	push rbp
	mov rbp, rsp
	sub rsp, 16
	mov dword [rbp - 4], 1
	cmp qword [rbp - 4], 0
	jne .jump1
	mov rdi, [rel const0] ; This should not be printed
	call _fail_assertion
.jump1:
	mov rbx, [rel const1] ; 0
	mov [rbp - 12], rbx
	mov rax, [rbp - 12]
	add rsp, 16
	pop rbp
	ret
Compilation succeeded: assembly complete
