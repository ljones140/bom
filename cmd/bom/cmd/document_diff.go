package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sigs.k8s.io/bom/pkg/spdx"
)

func AddDiff(parent *cobra.Command) {
	diffCmd := &cobra.Command{
		PersistentPreRunE: initLogging,
		Short:             "bom document diff → Compare two SPDX documents",
		Long: `bom document diff → Compare two SPDX documents",

		this commands does some diffing
		`,
		Use:           "diff sbom.spdx.json sbom.spdx.json",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			fmt.Println("diffing")
			if len(args) != 2 {
				return fmt.Errorf("need two args")
			}
			docA, err := spdx.OpenDoc(args[0])
			if err != nil {
				return fmt.Errorf("opening docA: %w", err)
			}
			docB, err := spdx.OpenDoc(args[1])
			if err != nil {
				return fmt.Errorf("opening docB: %w", err)
			}

			outputA, err := spdx.DiffPackages(args[0], docA, docB)
			if err != nil {
				return fmt.Errorf("Error diffing A against B %v", err)
			}
			fmt.Println(outputA)

			outputB, err := spdx.DiffPackages(args[1], docB, docA)
			if err != nil {
				return fmt.Errorf("Error diffing B against A %v", err)
			}
			fmt.Println(outputB)
			return nil
		},
	}

	parent.AddCommand(diffCmd)
}
