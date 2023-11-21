package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/working-group-two/number-lookup-provider-tester/internal/server"
	"os"
)

const (
	flagPort           = "port"
	flagRps            = "rps"
	flagPrintRequests  = "print-requests"
	flagPrintResponses = "print-responses"
	flagPrintProgress  = "print-progress"
	flagPhoneNumber    = "phone-number"
)

var cmd = &cobra.Command{
	Use:   "number-lookup",
	Short: "Number Lookup Provider Tester",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt(flagPort)
		rps, _ := cmd.Flags().GetUint32(flagRps)
		phoneNumber, _ := cmd.Flags().GetStringSlice(flagPhoneNumber)
		requireFlag(cmd, flagPort)
		requireFlag(cmd, flagRps)
		requireFlag(cmd, flagPhoneNumber)

		printRequests, _ := cmd.Flags().GetBool(flagPrintRequests)
		printResponses, _ := cmd.Flags().GetBool(flagPrintResponses)
		printProgress, _ := cmd.Flags().GetBool(flagPrintProgress)

		printOptions := &server.PrintOptions{
			Requests:  printRequests,
			Responses: printResponses,
			Progress:  printProgress,
		}

		listener := fmt.Sprintf(":%d", port)

		server.Start(
			listener,
			rps,
			phoneNumber,
			printOptions,
		)
	},
}

func init() {
	cmd.Flags().Int(flagPort, 0, "Port to run the application on")
	cmd.Flags().Uint32(flagRps, 0, "Requests per second")
	cmd.Flags().StringSlice(flagPhoneNumber, []string{}, "Phone number to use in requests")

	cmd.Flags().Bool(flagPrintRequests, false, "Print requests")
	cmd.Flags().Bool(flagPrintResponses, false, "Print responses")
	cmd.Flags().Bool(flagPrintProgress, false, "Print progress")
}

func requireFlag(cmd *cobra.Command, flagName string) {
	if cmd.Flags().Changed(flagName) == false {
		cmd.Help()
		fmt.Printf("\nflag is required: --%s\n", flagName)
		os.Exit(1)
	}
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
