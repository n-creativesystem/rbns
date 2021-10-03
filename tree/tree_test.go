package tree

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestInsertTree(t *testing.T) {
	tree := NewTree()
	cases := []struct {
		method      string
		url         string
		permissions []string
	}{
		{
			method:      http.MethodGet,
			url:         "/",
			permissions: []string{"root:get"},
		},
		{
			method:      http.MethodPost,
			url:         "/",
			permissions: []string{"root:post"},
		},
		{
			method:      http.MethodPut,
			url:         "/",
			permissions: []string{"root:put"},
		},
		{
			method:      http.MethodDelete,
			url:         "/",
			permissions: []string{"root:delete"},
		},
	}
	for _, c := range cases {
		if err := tree.Insert(c.method, c.url, c.permissions...); err != nil {
			t.Errorf("err: %v\n", err)
		}
		p := tree.Search(c.method, c.url)
		if !reflect.DeepEqual(c.permissions, p) {
			t.Errorf("err: deepEqual %s", strings.Join(p, ","))
		}
	}
}

func TestSearch(t *testing.T) {
	tree := NewTree()
	routes := [...]string{
		"/",
		"/cmd/:tool/",
		"/cmd/:tool/:sub",
		"/cmd/whoami",
		"/cmd/whoami/root",
		"/cmd/whoami/root/",
		"/search/",
		"/search/:query",
		"/search/gin-gonic",
		"/search/google",
		"/user_:name",
		"/user_:name/about",
		"/doc/",
		"/doc/go_faq.html",
		"/doc/go1.html",
		"/info/:user/public",
		"/info/:user/project/:project",
		"/info/:user/project/golang",
		"/:cc",
		"/:cc/cc",
		"/:cc/:dd/ee",
		"/:cc/:dd/:ee/ff",
		"/:cc/:dd/:ee/:ff/gg",
		"/:cc/:dd/:ee/:ff/:gg/hh",
		"/get/test/abc/",
		"/get/:param/abc/",
		"/something/:paramname/thirdthing",
		"/something/secondthing/test",
		"/get/abc",
		"/get/:param",
		"/get/abc/123abc",
		"/get/abc/:param",
		"/get/abc/123abc/xxx8",
		"/get/abc/123abc/:param",
		"/get/abc/123abc/xxx8/1234",
		"/get/abc/123abc/xxx8/:param",
		"/get/abc/123abc/xxx8/1234/ffas",
		"/get/abc/123abc/xxx8/1234/:param",
		"/get/abc/123abc/xxx8/1234/kkdd/12c",
		"/get/abc/123abc/xxx8/1234/kkdd/:param",
		"/get/abc/:param/test",
		"/get/abc/123abd/:param",
		"/get/abc/123abddd/:param",
		"/get/abc/123/:param",
		"/get/abc/123abg/:param",
		"/get/abc/123abf/:param",
		"/get/abc/123abfff/:param",
		"/files/:dir/:filepath",
		"/aa/:xx",
	}
	addRoute := func(method, url string) {
		_ = tree.Insert(method, url, fmt.Sprintf("route:%s", url))
	}
	for _, route := range routes {
		addRoute(http.MethodGet, route)
	}
	type testRequests []struct {
		path       string
		permission string
	}
	v := testRequests{
		// {"/", "/"},
		// {"/cmd/test", "/cmd/:tool/"},
		// {"/cmd/test/", "/cmd/:tool/"},
		// {"/cmd/test/3", "/cmd/:tool/:sub"},
		// {"/cmd/who", "/cmd/:tool/"},
		// {"/cmd/who/", "/cmd/:tool/"},
		// {"/cmd/whoami", "/cmd/whoami"},
		// {"/cmd/whoami/", "/cmd/whoami"},
		// {"/cmd/whoami/r", ""},
		// {"/cmd/whoami/r/", ""},
		// {"/cmd/whoami/root", "/cmd/whoami/root/"},
		// {"/cmd/whoami/root/", "/cmd/whoami/root/"},
		// {"/src/", "/:cc"},
		// {"/src/some/file.png", ""},
		// {"/search/", "/search/"},
		// {"/search/someth!ng+in+ünìcodé", "/search/:query"},
		// {"/search/someth!ng+in+ünìcodé/", "/search/:query"},
		// {"/search/gin", "/search/:query"},
		// {"/search/gin-gonic", "/search/gin-gonic"},
		// {"/search/google", "/search/google"},
		{"/files/js/inc/framework.js", "/files/:dir/:filepath"},
		{"/info/gordon/public", "/info/:user/public"},
		{"/info/gordon/project/go", "/info/:user/project/:project"},
		{"/info/gordon/project/golang", "/info/:user/project/golang"},
		{"/aa/aa", "/aa/:xx"},
		{"/ab/ab", "/ab/:xx"},
		{"/a", "/:cc"},
		{"/all", "/:cc"},
		{"/d", "/:cc"},
		{"/ad", "/:cc"},
		{"/dd", "/:cc"},
		{"/dddaa", "/:cc"},
		{"/aa", "/:cc"},
		{"/aaa", "/:cc"},
		{"/aaa/cc", "/:cc/cc"},
		{"/ab", "/:cc"},
		{"/abb", "/:cc"},
		{"/abb/cc", "/:cc/cc"},
		{"/allxxxx", "/:cc"},
		{"/alldd", "/:cc"},
		{"/all/cc", "/:cc/cc"},
		{"/a/cc", "/:cc/cc"},
		{"/cc/cc", "/:cc/cc"},
		{"/ccc/cc", "/:cc/cc"},
		{"/deedwjfs/cc", "/:cc/cc"},
		{"/acllcc/cc", "/:cc/cc"},
		{"/get/test/abc/", "/get/test/abc/"},
		{"/get/te/abc/", "/get/:param/abc/"},
		{"/get/testaa/abc/", "/get/:param/abc/"},
		{"/get/xx/abc/", "/get/:param/abc/"},
		{"/get/tt/abc/", "/get/:param/abc/"},
		{"/get/a/abc/", "/get/:param/abc/"},
		{"/get/t/abc/", "/get/:param/abc/"},
		{"/get/aa/abc/", "/get/:param/abc/"},
		{"/get/abas/abc/", "/get/:param/abc/"},
		{"/something/secondthing/test", "/something/secondthing/test"},
		{"/something/abcdad/thirdthing", "/something/:paramname/thirdthing"},
		{"/something/secondthingaaaa/thirdthing", "/something/:paramname/thirdthing"},
		{"/something/se/thirdthing", "/something/:paramname/thirdthing"},
		{"/something/s/thirdthing", "/something/:paramname/thirdthing"},
		{"/c/d/ee", "/:cc/:dd/ee"},
		{"/c/d/e/ff", "/:cc/:dd/:ee/ff"},
		{"/c/d/e/f/gg", "/:cc/:dd/:ee/:ff/gg"},
		{"/c/d/e/f/g/hh", "/:cc/:dd/:ee/:ff/:gg/hh"},
		{"/cc/dd/ee/ff/gg/hh", "/:cc/:dd/:ee/:ff/:gg/hh"},
		{"/get/abc", "/get/abc"},
		{"/get/a", "/get/:param"},
		{"/get/abz", "/get/:param"},
		{"/get/12a", "/get/:param"},
		{"/get/abcd", "/get/:param"},
		{"/get/abc/123abc", "/get/abc/123abc"},
		{"/get/abc/12", "/get/abc/:param"},
		{"/get/abc/123ab", "/get/abc/:param"},
		{"/get/abc/xyz", "/get/abc/:param"},
		{"/get/abc/123abcddxx", "/get/abc/:param"},
		{"/get/abc/123abc/xxx8", "/get/abc/123abc/xxx8"},
		{"/get/abc/123abc/x", "/get/abc/123abc/:param"},
		{"/get/abc/123abc/xxx", "/get/abc/123abc/:param"},
		{"/get/abc/123abc/abc", "/get/abc/123abc/:param"},
		{"/get/abc/123abc/xxx8xxas", "/get/abc/123abc/:param"},
		{"/get/abc/123abc/xxx8/1234", "/get/abc/123abc/xxx8/1234"},
		{"/get/abc/123abc/xxx8/1", "/get/abc/123abc/xxx8/:param"},
		{"/get/abc/123abc/xxx8/123", "/get/abc/123abc/xxx8/:param"},
		{"/get/abc/123abc/xxx8/78k", "/get/abc/123abc/xxx8/:param"},
		{"/get/abc/123abc/xxx8/1234xxxd", "/get/abc/123abc/xxx8/:param"},
		{"/get/abc/123abc/xxx8/1234/ffas", "/get/abc/123abc/xxx8/1234/ffas"},
		{"/get/abc/123abc/xxx8/1234/f", "/get/abc/123abc/xxx8/1234/:param"},
		{"/get/abc/123abc/xxx8/1234/ffa", "/get/abc/123abc/xxx8/1234/:param"},
		{"/get/abc/123abc/xxx8/1234/kka", "/get/abc/123abc/xxx8/1234/:param"},
		{"/get/abc/123abc/xxx8/1234/ffas321", "/get/abc/123abc/xxx8/1234/:param"},
		{"/get/abc/123abc/xxx8/1234/kkdd/12c", "/get/abc/123abc/xxx8/1234/kkdd/12c"},
		{"/get/abc/123abc/xxx8/1234/kkdd/1", "/get/abc/123abc/xxx8/1234/kkdd/:param"},
		{"/get/abc/123abc/xxx8/1234/kkdd/12", "/get/abc/123abc/xxx8/1234/kkdd/:param"},
		{"/get/abc/123abc/xxx8/1234/kkdd/12b", "/get/abc/123abc/xxx8/1234/kkdd/:param"},
		{"/get/abc/123abc/xxx8/1234/kkdd/34", "/get/abc/123abc/xxx8/1234/kkdd/:param"},
		{"/get/abc/123abc/xxx8/1234/kkdd/12c2e3", "/get/abc/123abc/xxx8/1234/kkdd/:param"},
		{"/get/abc/12/test", "/get/abc/:param/test"},
		{"/get/abc/123abdd/test", "/get/abc/:param/test"},
		{"/get/abc/123abdddf/test", "/get/abc/:param/test"},
		{"/get/abc/123ab/test", "/get/abc/:param/test"},
		{"/get/abc/123abgg/test", "/get/abc/:param/test"},
		{"/get/abc/123abff/test", "/get/abc/:param/test"},
		{"/get/abc/123abffff/test", "/get/abc/:param/test"},
		{"/get/abc/123abd/test", "/get/abc/123abd/:param"},
		{"/get/abc/123abddd/test", "/get/abc/123abddd/:param"},
		{"/get/abc/123/test22", "/get/abc/123/:param"},
		{"/get/abc/123abg/test", "/get/abc/123abg/:param"},
		{"/get/abc/123abf/testss", "/get/abc/123abf/:param"},
		{"/get/abc/123abfff/te", "/get/abc/123abfff/:param"},
	}
	checkRoute := func(method, url, permission string) {
		p := tree.Search(method, url)
		permissions := []string{}
		if permission != "" {
			permissions = append(permissions, fmt.Sprintf("route:%s", permission))
		}
		if len(p) == 0 && len(permissions) == 0 {
			return
		}
		if !reflect.DeepEqual(p, permissions) {
			t.Errorf("err: deepEqual %s, %v, %v", url, p, permissions)
		}
	}
	for _, route := range v {
		checkRoute(http.MethodGet, route.path, route.permission)
	}
}
