package staging

import (
	"github.com/iyear/metacall-builder-demo/pkg/builder/env"
	"github.com/moby/buildkit/client/llb"
)

type deps struct{}

var Deps = deps{}

func (deps) Base(base llb.State, branch string) llb.State {
	return env.New(base).
		Base().MetaCallClone(branch).MetaCallCompile().Root()
}

func (deps) Languages(base llb.State, languages []string) llb.State {
	state := env.New(base)

	for _, lang := range languages {
		switch lang {
		case "python":
			state = state.Python()
		}
	}

	return state.Root()
}
