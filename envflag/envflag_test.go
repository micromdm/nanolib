package envflag

import (
	"flag"
	"testing"
)

type stringFlag struct {
	str string
}

func (f *stringFlag) String() string {
	return f.str
}

func (f *stringFlag) Set(s string) error {
	f.str = s
	return nil
}

func TestEnvFlag(t *testing.T) {

	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	testStr1 := fs.String("str-1", "def-1", "use-1")
	testInt2 := fs.Int64("int-2", 2, "use-2")
	testStr3 := fs.String("str-3", "", "use-3")
	testStr4 := fs.String("str-4", "def-4", "use-4")
	testBol5 := fs.Bool("bol-5", false, "use-5")
	testBol6 := &stringFlag{}
	fs.Var(testBol6, "str-6", "use-6")

	env := []string{
		"TEST_STR_1=env-1", // "-str-1"
		"TEST_STR_4=env-4", // "-str-4"
		"TEST_BOL_5=true",  // "-bol-5"
		"TEST_STR_6=env-6", // "-str-6"
	}

	err := PreParseFlagSet(fs, "TEST_", env, []string{"str-5"})
	if err != nil {
		t.Fatal(err)
	}

	args := []string{"-str-4", "arg-4"}

	err = fs.Parse(args)
	if err != nil {
		t.Fatal(err)
	}
	// fs.Usage()

	// no args, has envvar, has default: should use envvar (envvar overrides default, but not args)
	if have, want := *testStr1, "env-1"; have != want {
		t.Errorf("have: %q, want: %q", have, want)
	}

	// no args, no envvar, has default: should use default
	if have, want := *testInt2, int64(2); have != want {
		t.Errorf("have: %q, want: %q", have, want)
	}

	// no args, no envvar, no default: should be empty
	if have, want := *testStr3, ""; have != want {
		t.Errorf("have: %q, want: %q", have, want)
	}

	// has args, has envvar, has default: should be args (args overrides all)
	if have, want := *testStr4, "arg-4"; have != want {
		t.Errorf("have: %q, want: %q", have, want)
	}

	// no args, has envvar, no default: should be use envvar (bool test)
	if have, want := *testBol5, true; have != want {
		t.Errorf("have: %t, want: %t", have, want)
	}

	// no args, has envvar, no default: should be use envvar (var test)
	if have, want := testBol6.str, "env-6"; have != want {
		t.Errorf("have: %q, want: %q", have, want)
	}

}

func TestNameToEnv(t *testing.T) {
	for k, v := range map[string]string{
		"name":      "NAME",
		"with-dash": "WITH_DASH",
		"o":         "O",
		"one1":      "ONE1",
		".":         "_",
	} {
		if have, want := nameToEnv(k), v; have != want {
			t.Errorf("for: %q, want: %q, have: %q", k, want, have)
		}
	}
}
