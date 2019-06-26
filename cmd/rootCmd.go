package cmd

import (
	"fmt"
	"os"
	
	"github.com/spf13/cobra"
	
	"github.com/virepri/Spinner/common"
)

var rootCmd = &cobra.Command{
	Use: "spinner",
	Short: "Social-game discord bot to randomly connect users",
	Long: `Social-game discord bot to randomly connect users
Built with love by Virepri (github.com/Virepri/Spinner)`,
	Run: func(cmd *cobra.Command, args []string) {
		lcm := common.GetLifecycleManager()
		
		if err := sf.Verify(args); err != nil {
			fmt.Println(common.EExitCode.FailedVerify())
			fmt.Println(err)
			os.Exit(common.EExitCode.FailedVerify().Code)
		}
		
		// TODO: Start bot & CLI frontend routines
		
		lcm.SurrenderControl()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&sf.OAuthToken, "token", "", "Supply your bot's OAuth token")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
