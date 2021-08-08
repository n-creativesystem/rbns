package model

import (
	"testing"
)

type testRequests []struct {
	path       string
	nilHandler bool
	route      string
}

func checkRequests(t *testing.T, tree *node, requests testRequests) {
	for _, request := range requests {
		value := tree.getValue(request.path)

		if value.handlers == nil {
			if !request.nilHandler {
				t.Errorf("handle mismatch for route '%s': Expected non-nil handle", request.path)
			}
		} else if request.nilHandler {
			t.Errorf("handle mismatch for route '%s': Expected nil handle", request.path)
		} else {
			if value.handlers == nil {
				t.Errorf("handle mismatch for route '%s': Wrong handle (%s)", request.path, request.route)
			}
		}
	}
}

func TestTreeAddAndGet(t *testing.T) {
	tree := &node{}

	routes := [...]string{
		"/hi",
		"/contact",
		"/co",
		"/c",
		"/a",
		"/ab",
		"/doc/",
		"/doc/go_faq.html",
		"/doc/go1.html",
		"/α",
		"/β",
	}
	for _, route := range routes {
		f := func() []string {
			return []string{"read:handler"}
		}
		tree.addRoute(route, f)
	}

	checkRequests(t, tree, testRequests{
		{"/a", false, "/a"},
		{"/", true, ""},
		{"/hi", false, "/hi"},
		{"/contact", false, "/contact"},
		{"/co", false, "/co"},
		{"/con", true, ""},  // key mismatch
		{"/cona", true, ""}, // key mismatch
		{"/no", true, ""},   // no matching child
		{"/ab", false, "/ab"},
		{"/α", false, "/α"},
		{"/β", false, "/β"},
	})
}

