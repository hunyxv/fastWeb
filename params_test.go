package fastweb

import (
	"testing"

	"github.com/valyala/fasthttp"
)

type req struct {
	Username string `valid:"username,required,minlength=1,maxlength=18"`
	Password string `valid:"passwd,required,minlength=6,maxlength=18"`
	Age      int    `valid:"age"`
	Email    string `valid:"e-mail,required,re=^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$"`
}

func TestScan(t *testing.T) {
	r := &req{}
	p, err := scan(r)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", p)
}

func TestPadding(t *testing.T) {
	r := &req{}
	p, err := scan(r)
	if err != nil {
		t.Fatal(err)
	}

	args := fasthttp.AcquireArgs()
	args.Set("username", "zhangsan")
	args.Set("passwd", "password12345")
	args.Set("age", "18")
	args.Set("e-mail", "1334435$#3djsd@gmail.com")

	args.VisitAll(func(key, val []byte){
		err = p.padding(key, val, r)
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := p.valid(r); err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", r)
}