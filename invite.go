package main

import "github.com/bwmarrin/discordgo"

const INVITE_URL = "https://discord.com/oauth2/authorize?client_id=931567437437075496&permissions=2048&scope=bot%20applications.commands"

var InviteCommand = &discordgo.ApplicationCommand{
	Name:        "invite",
	Description: "Get the invite link of the bot",
}

func Invite(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: INVITE_URL,
		},
	})
}