func TestTreeWildcard(t *testing.T) {
	tree := &node{}

	routes := [...]string{
		"/",
		"/cmd/:tool/",
		"/cmd/:tool/:sub",
		"/cmd/whoami",
		"/cmd/whoami/root",
		"/cmd/whoami/root/",
		"/src/*filepath",
		"/search/",
		"/search/:query",
		"/search/gin-gonic",
		"/search/google",
		"/user_:name",
		"/user_:name/about",
		"/files/:dir/*filepath",
		"/doc/",
		"/doc/go_faq.html",
		"/doc/go1.html",
		"/info/:user/public",
		"/info/:user/project/:project",
		"/info/:user/project/golang",
		"/aa/*xx",
		"/ab/*xx",
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
	}
	for _, route := range routes {
		f := func() []string {
			return []string{"read:handler"}
		}
		tree.addRoute(route, f)
	}

	checkRequests(t, tree, testRequests{
		{"/", false, "/"},
		{"/cmd/test", true, "/cmd/:tool/"},
		{"/cmd/test/", false, "/cmd/:tool/"},
		{"/cmd/test/3", false, "/cmd/:tool/:sub"},
		{"/cmd/who", true, "/cmd/:tool/"},
		{"/cmd/who/", false, "/cmd/:tool/"},
		{"/cmd/whoami", false, "/cmd/whoami"},
		{"/cmd/whoami/", true, "/cmd/whoami"},
		{"/cmd/whoami/r", false, "/cmd/:tool/:sub"},
		{"/cmd/whoami/r/", true, "/cmd/:tool/:sub"},
		{"/cmd/whoami/root", false, "/cmd/whoami/root"},
		{"/cmd/whoami/root/", false, "/cmd/whoami/root/"},
		{"/src/", false, "/src/*filepath"},
		{"/src/some/file.png", false, "/src/*filepath"},
		{"/search/", false, "/search/"},
		{"/search/someth!ng+in+ünìcodé", false, "/search/:query"},
		{"/search/someth!ng+in+ünìcodé/", true, ""},
		{"/search/gin", false, "/search/:query"},
		{"/search/gin-gonic", false, "/search/gin-gonic"},
		{"/search/google", false, "/search/google"},
		{"/user_gopher", false, "/user_:name"},
		{"/user_gopher/about", false, "/user_:name/about"},
		{"/files/js/inc/framework.js", false, "/files/:dir/*filepath"},
		{"/info/gordon/public", false, "/info/:user/public"},
		{"/info/gordon/project/go", false, "/info/:user/project/:project"},
		{"/info/gordon/project/golang", false, "/info/:user/project/golang"},
		{"/aa/aa", false, "/aa/*xx"},
		{"/ab/ab", false, "/ab/*xx"},
		{"/a", false, "/:cc"},
		// * Error with argument being intercepted
		// new PR handle (/all /all/cc /a/cc)
		// fix PR: https://github.com/gin-gonic/gin/pull/2796
		{"/all", false, "/:cc"},
		{"/d", false, "/:cc"},
		{"/ad", false, "/:cc"},
		{"/dd", false, "/:cc"},
		{"/dddaa", false, "/:cc"},
		{"/aa", false, "/:cc"},
		{"/aaa", false, "/:cc"},
		{"/aaa/cc", false, "/:cc/cc"},
		{"/ab", false, "/:cc"},
		{"/abb", false, "/:cc"},
		{"/abb/cc", false, "/:cc/cc"},
		{"/allxxxx", false, "/:cc"},
		{"/alldd", false, "/:cc"},
		{"/all/cc", false, "/:cc/cc"},
		{"/a/cc", false, "/:cc/cc"},
		{"/cc/cc", false, "/:cc/cc"},
		{"/ccc/cc", false, "/:cc/cc"},
		{"/deedwjfs/cc", false, "/:cc/cc"},
		{"/acllcc/cc", false, "/:cc/cc"},
		{"/get/test/abc/", false, "/get/test/abc/"},
		{"/get/te/abc/", false, "/get/:param/abc/"},
		{"/get/testaa/abc/", false, "/get/:param/abc/"},
		{"/get/xx/abc/", false, "/get/:param/abc/"},
		{"/get/tt/abc/", false, "/get/:param/abc/"},
		{"/get/a/abc/", false, "/get/:param/abc/"},
		{"/get/t/abc/", false, "/get/:param/abc/"},
		{"/get/aa/abc/", false, "/get/:param/abc/"},
		{"/get/abas/abc/", false, "/get/:param/abc/"},
		{"/something/secondthing/test", false, "/something/secondthing/test"},
		{"/something/abcdad/thirdthing", false, "/something/:paramname/thirdthing"},
		{"/something/secondthingaaaa/thirdthing", false, "/something/:paramname/thirdthing"},
		{"/something/se/thirdthing", false, "/something/:paramname/thirdthing"},
		{"/something/s/thirdthing", false, "/something/:paramname/thirdthing"},
		{"/c/d/ee", false, "/:cc/:dd/ee"},
		{"/c/d/e/ff", false, "/:cc/:dd/:ee/ff"},
		{"/c/d/e/f/gg", false, "/:cc/:dd/:ee/:ff/gg"},
		{"/c/d/e/f/g/hh", false, "/:cc/:dd/:ee/:ff/:gg/hh"},
		{"/cc/dd/ee/ff/gg/hh", false, "/:cc/:dd/:ee/:ff/:gg/hh"},
		{"/get/abc", false, "/get/abc"},
		{"/get/a", false, "/get/:param"},
		{"/get/abz", false, "/get/:param"},
		{"/get/12a", false, "/get/:param"},
		{"/get/abcd", false, "/get/:param"},
		{"/get/abc/123abc", false, "/get/abc/123abc"},
		{"/get/abc/12", false, "/get/abc/:param"},
		{"/get/abc/123ab", false, "/get/abc/:param"},
		{"/get/abc/xyz", false, "/get/abc/:param"},
		{"/get/abc/123abcddxx", false, "/get/abc/:param"},
		{"/get/abc/123abc/xxx8", false, "/get/abc/123abc/xxx8"},
		{"/get/abc/123abc/x", false, "/get/abc/123abc/:param"},
		{"/get/abc/123abc/xxx", false, "/get/abc/123abc/:param"},
		{"/get/abc/123abc/abc", false, "/get/abc/123abc/:param"},
		{"/get/abc/123abc/xxx8xxas", false, "/get/abc/123abc/:param"},
		{"/get/abc/123abc/xxx8/1234", false, "/get/abc/123abc/xxx8/1234"},
		{"/get/abc/123abc/xxx8/1", false, "/get/abc/123abc/xxx8/:param"},
		{"/get/abc/123abc/xxx8/123", false, "/get/abc/123abc/xxx8/:param"},
		{"/get/abc/123abc/xxx8/78k", false, "/get/abc/123abc/xxx8/:param"},
		{"/get/abc/123abc/xxx8/1234xxxd", false, "/get/abc/123abc/xxx8/:param"},
		{"/get/abc/123abc/xxx8/1234/ffas", false, "/get/abc/123abc/xxx8/1234/ffas"},
		{"/get/abc/123abc/xxx8/1234/f", false, "/get/abc/123abc/xxx8/1234/:param"},
		{"/get/abc/123abc/xxx8/1234/ffa", false, "/get/abc/123abc/xxx8/1234/:param"},
		{"/get/abc/123abc/xxx8/1234/kka", false, "/get/abc/123abc/xxx8/1234/:param"},
		{"/get/abc/123abc/xxx8/1234/ffas321", false, "/get/abc/123abc/xxx8/1234/:param"},
		{"/get/abc/123abc/xxx8/1234/kkdd/12c", false, "/get/abc/123abc/xxx8/1234/kkdd/12c"},
		{"/get/abc/123abc/xxx8/1234/kkdd/1", false, "/get/abc/123abc/xxx8/1234/kkdd/:param"},
		{"/get/abc/123abc/xxx8/1234/kkdd/12", false, "/get/abc/123abc/xxx8/1234/kkdd/:param"},
		{"/get/abc/123abc/xxx8/1234/kkdd/12b", false, "/get/abc/123abc/xxx8/1234/kkdd/:param"},
		{"/get/abc/123abc/xxx8/1234/kkdd/34", false, "/get/abc/123abc/xxx8/1234/kkdd/:param"},
		{"/get/abc/123abc/xxx8/1234/kkdd/12c2e3", false, "/get/abc/123abc/xxx8/1234/kkdd/:param"},
		{"/get/abc/12/test", false, "/get/abc/:param/test"},
		{"/get/abc/123abdd/test", false, "/get/abc/:param/test"},
		{"/get/abc/123abdddf/test", false, "/get/abc/:param/test"},
		{"/get/abc/123ab/test", false, "/get/abc/:param/test"},
		{"/get/abc/123abgg/test", false, "/get/abc/:param/test"},
		{"/get/abc/123abff/test", false, "/get/abc/:param/test"},
		{"/get/abc/123abffff/test", false, "/get/abc/:param/test"},
		{"/get/abc/123abd/test", false, "/get/abc/123abd/:param"},
		{"/get/abc/123abddd/test", false, "/get/abc/123abddd/:param"},
		{"/get/abc/123/test22", false, "/get/abc/123/:param"},
		{"/get/abc/123abg/test", false, "/get/abc/123abg/:param"},
		{"/get/abc/123abf/testss", false, "/get/abc/123abf/:param"},
		{"/get/abc/123abfff/te", false, "/get/abc/123abfff/:param"},
	})
}

