package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stevegood/flare/bot"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var discordToken string
var rootCmd = &cobra.Command{
	Use:   "flare",
	Short: "Flare is a discord bot for the Solar gaming group in San Diego",
	Long: `Take control of channel moderation while also adding fun tools
			like Lumpy gifs and Gonk info.`,
	Run: func(cmd *cobra.Command, args []string) {
		discordBot := bot.NewDiscord(discordToken)
		err := discordBot.Connect()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	discordToken = os.Getenv("DISCORD_TOKEN")
	if discordToken == "" {
		viper.SetConfigType("yaml")

		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.flare")
		viper.AddConfigPath("/opt/flare")

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Println("Could not find config.yml")
			}
			log.Fatal(err)
		}

		discordToken = viper.GetString("discord.token")
	}
}

// Execute is part of the Cobra interface
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
