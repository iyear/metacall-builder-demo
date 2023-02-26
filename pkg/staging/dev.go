package staging

import (
	"fmt"
	"github.com/iyear/metacall-builder-demo/pkg/builder"
	"github.com/moby/buildkit/client/llb"
	"runtime"
	"strings"
)

// DevConfigure generates cmake options for the given languages
//
// Ref: https://github.com/metacall/core/blob/e6253938574a15735a734d101bc1284994f62f73/tools/metacall-configure.sh
func DevConfigure(base llb.State, opts DevOptions) llb.State {
	options := []string{
		"-DOPTION_BUILD_LOG_PRETTY=Off",
		"-DOPTION_BUILD_LOADERS=On",
		"-DOPTION_BUILD_LOADERS_MOCK=On",
	}

	for _, lang := range opts.Languages {
		options = append(options, builder.Configurations[lang](opts.Scripts, opts.Ports)...)
	}

	if opts.Scripts {
		options = append(options, "-DOPTION_BUILD_SCRIPTS=On")
	} else {
		options = append(options, "-DOPTION_BUILD_SCRIPTS=Off")
	}

	if opts.Examples {
		options = append(options, "-DOPTION_BUILD_EXAMPLES=On")
	} else {
		options = append(options, "-DOPTION_BUILD_EXAMPLES=Off")
	}

	if opts.Tests {
		options = append(options, "-DOPTION_BUILD_TESTS=On")
	} else {
		options = append(options, "-DOPTION_BUILD_TESTS=Off")
	}

	if opts.Benchmarks {
		options = append(options, "-DOPTION_BUILD_BENCHMARKS=On")
	} else {
		options = append(options, "-DOPTION_BUILD_BENCHMARKS=Off")
	}

	if opts.Ports {
		options = append(options, "-DOPTION_BUILD_PORTS=On")
	} else {
		options = append(options, "-DOPTION_BUILD_PORTS=Off")
	}

	if opts.Coverage {
		options = append(options, "-DOPTION_BUILD_COVERAGE=On")
	} else {
		options = append(options, "-DOPTION_BUILD_COVERAGE=Off")
	}

	if opts.Sanitizer {
		options = append(options, "-DOPTION_BUILD_SANITIZER=On")
	} else {
		options = append(options, "-DOPTION_BUILD_SANITIZER=Off")
	}

	if opts.ThreadSanitizer {
		options = append(options, "-DOPTION_BUILD_THREAD_SANITIZER=On")
	} else {
		options = append(options, "-DOPTION_BUILD_THREAD_SANITIZER=Off")
	}

	options = append(options, "-DCMAKE_BUILD_TYPE="+opts.Type)

	return base.
		Run(llb.Shlexf(`cmake -Wno-dev -DOPTION_GIT_HOOKS=Off %s ..`, strings.Join(options, " "))).
		Root()
}

type DevOptions struct {
	Type            string // Debug/Release/RelWithDebInfo
	Languages       []builder.Language
	Scripts         bool
	Examples        bool
	Benchmarks      bool
	Ports           bool
	Sanitizer       bool
	ThreadSanitizer bool
	Static          bool // not used
	Dynamic         bool // not used
	Root            bool // TODO(iyear): do we need to run as root?
	Tests           bool
	Coverage        bool
	Install         bool

	// depend on deps
	DepsOptions
}

// DevBuild TODO(iyear): split into multiple stages
func DevBuild(base llb.State, opts DevOptions) llb.State {
	commands := make([]string, 0)

	commands = append(commands, fmt.Sprintf("make -k -j%v", runtime.NumCPU()))

	// coverage needs to run the tests
	if opts.Tests || opts.Coverage {
		commands = append(commands, fmt.Sprintf("ctest -j%v --output-on-failure --test-output-size-failed 3221000000 -C %s",
			runtime.NumCPU()+1, opts.Type))
	}

	if opts.Coverage {
		commands = append(commands,
			"make -k gcov",
			"make -k lcov",
			"make -k lcov-genhtml")
	}

	if opts.Install {
		// TODO(iyear): do we need to run as root?
		if opts.Root {
			commands = append(commands, "make install")
		} else {
			// Needed for rustup in order to install rust loader properly
			commands = append(commands, `sudo HOME="$HOME" make install`)
		}
	}

	for _, command := range commands {
		base = base.Run(llb.Shlex(command)).Root()
	}

	return base
}
