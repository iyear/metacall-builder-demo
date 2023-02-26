package builder

import "github.com/moby/buildkit/client/llb"

var Runtimes = map[Language]func(llb.State) llb.State{
	Python: pythonRuntime,
	Node:   nodeRuntime,
}

func pythonRuntime(base llb.State) llb.State {
	return base.
		Run(llb.Shlex("apt-get -y install --no-install-recommends libpython3.9")).
		Run(llb.Shlex("apt-mark hold libpython3.9")).Root()
}

func nodeRuntime(base llb.State) llb.State {
	return base.
		Run(llb.Shlex("apt-get -y --no-install-recommends install libnode72")).
		Run(llb.Shlex("apt-mark hold libnode72")).Root()
}