func TestUnescapeParameters(t *testing.T) {
	tree := &node{}

	routes := [...]string{
		"/",
		"/cmd/:tool/:sub",
		"/cmd/:tool/",
		"/src/*filepath",
		"/search/:query",
		"/files/:dir/*filepath",
		"/info/:user/project/:project",
		"/info/:user",
	}
	for _, route := range routes {
		f := func() []string {
			return []string{"read:handler"}
		}
		tree.addRoute(route, f)
	}

	checkRequests(t, tree, testRequests{
		{"/", false, "/"},
		{"/cmd/test/", false, "/cmd/:tool/"},
		{"/cmd/test", true, ""},
		{"/src/some/file.png", false, "/src/*filepath"},
		{"/src/some/file+test.png", false, "/src/*filepath"},
		{"/src/some/file++++%%%%test.png", false, "/src/*filepath"},
		{"/src/some/file%2Ftest.png", false, "/src/*filepath"},
		{"/search/someth!ng+in+ünìcodé", false, "/search/:query"},
		{"/info/gordon/project/go", false, "/info/:user/project/:project"},
		{"/info/slash%2Fgordon", false, "/info/:user"},
		{"/info/slash%2Fgordon/project/Project%20%231", false, "/info/:user/project/:project"},
		{"/info/slash%%%%", false, "/info/:user"},
		{"/info/slash%%%%2Fgordon/project/Project%%%%20%231", false, "/info/:user/project/:project"},
	})
}
