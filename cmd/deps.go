package cmd

import (
	"context"
	"github.com/iyear/metacall-builder-demo/pkg/builder"
	"github.com/moby/buildkit/client/llb"
	"github.com/spf13/cobra"
)

func NewDepsCmd() *cobra.Command {
	type options struct {
		Languages []string
		Branch    string // git branch
	}
	opts := options{}

	cmd := &cobra.Command{
		Use:   "deps",
		Short: "Generate dependencies images for languages",
		Run: func(cmd *cobra.Command, args []string) {
			base := cmd.Context().Value(baseKey{}).(llb.State)

			opts.Languages = cmd.Context().Value(languagesKey{}).([]string)

			depsBase := builder.Environment(base).Base().MetaCallClone(opts.Branch).Root()
			deps := builder.Environment(depsBase).MetaCallCompile().Languages(opts.Languages).Root()

			// remove deps base from final image to reduce size
			deps = llb.Diff(depsBase, deps)

			// set final state
			cmd.SetContext(context.WithValue(cmd.Context(), finalKey{}, deps))
		},
	}

	cmd.Flags().StringVarP(&opts.Branch, "branch", "B", "develop", "core git branch to use")

	return cmd
}
