package envflag_test

import (
	"flag"
	"fmt"

	"github.com/micromdm/nanolib/envflag"
)

func ExampleParse() {
	var (
		flExampleString = flag.String("example", "default", "usage for example")
		flNoEnv         = flag.Bool("no-env", false, "usage for no-env")
	)

	// if an ENVFLAG_EXAMPLE environment variable exists then
	// *flExampleString will contain it after Parse assuming no
	// command-line flags are provided (which override environment
	// variables).

	// *flNoEnv is excluded from environment variable handling here.

	envflag.Parse("ENVFLAG_", []string{"no-env"})

	// as well when usage is printed it should print something along
	// these lines:
	//
	// Usage of example:
	//   -example string
	//     	usage for example [ENVFLAG_EXAMPLE] (default "default")
	//   -no-env
	//     	usage for no-env

	fmt.Printf("Example: %s\nNo-Env: %t\n", *flExampleString, *flNoEnv)
}
