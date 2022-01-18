package general

import "github.com/bwmarrin/discordgo"

var ContributeCommand = &discordgo.ApplicationCommand{
	Name:        "contribute",
	Description: "Display contribute information",
}

func Contribute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := discordgo.MessageEmbed{
		Title: "Contribute",
		Description: `This bot is open source with the [GPL-3.0 License](https://github.com/JasperSurmont/rs-drop-emulator/blob/main/LICENSE).
			You can contribute by heading over to the [GitHub](https://github.com/JasperSurmont/rs-drop-emulator) page.
			Below is some extra info about me, the creator.`,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Who I am",
				Value: `I am Jasper, a computer science student from Belgium. I've been playing Runescape on-off for quite some time now, with username **Pai Sho**.`,
			},
			{
				Name: "Why I made this bot",
				Value: `I actually just wanted to increase my knowledge in the **Go** programming language. However, I never knew what to make.
					Then I saw the discordgo package, and this simulator was the first thing that popped in my head. Hence, this bot was born.`,
			},
			{
				Name: "My experience",
				Value: `I've worked on many different projects (both for school and personal). These were mostly in JS/TS, Python, C(++) and Java. 
					This is my first _big_ Go project. That is also why an experienced Go developer might look at my code and start to vomit.
					However, I'm open for any feedback or contributions, as long as we keep it civil :).`,
			},
			{
				Name:  "Contact information",
				Value: `You can contact me via my email [surmontjasper@gmail.com](mailto:surmontjasper@gmail.com), or by opening an issue on GitHub.`,
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
