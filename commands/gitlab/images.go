package gitlab

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dryRun bool

func NewCmdImages() *cobra.Command {
	var imagesCmd = &cobra.Command{
		Use:   "images",
		Short: "Perform prune images operation on GitLab Registry (registry.gitlab.com)",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Gl images")
		},
	}
	imagesCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run results")
	return imagesCmd
}
