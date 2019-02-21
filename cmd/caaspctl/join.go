package main

import (
	"log"

	"github.com/spf13/cobra"
	"suse.com/caaspctl/internal/pkg/caaspctl/deployments/salt"
	"suse.com/caaspctl/pkg/caaspctl/actions/join"
)

type JoinOptions struct {
	Role string
}

func newJoinCmd() *cobra.Command {
	joinOptions := JoinOptions{}

	cmd := &cobra.Command{
		Use:   "join",
		Short: "joins a new node to the cluster",
		Run: func(cmd *cobra.Command, targets []string) {
			user, err := cmd.Flags().GetString("user")
			if err != nil {
				log.Fatalf("Unable to parse user flag: %v", err)
			}
			sudo, err := cmd.Flags().GetBool("sudo")
			if err != nil {
				log.Fatalf("Unable to parse sudo flag: %v", err)
			}

			target := salt.Target{
				Node: targets[0],
				User: user,
				Sudo: sudo,
			}

			var role join.Role
			switch joinOptions.Role {
			case "master":
				role = join.MasterRole
			case "worker":
				role = join.WorkerRole
			default:
				log.Fatalf("Invalid role provided: %q, 'master' or 'worker' are the only accepted roles", joinOptions.Role)
			}

			join.Join(target, role)
		},
		Args: cobra.ExactArgs(1),
	}

	cmd.Flags().StringP("user", "u", "root", "user identity used to connect to target")
	cmd.Flags().Bool("sudo", false, "run remote command via sudo")
	cmd.Flags().StringVarP(&joinOptions.Role, "role", "", "", "Role that this node will have in the cluster (master|worker)")

	cmd.MarkFlagRequired("role")

	return cmd
}