package app

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/tomhjx/cnet/cmd/cnet/app/options"
	"github.com/tomhjx/cnet/cmd/cnet/app/version"
	"github.com/tomhjx/cnet/pkg/flow"
)

func NewCommand() *cobra.Command {

	o, err := options.NewOptions()
	if err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use: "cnet",
		// Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.Config.Version {
				version.Print()
				return nil
			}
			o.Config.OnChanged = func(c *options.Config) {
				NotifyLoadJob(context.Background(), c)
			}
			return Run(context.Background(), o.Config, flow.SetupSignalHandler())
		}}

	cmd.SuggestionsMinimumDistance = 1

	fs := cmd.Flags()
	// fs := cmd.PersistentFlags()
	o.Config.AddFlags(fs)

	return cmd
}
