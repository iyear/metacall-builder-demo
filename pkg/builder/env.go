package builder

import (
	_ "embed"
	"fmt"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/util/appcontext"
	"runtime"
)

var languageEnvs = map[Language]func(llb.State) llb.State{
	Python: pythonEnv,
	Node:   nodeEnv,
	// TODO(iyear): Add more languages
}

type Env struct {
	state llb.State
}

func Environment(base llb.State) Env {

	fmt.Println(base.GetEnv(appcontext.Context(), "APT_CACHE_DIR"))

	return Env{
		state: base,
	}

}

func (e Env) Base(branch string) Env {
	e.state = e.state.Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install -y --no-install-recommends build-essential git cmake libgtest-dev wget apt-utils apt-transport-https gnupg dirmngr ca-certificates")).
		Run(llb.Shlexf("git clone --depth 1 --single-branch --branch=%v https://github.com/metacall/core.git", branch)).
		Run(llb.Shlex("mkdir core/build")).
		Dir("core/build").
		Run(llb.Shlex("cmake -DOPTION_BUILD_SCRIPTS=OFF -DOPTION_BUILD_EXAMPLES=OFF -DOPTION_BUILD_TESTS=OFF -DOPTION_BUILD_DOCS=OFF -DOPTION_FORK_SAFE=OFF ..")).
		Run(llb.Shlexf("cmake --build . -j %v --target install", runtime.NumCPU())).Root()

	return e
}

func (e Env) Swig() Env {
	return e
}

func (e Env) Languages(languages []Language) Env {
	for _, lang := range languages {
		e.state = languageEnvs[lang](e.state)
	}
	return e
}

func (e Env) Root() llb.State {
	return e.state
}

func pythonEnv(base llb.State) llb.State {
	return base.
		Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get -y --no-install-recommends install python3 python3-dev python3-pip")).Root()
}

func nodeEnv(base llb.State) llb.State {
	return base.
		Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get -y --no-install-recommends install python3 g++ make nodejs npm")).Root()
}
