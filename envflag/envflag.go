package envflag

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// PreParseFlagSet prepares fs for use with environment variables.
// A prefix is prepended to the environment variable name (e.g. "PROG_").
// The environment from environ is parsed and any matching flags have their values set.
// Flag names in excluded are skipped over.
// Environment variables names are generated from the flag names.
// The usage text of each flag in fs is modified to include these
// generated names and whether they are currently set.
// Assumed to only be run once per fs.
func PreParseFlagSet(fs *flag.FlagSet, prefix string, environ []string, excluded []string) error {
	env := make(map[string]string)
	for _, e := range environ {
		es := strings.SplitN(e, "=", 2)
		if len(es) < 2 {
			continue
		}
		env[es[0]] = es[1]
	}

	flagsToSet := make(map[string]string)

	fs.VisitAll(func(f *flag.Flag) {
		for _, e := range excluded {
			if f.Name == e {
				// skip excluded flag names
				return
			}
		}

		envName := nameToEnv(prefix + f.Name)

		var envIsSet string
		if value, ok := env[envName]; ok {
			// this flag has a currently set envvar
			flagsToSet[f.Name] = value

			// indicate a set envvar to user
			// don't show actual value; may have sensitive flags
			envIsSet = " is set"
		}

		f.Usage += fmt.Sprintf(" [%s%s]", envName, envIsSet)
	})

	for name, value := range flagsToSet {
		if err := fs.Set(name, value); err != nil {
			return fmt.Errorf("setting flag -%s from environment: %w", name, err)
		}
	}

	return nil
}

// ParseFlagSet combines a [PreParseFlagSet] and flagset fs Parse.
// Errors with pre-parse are handled similiarly to flagset fs Parse.
// I.e. continuing, exiting, or panicing.
// Pre-parse uses fs, prefix, environ, and excluded whereas fs Parse is given args.
// See [PreParseFlagSet] for parameter descriptions.
func ParseFlagSet(fs *flag.FlagSet, args []string, prefix string, environ []string, excluded []string) error {
	if fs.Parsed() {
		return errors.New("already parsed")
	}
	err := PreParseFlagSet(fs, prefix, environ, excluded)
	if err != nil {
		// essentially taken from flag failf/sprintf methods
		fmt.Fprintln(fs.Output(), err.Error())
		fs.Usage()
		// switch behavior taken from stdlib flagset Parse() method
		switch fs.ErrorHandling() {
		case flag.ContinueOnError:
			return err
		case flag.ExitOnError:
			if err == flag.ErrHelp {
				os.Exit(0)
			}
			os.Exit(2)
		case flag.PanicOnError:
			panic(err)
		}
	}
	return fs.Parse(args)
}

// Parse parses the command-line flags from [os.Args][1:] and [os.Environ].
// Must be called after all flags are defined and before flags are accessed by the program.
// See [PreParseFlagSet] for parameter descriptions.
func Parse(prefix string, excluded []string) {
	// Ignore errors; CommandLine is set for ExitOnError.
	ParseFlagSet(flag.CommandLine, os.Args[1:], prefix, os.Environ(), excluded)
}

// nameToEnv converts a flag name to an environment variable name.
// Conversion is uppercasing name and any non A-Z, 0-9 characters are
// replaced with underscores.
func nameToEnv(name string) string {
	var b strings.Builder
	for _, r := range name {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else if r >= 'a' && r <= 'z' {
			b.WriteRune(r - 32) // uppercase
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}
