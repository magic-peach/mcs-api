/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package conformance

import "testing"

func TestFirstLine(t *testing.T) {
	cases := map[string]struct {
		in   string
		want string
	}{
		"single line":      {in: "just one line", want: "just one line"},
		"multiple lines":   {in: "first\nsecond\nthird", want: "first"},
		"leading trailing": {in: "  padded  \nrest", want: "padded"},
		"empty":            {in: "", want: ""},
		"only newline":     {in: "\nsecond", want: ""},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := firstLine(c.in); got != c.want {
				t.Errorf("firstLine(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestLookupSpecRef(t *testing.T) {
	original := specRefRegistry
	specRefRegistry = map[string]string{
		"exports a service": "kep-1645-ref-1",
	}
	t.Cleanup(func() { specRefRegistry = original })

	if got := lookupSpecRef("a test that exports a service correctly"); got != "kep-1645-ref-1" {
		t.Errorf("matching text: got %q, want %q", got, "kep-1645-ref-1")
	}
	if got := lookupSpecRef("a totally unrelated test"); got != "" {
		t.Errorf("no match: got %q, want empty string", got)
	}
}

func TestParseFailureMessage(t *testing.T) {
	withDesc := "my description\nUnexpected error:\n    <*errors.errorString | 0x1>: \n    \"foo\" not found\n    {s: \"foo not found\"}\noccurred"
	if got := parseFailureMessage(withDesc); got != `my description: "foo" not found` {
		t.Errorf("with description: got %q", got)
	}

	noDesc := "Unexpected error:\n    <*errors.errorString | 0x1>: \n    boom\n    {s: \"boom\"}\noccurred"
	if got := parseFailureMessage(noDesc); got != "boom" {
		t.Errorf("without description: got %q, want %q", got, "boom")
	}

	descWithPeriod := "context deadline exceeded.\nUnexpected error:\n    <*errors.errorString | 0x1>: \n    timeout\n    {s: \"timeout\"}\noccurred"
	if got := parseFailureMessage(descWithPeriod); got != "context deadline exceeded: timeout" {
		t.Errorf("description ending in a period: got %q", got)
	}

	fallback := "some plain multi line\nfailure text\nwith no gomega markers"
	if got := parseFailureMessage(fallback); got != "some plain multi line" {
		t.Errorf("fallback to first line: got %q, want %q", got, "some plain multi line")
	}
}
