package cmd

import (
	"context"
	"github.com/iyear/metacall-builder-demo/pkg/builder"
	"github.com/iyear/metacall-builder-demo/pkg/staging"
	"github.com/moby/buildkit/client/llb"
	"github.com/spf13/cobra"
)

func NewDevCmd() *cobra.Command {
	opts := staging.DevOptions{}

	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Generate development images for languages",
		Run: func(cmd *cobra.Command, args []string) {
			base := cmd.Context().Value(baseKey{}).(llb.State)

			opts.Languages = cmd.Context().Value(languagesKey{}).([]builder.Language)

			depsBase := staging.DepsBase(base, opts.Branch)
			deps := staging.DepsLang(depsBase, opts.Languages)
			dev := staging.DevConfigure(deps, opts)
			dev = staging.DevBuild(dev, opts)

			// remove depsBase from final image to reduce size
			dev = llb.Diff(depsBase, dev)

			// set final state
			cmd.SetContext(context.WithValue(cmd.Context(), finalKey{}, dev))
		},
	}

	// TODO(iyear): modify default values
	cmd.Flags().StringVarP(&opts.Type, "type", "t", "Release", "Debug/Release/RelWithDebInfo")
	cmd.Flags().BoolVar(&opts.Scripts, "scripts", false, "build all scripts")
	cmd.Flags().BoolVar(&opts.Tests, "tests", false, "build all tests")
	cmd.Flags().BoolVar(&opts.Examples, "examples", false, "build all examples")
	cmd.Flags().BoolVar(&opts.Benchmarks, "benchmarks", false, "build all benchmarks")
	cmd.Flags().BoolVar(&opts.Ports, "ports", false, "build all ports")
	cmd.Flags().BoolVar(&opts.Coverage, "coverage", false, "build all coverage")
	cmd.Flags().BoolVar(&opts.Sanitizer, "sanitizer", false, "build all sanitizer")
	cmd.Flags().BoolVar(&opts.ThreadSanitizer, "thread-sanitizer", false, "build all thread-sanitizer")
	cmd.Flags().BoolVar(&opts.Root, "root", false, "build as root")
	cmd.Flags().BoolVar(&opts.Install, "install", false, "install all dependencies")

	// deps options
	cmd.Flags().StringVarP(&opts.Branch, "branch", "B", "develop", "metacall/core branch to use")
	return cmd
}
