package simulations

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type DropFunc func(int64, dropTables, *discordgo.InteractionCreate) map[string]*Drop

// Use this function to do a lot of the discord handling and message sending. Pass the drop function and the droptables needed.
func simulateDrop(s *discordgo.Session, i *discordgo.InteractionCreate, dropFunc DropFunc, tables dropTables, url string, name string, prettyName string) {
	var amount int64 = 1
	if amountOpt := getOptionWithName(i.ApplicationCommandData().Options, "amount"); amountOpt.Name != "" {
		amount = amountOpt.IntValue()
	}

	// Replace this with max value later
	if amount > MAX_AMOUNT_ROLLS || amount < 1 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The amount has to be between 1 and %v", MAX_AMOUNT_ROLLS),
			},
		})
		return
	}

	drops := dropFunc(amount, tables, i)

	dropWithPrice, total, ok := amountToPrice(drops)
	sortDrops(&dropWithPrice)
	content := makeDropList(dropWithPrice, drops, total, ok)

	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: url,
		},
		Title:       fmt.Sprintf("You killed %v %v times and you got:", prettyName, amount),
		Description: content,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})
	log.Info("command executed", "command", name)
}
