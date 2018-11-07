package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const ppCommandLongDesc = `Postman Pat is here to try and ease your day to day work with Postman.

First use case will be in splitting and joining collection files so they can be better used via source control.`

//Execute is the main entrance point for the porgrame
func Execute() {
	rootCmd := &cobra.Command{
		Use:   "postmanpat",
		Short: "postmanpat : Postman helper",
		Long:  ppCommandLongDesc,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	// 'split' subcommand.
	splitCommand := &cobra.Command{
		Use:   "split [/path/to/postmancollection.json]",
		Short: "Splits an existing collection into request files.",
		Long: `Split converts an existing postman collection file into the folder structure.

		The Collection Name is used to name the folder.

		The file name for each request is Request:<Request Name>:<HTTP Method>
	`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("you must provide a collection file path")
			}
			return nil
		},
		Run: splitCmd,
	}
	rootCmd.AddCommand(splitCommand)

	// 'join' subcommand.
	joinCommand := &cobra.Command{
		Use:   "join [/path/to/postmanpatfolder]",
		Short: "Joins a split collection from folder into executable collection.",
		Long: `Join converts an split postman collection folder back into a postman collection file.
	`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("you must provide a collection folder")
			}
			return nil
		},
		Run: joinCmd,
	}
	rootCmd.AddCommand(joinCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
