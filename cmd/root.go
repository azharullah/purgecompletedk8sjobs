package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	purgeCompletedK8sJobs "github.com/azharullah/purge-completed-k8s-jobs/pkg/purge-completed-jobs"
)

var cmdDescription = "This utility helps purge all kubernetes jobs that have completed their execution " +
	"more than x hours ago. Optionally, this also write the job spec and events to a provided log file path"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "purge-completed-jobs",
	Short: cmdDescription,
	Long:  cmdDescription,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Starting executions")

		namespace, err := cmd.Flags().GetString("namespace")
		beforeHours, err := cmd.Flags().GetInt16("before-hours")
		options := map[string]string{}

		eventsLogFile, err := cmd.Flags().GetString("events-log-file")
		if eventsLogFile != "" {
			options["eventsLogFile"] = eventsLogFile
		}

		specLogFile, err := cmd.Flags().GetString("events-log-file")
		if specLogFile != "" {
			options["specLogFile"] = specLogFile
		}

		if err != nil {
			panic(err.Error())
		}

		resp := purgeCompletedK8sJobs.PurgeJobs(namespace, beforeHours, options)
		if resp.Success {
			log.Print(resp.Msg)
		} else {
			log.Printf("Failed to delete some / all the compeleted job(s), error: %v", resp.Err.Error())
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Define the cmd flags
	rootCmd.Flags().StringP("namespace", "n", "default", "Namespace in which the operations are to be performed")
	rootCmd.Flags().Int16P("before-hours", "b", 1, "Query and delete jobs that were complete before x hours")
	rootCmd.Flags().StringP("events-log-file", "e", "", "Log file to write the job events to")
	rootCmd.Flags().StringP("spec-log-file", "s", "", "Log file to write the job spec to")

	rootCmd.MarkFlagFilename("events-log-file")
	rootCmd.MarkFlagFilename("spec-log-file")
}
