package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/virepri/Spinner/botFE"
	"github.com/virepri/Spinner/cliFE"
	"github.com/virepri/Spinner/common"
)

var rootCmd = &cobra.Command{
	Use:   "spinner",
	Short: "Social-game discord bot to randomly connect users",
	Long: `Social-game discord bot to randomly connect users
Built with love by Virepri (github.com/Virepri/Spinner)`,
	Run: func(cmd *cobra.Command, args []string) {
		lcm := common.GetLifecycleManager()

		if err := sf.Cook(args, lcm); err != nil {
			lcm.Log(fmt.Sprintf("%s\n %s", common.EExitCode.FailedVerify(), err), common.ELogLevel.Fatal())
			os.Exit(common.EExitCode.FailedVerify().Code)
		}

		lcm.Log("Spinner "+common.Version, common.ELogLevel.Information())
		lcm.CreateRoutine(cliFE.RunCLI)
		lcm.CreateRoutine(botFE.BakeRunBot(botFE.BotParameters{
			Token: sf.OAuthToken,
		}))

		x, _ := lcm.SurrenderControl()
		if x != nil {
			os.Exit(x.Code)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&sf.OAuthToken, "token", "", "Supply your bot's OAuth token")
	rootCmd.PersistentFlags().StringVar(&sf.LogLevel, "log-level", "", "Set the minimum logging level.")
	rootCmd.PersistentFlags().StringVar(&sf.LogLocation, "log-location", "", "Set where logs are dumped. Appends to individual files, and creates new files in directories.")
	rootCmd.PersistentFlags().BoolVar(&sf.LogfileOnly, "logfile-only", false, "Log ONLY to the specified file.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
