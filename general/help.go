package general

import "github.com/bwmarrin/discordgo"

var HelpCommand = &discordgo.ApplicationCommand{
	Name:        "help",
	Description: "Display help and other useful information about the bot",
}

func Help(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := discordgo.MessageEmbed{
		Title: "Help",
		Description: `This bot simulates the drops one might get from doing an action a certain amount of times.
			Depending on the action, you can have multiple options like adjusting the luck, adjusting reputation, etc.
			Below are some of the Frequently Asked Questions.`,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "How does it work?",
				Value: `Every action has its own algorithm (although some are the same) containing some (pseudo)randomness to randomly determine a drop.
					Since Jagex doesn't release all the specific drop algorithms, most of them here are made with logical thinking.`,
			},
			{
				Name: "How representative is it?",
				Value: `All the information is taken from the [wiki](rs.wiki). This, together with the fact that especially the common / uncommon drop rates are not precisely known, makes this emulator not 100% representative.
					For example: In the drop tables we often just see **common** or **uncommon**. In this example we take the ratio of common to uncommon to be 2 / 3, which might differ from the actual game.`,
			},
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})
}
