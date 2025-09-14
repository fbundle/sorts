package el_sorts

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fbundle/sorts/form"
)

var TypeErr = fmt.Errorf("type_err")

func mustMatchHead(H form.Name, list form.List) {
	if H != list[0] {
		panic(TypeErr)
	}
}

func not(b bool) bool {
	return !b
}

func mustTermOf(ctx Context, x Sort, X Sort) {
	if !ctx.LessEqual(ctx.Parent(x), X) {
		panic(TypeErr)
	}
}

var rs = rand.New(rand.NewSource(time.Now().UnixNano()))

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rs.Intn(len(letters))]
	}
	return string(b)
}
