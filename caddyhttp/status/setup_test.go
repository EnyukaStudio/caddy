package status

import (
	"testing"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func TestSetup(t *testing.T) {
	c := caddy.NewTestController("http", `status 404 /foo`)
	err := setup(c)
	if err != nil {
		t.Errorf("Expected no errors, but got: %v", err)
	}
	mids := httpserver.GetConfig(c).Middleware()
	if len(mids) == 0 {
		t.Fatal("Expected middleware, had 0 instead")
	}

	handler := mids[0](httpserver.EmptyNext)
	myHandler, ok := handler.(Status)
	if !ok {
		t.Fatalf("Expected handler to be type Status, got: %#v",
			handler)
	}

	if !httpserver.SameNext(myHandler.Next, httpserver.EmptyNext) {
		t.Error("'Next' field of handler was not set properly")
	}
}

func TestStatusParse(t *testing.T) {
	tests := []struct {
		input     string
		shouldErr bool
		expected  map[string]int
	}{
		{`status`, true, map[string]int{}},
		{`status /foo`, true, map[string]int{}},
		{`status bar /foo`, true, map[string]int{}},
		{`status 404 /foo bar`, true, map[string]int{}},
		{`status 404 /foo`, false, map[string]int{"/foo": 404}},
		{`status {
		 }`,
			true,
			map[string]int{},
		},
		{`status 404 {
		 }`,
			true,
			map[string]int{},
		},
		{`status 404 {
			/foo
			/foo
		 }`,
			true,
			map[string]int{},
		},
		{`status 404 {
			404 /foo
		 }`,
			true,
			map[string]int{},
		},
		{`status 404 {
			/foo
			/bar
		 }`,
			false,
			map[string]int{"/foo": 404, "/bar": 404},
		},
	}

	for i, test := range tests {
		actual, err := statusParse(caddy.NewTestController("http",
			test.input))

		if err == nil && test.shouldErr {
			t.Errorf("Test %d didn't error, but it should have", i)
		} else if err != nil && !test.shouldErr {
			t.Errorf("Test %d errored, but it shouldn't have; got '%v'",
				i, err)
		} else if err != nil && test.shouldErr {
			continue
		}

		if len(actual) != len(test.expected) {
			t.Fatalf("Test %d expected %d rules, but got %d",
				i, len(test.expected), len(actual))
		}

		for expectedPath, expectedCode := range test.expected {
			actualCode, ok := actual[expectedPath]
			if !ok {
				t.Fatalf("Test %d: Path '%s' not found in parsed rules",
					i, expectedPath)
			}

			if actualCode != expectedCode {
				t.Errorf("Test %d: Expected status code %d for path '%s'. Got %d",
					i, expectedCode, expectedPath, actualCode)
			}
		}
	}
}
