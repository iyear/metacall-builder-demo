package staging

import (
	"github.com/iyear/metacall-builder-demo/pkg/builder"
	"github.com/moby/buildkit/client/llb"
)

type DepsOptions struct {
	Languages []builder.Language
	Branch    string // git branch

}

type deps struct{}

var Deps = deps{}

func (deps) Base(base llb.State, branch string) llb.State {
	return builder.Environment(base).Base(branch).Root()
}

func (deps) Languages(base llb.State, languages []builder.Language) llb.State {
	return builder.Environment(base).Languages(languages).Root()
}
