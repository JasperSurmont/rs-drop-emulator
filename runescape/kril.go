package runescape

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	krilUrl             = "https://runescape.wiki/images/K%27ril_Tsutsaroth.png?11873"
	krilCommonDroptable = []Drop{
		{
			Name:   "Infernal ashes",
			Amount: 5,
		},
		{
			Name:   "Medium bladed rune salvage",
			Amount: 1,
		},
		{
			Name:   "Large plated rune salvage",
			Amount: 1,
		},
		{
			Name:   "Huge plated adamant salvage",
			Amount: 1,
		},
		{
			Name:   "Super attack (3)",
			Amount: 3,
			OtherDrops: []Drop{
				{
					Name:   "Super strength (3)",
					Amount: 3,
				},
			},
		},
		{
			Name:   "Super restore (3)",
			Amount: 1,
			OtherDrops: []Drop{
				{
					Name:   "Zamorak brew (3)",
					Amount: 3,
				},
			},
		},
	}
	krilUncommonDroptable = []Drop{
		{
			Name:   "Lantadyme seed",
			Amount: 3,
		},
		{
			Name:   "Grimy lantadyme",
			Amount: 10,
		},
		{
			Name:   "Orichalcite stone spirit",
			Amount: 3,
		},
		{
			Name:        "Wine of Zamorak",
			AmountRange: [2]int{2, 10},
		},
	}
	krilRareDroptable = []Drop{
		{
			Name: "Hood of subjugation",
			Rate: 1.0 / 512.0,
			Bold: true,
		},
		{
			Name: "Garb of subjugation",
			Rate: 1.0 / 512.0,
			Bold: true,
		},
		{
			Name: "Gown of subjugation",
			Rate: 1.0 / 512.0,
			Bold: true,
		},
		{
			Name: "Gloves of subjugation",
			Rate: 1.0 / 512.0,
			Bold: true,
		},
		{
			Name: "Boots of subjugation",
			Rate: 1.0 / 512.0,
			Bold: true,
		},
		{
			Name: "Ward of subjugation",
			Rate: 1.0 / 512.0,
			Bold: true,
		},
		{
			Name: "Zamorakian spear",
			Rate: 1.0 / 512.0,
			Bold: true,
		},
		{
			Name: "Steam battlestaff",
			Rate: 1.0 / 512.0,
			Bold: true,
		},
	}
	zamorakHilt = Drop{
		Name:   "Zamorak hilt",
		Amount: 1,
	}
	krilAlwaysDroptable = []Drop{
		{
			Name:   "Infernal ashes",
			Amount: 1,
		},
	}
)

var KrilCommand = &discordgo.ApplicationCommand{
	Name:        "kril",
	Description: "Emulate a General Kril drop",
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

func Kril(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var amount int64 = 1
	if amountOpt := GetOptionWithName(i.ApplicationCommandData().Options, "amount"); amountOpt.Name != "" {
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

	drops := EmulateDropGwd1(amount, krilRareDroptable, krilUncommonDroptable, krilCommonDroptable, bandosHilt)
	if enableGuarantees := GetOptionWithName(i.ApplicationCommandData().Options, "enable-guarantees"); enableGuarantees.Name == "" || enableGuarantees.BoolValue() {
		AddAlwaysDroptable(amount, &drops, krilAlwaysDroptable)
	}

	dropWithPrice, total, ok := AmountToPrice(drops)
	SortDrops(&dropWithPrice)
	content := MakeDropList(dropWithPrice, drops, total, ok)

	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: krilUrl,
		},
		Title:       fmt.Sprintf("You killed General Kril %v times and you got:", amount),
		Description: content,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})
}
