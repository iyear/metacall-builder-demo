package cmd

import (
	"context"
	"fmt"
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
		Use:   "builder",
		Short: "builder is a tool for building MetaCall images",
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if _, ok := builder.Languages[arg]; !ok {
					return fmt.Errorf("invalid language %s", arg)
				}
			}
			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// set base image
			cmd.SetContext(context.WithValue(cmd.Context(), baseKey{}, builder.Base(baseImg, basePost)))
			// parse args to languages
			languages := make([]builder.Language, 0, len(args))
			for _, arg := range args {
				languages = append(languages, builder.Languages[arg])
			}
			cmd.SetContext(context.WithValue(cmd.Context(), languagesKey{}, languages))
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

	cmd.AddCommand(NewDepsCmd(), NewDevCmd())

	cmd.PersistentFlags().StringVarP(&baseImg, "base-image", "b", "debian:bullseye-slim", "base image")
	cmd.PersistentFlags().StringVar(&basePost, "base-post", "", "post command to run on base image (e.g. apt-get update)")

	return cmd
}
