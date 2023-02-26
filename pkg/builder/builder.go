package builder

import (
	_ "embed"
	"github.com/moby/buildkit/client/llb"
)

func Base(image string, post string) llb.State {
	state := llb.Image(image)
	if post != "" {
		state = state.Run(llb.Shlex(post)).Root()
	}
	return state
}
