#line 3 "C:\\Users\\stian\\Documents\\GitHub\\go\\src\\Heisprosjekt\\TEST\\test_go.go"

#include "test.h"



// Usual nonsense: if x and y are not equal, the type will be invalid
// (have a negative array count) and an inscrutable error will come
// out of the compiler and hopefully mention "name".
#define __cgo_compile_assert_eq(x, y, name) typedef char name[(x-y)*(x-y)*-2+1];

// Check at compile time that the sizes we use match our expectations.
#define __cgo_size_assert(t, n) __cgo_compile_assert_eq(sizeof(t), n, _cgo_sizeof_##t##_is_not_##n)

__cgo_size_assert(char, 1)
__cgo_size_assert(short, 2)
__cgo_size_assert(int, 4)
typedef long long __cgo_long_long;
__cgo_size_assert(__cgo_long_long, 8)
__cgo_size_assert(float, 4)
__cgo_size_assert(double, 8)

extern char* _cgo_topofstack(void);

#include <errno.h>
#include <string.h>

void
_cgo_56271688278c_Cfunc_test_func(void *v)
{
	struct {
		int r;
		char __pad4[4];
	} __attribute__((__packed__, __gcc_struct__)) *a = v;
	char *stktop = _cgo_topofstack();
	__typeof__(a->r) r = test_func();
	a = (void*)((char*)a + (_cgo_topofstack() - stktop));
	a->r = r;
}

