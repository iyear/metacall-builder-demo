package staging

import (
	"github.com/iyear/metacall-builder-demo/pkg/builder"
	"github.com/moby/buildkit/client/llb"
	"runtime"
)

type DepsOptions struct {
	Languages []builder.Language
	Branch    string // git branch
	// TODO(iyear): Add more options
}

func DepsBase(base llb.State, branch string) llb.State {
	return base.Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install -y --no-install-recommends build-essential git cmake libgtest-dev wget apt-utils apt-transport-https gnupg dirmngr ca-certificates")).
		Run(llb.Shlexf("git clone --depth 1 --single-branch --branch=%v https://github.com/metacall/core.git", branch)).
		Run(llb.Shlex("mkdir core/build")).
		Dir("core/build").
		Run(llb.Shlex("cmake -DOPTION_BUILD_SCRIPTS=OFF -DOPTION_BUILD_EXAMPLES=OFF -DOPTION_BUILD_TESTS=OFF -DOPTION_BUILD_DOCS=OFF -DOPTION_FORK_SAFE=OFF ..")).
		Run(llb.Shlexf("cmake --build . -j %v --target install", runtime.NumCPU())).Root()
}

func DepsLang(base llb.State, languages []builder.Language) llb.State {
	// TODO(iyear): Add more commands
	for _, lang := range languages {
		base = builder.Envs[lang](base)
	}
	return base
}

// TODO(iyear): more stages
