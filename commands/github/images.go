package github

import (
	"context"
	"fmt"
	"time"

	"github.com/ci-monk/drprune/internal/constants"
	log "github.com/ci-monk/drprune/internal/log"
	"github.com/ci-monk/drprune/internal/utils"
	gh "github.com/ci-monk/drprune/pkg/github"

	"github.com/spf13/cobra"
)

const Day = 24

var dryRun, untagged bool
var olderThan int

func NewCmdImages() *cobra.Command {
	var imagesCmd = &cobra.Command{
		Use:   "images",
		Short: "Perform prune images operation on GitHub Registry (ghcr.io)",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(constants.ASCIIPrune)
			checkCmdParams()

			// Auth in github client
			ctx := context.Background()
			client, err := gh.NewClient(ctx, token, name, "")
			if err != nil {
				log.Fatal(err)
			}

			pkgVersions, err := client.GetUserAllContainerPackageVersions(ctx, container)
			if err != nil {
				log.Fatal(err)
			}

			size := len(pkgVersions)

			log.Infof("Package: %s", utils.DecodeParam(container))
			log.Infof("We have %v versions in this package", size)

			now := time.Now().UTC()

			// Loop in the list of package versions
			for _, pkg := range pkgVersions {
				tagged := false
				old := false

				if pkg.CreatedAt != nil {
					// Get the time difference between now and the package creation time
					then := *pkg.CreatedAt
					days := int(now.Sub(then.Time).Hours() / Day)

					if days > olderThan {
						old = true
					}
				}

				tags := pkg.Metadata.Container.Tags
				// Check untagged
				if !(len(tags) == 0) {
					// Is tagged
					tagged = true
				}

				if (untagged && !tagged) || (old && olderThan > 0) {
					client.DeleteContainerPackageVersion(ctx, container, pkg, dryRun)
				}
			}
		},
	}
	imagesCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Controlling whether to execute the action as a dry-run")
	imagesCmd.PersistentFlags().BoolVarP(&untagged, "untagged", "u", false, "Boolean controlling whether untagged versions should be pruned")
	imagesCmd.PersistentFlags().IntVarP(&olderThan, "older-than", "o", 0, "Minimum age in days of a version before it is pruned")
	return imagesCmd
}
