package builder

import (
	_ "embed"
	"github.com/moby/buildkit/client/llb"
)

var Envs = map[Language]func(llb.State) llb.State{
	Python: pythonEnv,
	Node:   nodeEnv,
	// TODO(iyear): Add more languages
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
