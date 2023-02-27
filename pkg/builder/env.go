package builder

import (
	_ "embed"
	"github.com/moby/buildkit/client/llb"
	"runtime"
)

/*
// languages
-[x] python
-[ ] ruby
-[ ] rust
-[ ] netcore
-[ ] netcore2
-[ ] netcore5
-[ ] netcore7
-[ ] rapidjson
-[ ] funchook
-[ ] v8/v8rep54
-[ ] v8rep57
-[ ] v8rep58
-[ ] v8rep52
-[ ] v8rep51
-[ ] nodejs
-[ ] typescript
-[ ] file
-[ ] rpc
-[ ] wasm
-[ ] java
-[ ] c
-[ ] cobol

-[ ] cache // maybe we don't need this
-[x] base
-[x] swig
-[x] metacall
-[x] pack
-[x] coverage
-[ ] clangformat
-[x] backtrace
*/

var languageEnvs = map[string]func(llb.State) llb.State{
	"python": pythonEnv,
	// TODO(iyear): Add more languages
}

type Env struct {
	state llb.State
}

func Environment(base llb.State) Env {
	return Env{
		state: base,
	}
}

func (e Env) Base() Env {
	e.state = e.state.Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install -y --no-install-recommends build-essential git cmake libgtest-dev wget apt-utils apt-transport-https gnupg dirmngr ca-certificates")).
		Root()

	return e
}

func (e Env) MetaCallClone(branch string) Env {
	e.state = e.state.Run(llb.Shlexf("git clone --depth 1 --single-branch --branch=%v https://github.com/metacall/core.git", branch)).
		Run(llb.Shlex("mkdir core/build")).Root()

	return e
}

func (e Env) MetaCallCompile() Env {
	e.state = e.state.Dir("core/build").
		Run(llb.Shlex("cmake -DOPTION_BUILD_SCRIPTS=OFF -DOPTION_BUILD_EXAMPLES=OFF -DOPTION_BUILD_TESTS=OFF -DOPTION_BUILD_DOCS=OFF -DOPTION_FORK_SAFE=OFF ..")).
		Run(llb.Shlexf("cmake --build . -j %v --target install", runtime.NumCPU())).Root()

	return e
}

func (e Env) Pack() Env {
	e.state = e.state.Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install -y --no-install-recommends rpm")).Root()

	return e
}

func (e Env) Coverage() Env {
	e.state = e.state.Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install -y --no-install-recommends lcov")).Root()

	return e
}

func (e Env) ClangFormat() Env {
	e.state = e.state.Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install -y --no-install-recommends clang-format")).Root()

	return e
}

func (e Env) Swig() Env {
	e.state = e.state.Run(llb.Shlex("apt-get install -y --no-install-recommends g++ libpcre3-dev tar")).
		Run(llb.Shlex("wget http://prdownloads.sourceforge.net/swig/swig-4.0.1.tar.gz")).
		Run(llb.Shlex("tar -xzf swig-4.0.1.tar.gz")).
		Dir("swig-4.0.1").
		Run(llb.Shlex("./configure --prefix=/usr/local")).
		Run(llb.Shlex("make")).
		Run(llb.Shlex("make install")).
		Dir("..").
		Run(llb.Shlex("rm -rf swig-4.0.1")).
		Run(llb.Shlex("pip3 install setuptools")).Root()

	return e
}

func (e Env) Backtrace() Env {
	e.state = e.state.Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install -y --no-install-recommends libdw-dev")).
		Root()

	return e
}

func (e Env) Languages(languages []string) Env {
	for _, lang := range languages {
		if env, ok := languageEnvs[lang]; ok {
			e.state = env(e.state)
		}
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
