package cmd

import (
	"context"
	"github.com/iyear/metacall-builder-demo/pkg/builder"
	"github.com/iyear/metacall-builder-demo/pkg/staging"
	"github.com/moby/buildkit/client/llb"
	"github.com/spf13/cobra"
)

func NewDepsCmd() *cobra.Command {
	opts := staging.DepsOptions{}

	cmd := &cobra.Command{
		Use:   "deps",
		Short: "Generate dependencies images for languages",
		Run: func(cmd *cobra.Command, args []string) {
			base := cmd.Context().Value(baseKey{}).(llb.State)

			opts.Languages = cmd.Context().Value(languagesKey{}).([]builder.Language)

			depsBase := staging.Deps.Base(base, opts.Branch)
			deps := staging.Deps.Languages(depsBase, opts.Languages)

			// remove deps base from final image to reduce size
			deps = llb.Diff(depsBase, deps)

			// set final state
			cmd.SetContext(context.WithValue(cmd.Context(), finalKey{}, deps))
		},
	}

	cmd.Flags().StringVarP(&opts.Branch, "branch", "B", "develop", "core git branch to use")

	return cmd
}
