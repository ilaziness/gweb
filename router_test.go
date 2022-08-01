package gweb

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAddRoute(t *testing.T) {
	r := newRouter()
	r.AddRoute("GET", "/", nil)
	r.AddRoute("GET", "/abc/ab", nil)
	r.AddRoute("GET", "/c", nil)
	r.AddRoute("GET", "/c/ab", nil)
	r.AddRoute("GET", "/d/ab", nil)
	r.AddRoute("GET", "/d", nil)

	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"/", "p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"/", "p", "*"})
	//ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"/", "p", "*name"})
	//fmt.Println(parsePattern("/p/:name"))
	//fmt.Println(parsePattern("/p/*"))
	//fmt.Println(parsePattern("/p/name/*"))
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestMatchRoute(t *testing.T) {
	r := newRouter()
	r.AddRoute("GET", "/hello/:name", nil)
	r.AddRoute("GET", "/hello/abc", nil)

	r.AddRoute("GET", "/:age", nil)
	r.AddRoute("GET", "/18", nil)

	n, ps := r.Match("GET", "/hello/abc")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	println(n.pattern)
	fmt.Printf("%+v\n", ps)
	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "abc" {
		t.Fatal("name should be equal to 'abc'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])

	r.AddRoute("GET", "/hello/*path", nil)
	n, ps = r.Match("GET", "/hello/cb/dfb/df")
	if n == nil {
		t.Fatal("2: nil shouldn't be returned")
	}
	println(n.pattern)
	fmt.Printf("%+v\n", ps)
	if n.pattern != "/hello/*path" {
		t.Fatal("should match /hello/*path")
	}
	if ps["path"] != "cb/dfb/df" {
		t.Fatal("name should be equal to 'abc'")
	}
}
