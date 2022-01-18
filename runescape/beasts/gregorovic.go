package beasts

import (
	"fmt"

	"rs-drop-emulator/runescape/core"

	"github.com/bwmarrin/discordgo"
)

var (
	gregorovicUrl             = "https://runescape.wiki/images/Gregorovic.png?d97e5"
	gregorovicCommonDroptable = []core.Drop{
		{
			Name:        "Drakolith stone spirit",
			AmountRange: [2]int{15, 25},
		},
		{
			Name:        "Orichalcite stone spirit",
			AmountRange: [2]int{15, 25},
		},
		{
			Name:        "Uncut diamond",
			AmountRange: [2]int{18, 22},
		},
		{
			Name:        "Grimy dwarf weed",
			AmountRange: [2]int{14, 25},
		},
		{
			Name:        "Raw shark",
			AmountRange: [2]int{45, 55},
		},
	}

	gregorovicUncommonDroptable = []core.Drop{
		{
			Name:        "Phasmatite stone spirit",
			AmountRange: [2]int{15, 25},
		},
		{
			Name:        "Necrite stone spirit",
			AmountRange: [2]int{15, 25},
		},
		{
			Name:        "Coins",
			AmountRange: [2]int{60000, 82816},
		},
		{
			Name:        "Uncut dragonstone",
			AmountRange: [2]int{8, 12},
		},
		{
			Name:        "Magic logs",
			AmountRange: [2]int{150, 250},
		},
		{
			Name:        "Large bladed rune salvage",
			AmountRange: [2]int{5, 10},
		},
		{
			Name:        "Medium plated rune salvage",
			AmountRange: [2]int{5, 10},
		},
		{
			Name:        "Battlestaff",
			AmountRange: [2]int{45, 60},
		},
	}

	gregorovicRareDroptable = []core.Drop{
		{
			Rate: 1.0 / 256.0,
			Name: "Dormant anima core helm",
			Bold: true,
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Dormant anima core body",
			Bold: true,
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Dormant anima core legs",
			Bold: true,
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Shadow glaive",
			Bold: true,
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Off-hand shadow glaive",
			Bold: true,
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Crest of Sliske",
			Bold: true,
		},
		{
			Rate: 1.0 / 64.0,
			Name: "Sliskean essence",
			Bold: true,
		},
	}

	gregorovicAlwaysDroptable = []core.Drop{}
)

var GregorovicCommand = &discordgo.ApplicationCommand{
	Name:        "gregorovic",
	Description: "Emulate a Gregorovic drop with full reputation",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "amount",
			Description: "Amount of kills",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "enable-guarantees",
			Description: "Set to false to remove the guaranteed drops",
			Required:    false,
		},
	},
}

func Gregorovic(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var amount int64 = 1
	if lenOpt := len(i.ApplicationCommandData().Options); lenOpt >= 1 {
		amount = i.ApplicationCommandData().Options[0].IntValue()
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

	drops := core.EmulateDropGwd2(amount, gregorovicRareDroptable, gregorovicUncommonDroptable, gregorovicCommonDroptable)
	if lenOpt := len(i.ApplicationCommandData().Options); lenOpt < 2 || i.ApplicationCommandData().Options[1].BoolValue() {
		core.AddAlwaysDroptable(amount, &drops, gregorovicAlwaysDroptable)
	}

	dropWithPrice, total, ok := core.AmountToPrice(drops)
	core.SortDrops(&dropWithPrice)
	content := core.MakeDropList(dropWithPrice, drops, total, ok)

	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: gregorovicUrl,
		},
		Title:       fmt.Sprintf("You killed Gregorovic %v times and you got:", amount),
		Description: content,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})
}
