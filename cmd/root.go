package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/user/procfile-run/procfile"
	"github.com/user/procfile-run/runner"
)

var (
	procfilePath string
	envFilePath  string
	port         int
)

var rootCmd = &cobra.Command{
	Use:   "procfile-run",
	Short: "Minimal Procfile runner for local development",
	Long:  `Run processes defined in a Procfile with colored output, process supervision, and port conflict detection.`,
	RunE:  run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&procfilePath, "procfile", "f", "Procfile", "path to the Procfile")
	rootCmd.Flags().StringVarP(&envFilePath, "env", "e", ".env", "path to the .env file")
}

func run(cmd *cobra.Command, args []string) error {
	processes, err := procfile.Parse(procfilePath)
	if err != nil {
		return fmt.Errorf("failed to parse Procfile: %w", err)
	}

	if len(processes) == 0 {
		return fmt.Errorf("no processes found in %s", procfilePath)
	}

	env, err := runner.LoadEnvFile(envFilePath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "warning: could not load env file: %v\n", err)
	}

	if err := runner.CheckPorts(processes); err != nil {
		color.Red("Port conflict detected: %v", err)
		return err
	}

	palette := runner.NewPalette()
	supervisor := runner.NewSupervisor()

	for name, command := range processes {
		logger := runner.NewLogger(name, palette.ColorFor(name))
		proc := runner.NewProcess(name, command, runner.EnvMapToSlice(env), logger)
		supervisor.Add(name, proc)
	}

	return supervisor.Start()
}
