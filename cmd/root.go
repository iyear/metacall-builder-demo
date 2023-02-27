package cmd

import (
	"context"
	"github.com/iyear/metacall-builder-demo/pkg/builder"
	"github.com/moby/buildkit/client/llb"
	"github.com/spf13/cobra"
	"os"
)

func NewRootCmd() *cobra.Command {
	var (
		baseImg  string
		basePost string
	)

	cmd := &cobra.Command{
		Use:           "builder",
		Short:         "builder is a tool for building MetaCall images",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// set base image
			cmd.SetContext(context.WithValue(cmd.Context(), baseKey{}, builder.Base(baseImg, basePost)))
			// set languages
			cmd.SetContext(context.WithValue(cmd.Context(), languagesKey{}, args))
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			state := cmd.Context().Value(finalKey{}).(llb.State)

			def, err := state.Marshal(cmd.Context(), llb.LinuxAmd64)
			if err != nil {
				return err
			}

			return llb.WriteTo(def, os.Stdout)
		},
	}

	cmd.AddCommand(NewDepsCmd())

	cmd.PersistentFlags().StringVarP(&baseImg, "base-image", "b", "debian:bullseye-slim", "base image")
	cmd.PersistentFlags().StringVar(&basePost, "base-post", "", "post command to run on base image (e.g. apt-get update)")

	return cmd
}
