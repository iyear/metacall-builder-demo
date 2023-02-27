package staging

import (
	"github.com/iyear/metacall-builder-demo/pkg/builder"
	"github.com/moby/buildkit/client/llb"
)

type deps struct{}

var Deps = deps{}

func (deps) Base(base llb.State, branch string) llb.State {
	return builder.Environment(base).
		Base().MetaCall(branch).Root()
}

func (deps) Languages(base llb.State, languages []string) llb.State {
	return builder.Environment(base).Languages(languages).Root()
}
