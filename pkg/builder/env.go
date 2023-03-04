package builder

import (
	_ "embed"
	"fmt"
	"github.com/moby/buildkit/client/llb"
	"runtime"
)

/*
// languages
-[x] python
-[x] ruby
-[x] rust
-[x] netcore
-[x] netcore2
-[x] netcore5
-[x] netcore7
-[x] rapidjson
-[x] funchook
-[ ] v8/v8rep54
-[ ] v8rep57
-[ ] v8rep58
-[ ] v8rep52
-[ ] v8rep51
-[x] nodejs
-[x] typescript
-[x] file
-[x] rpc
-[x] wasm
-[x] java
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

func (e Env) Python() Env {
	e.state = e.state.Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get -y --no-install-recommends install python3 python3-dev python3-pip")).
		Root()

	return e
}

func (e Env) Ruby() Env {
	e.state = e.state.Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install -y --no-install-recommends ruby ruby-dev")).
		Root()

	return e
}

func (e Env) Rust() Env {
	e.state = e.state.Run(llb.Shlex("apt-get install -y --no-install-recommends curl")).
		Run(llb.Shlex("curl https://sh.rustup.rs -sSf | sh -s -- -y --default-toolchain nightly-2021-12-04 --profile default")).
		Root()

	return e
}

func (e Env) RapidJSON() Env {
	e.state = e.state.Run(llb.Shlex("git clone https://github.com/miloyip/rapidjson.git")).
		Dir("rapidjson").
		Run(llb.Shlex("git checkout v1.1.0")).
		Run(llb.Shlex("mkdir build")).
		Dir("build").
		Run(llb.Shlex("cmake -DRAPIDJSON_BUILD_DOC=Off -DRAPIDJSON_BUILD_EXAMPLES=Off -DRAPIDJSON_BUILD_TESTS=Off ..")).
		Run(llb.Shlex("make && make install")).
		Dir("../..").
		Run(llb.Shlex("rm -rf ./rapidjson")).
		Root()

	return e
}

func (e Env) FuncHook() Env {
	e.state = e.state.Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install -y --no-install-recommends cmake")).
		Root()

	return e
}

func (e Env) NetCore() Env {
	version := "1.1.11"
	e.state = e.state.
		Run(llb.Shlex("apt-get install -y --no-install-recommends libc6 libcurl3 libgcc1 libgssapi-krb5-2 libicu57 liblttng-ust0 libssl1.0.2 libstdc++6 libunwind8 libuuid1 zlib1g")).
		Run(llb.Shlexf("wget %s -O dotnet.tar.gz", fmt.Sprintf("https://dotnetcli.blob.core.windows.net/dotnet/Sdk/%s/dotnet-dev-debian.9-x64.%s.tar.gz", version, version))).
		Run(llb.Shlex("mkdir -p /usr/share/dotnet")).
		Run(llb.Shlex("tar -zxf dotnet.tar.gz -C /usr/share/dotnet")).
		Run(llb.Shlex("rm dotnet.tar.gz")).
		Run(llb.Shlex("ln -s /usr/share/dotnet/dotnet /usr/bin/dotnet")).
		Run(llb.Shlex("mkdir warmup")).
		Dir("warmup").
		Run(llb.Shlex("dotnet new")).
		Dir("..").
		Run(llb.Shlex("rm -rf warmup")).
		Run(llb.Shlex("rm -rf /tmp/NuGetScratch")).
		Root()

	return e
}

func (e Env) netCoreBase() Env {
	e.state = e.state.
		Run(llb.Shlex("wget https://packages.microsoft.com/config/debian/10/packages-microsoft-prod.deb -O packages-microsoft-prod.deb")).
		Run(llb.Shlex("dpkg -i packages-microsoft-prod.deb")).
		Run(llb.Shlex("rm packages-microsoft-prod.deb")).
		Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install -y --no-install-recommends apt-transport-https")).
		Root()

	return e
}

func (e Env) NetCore2() Env {
	e.state = e.netCoreBase().state.
		Run(llb.Shlex("apt-get install -y --no-install-recommends dotnet-sdk-2.2")).
		Root()

	return e
}

func (e Env) NetCore5() Env {
	e.state = e.netCoreBase().state.
		Run(llb.Shlex("apt-get install -y --no-install-recommends dotnet-sdk-5.0")).
		Root()

	return e
}

func (e Env) NetCore7() Env {
	e.state = e.netCoreBase().state.
		Run(llb.Shlex("apt-get install -y --no-install-recommends dotnet-sdk-7.0")).
		Root()

	return e
}

func (e Env) NodeJS() Env {
	e.state = e.state.
		Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install -y --no-install-recommends python3 g++ make nodejs npm curl")).
		Root()

	return e
}

func (e Env) TypeScript() Env {
	e.state = e.state.
		// Install React dependencies in order to run the tests
		Run(llb.Shlex("npm i react@latest -g")).
		Run(llb.Shlex("npm i react-dom@latest -g")).
		Root()

	return e
}

func (e Env) File() Env {
	return e
}

func (e Env) RPC() Env {
	e.state = e.state.
		Run(llb.Shlex("apt-get install -y --no-install-recommends libcurl4-openssl-dev")).
		Root()

	return e
}

func (e Env) WASM() Env {
	// TODO: add shell script

	return e
}

func (e Env) Java() Env {
	e.state = e.state.
		Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install -y --no-install-recommends default-jre default-jdk")).
		Root()

	return e
}

func (e Env) Root() llb.State {
	return e.state
}
