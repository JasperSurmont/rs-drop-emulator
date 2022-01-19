package beasts

import (
	"fmt"
	"rs-drop-emulator/runescape/core"

	"github.com/bwmarrin/discordgo"
)

//ref: https://archive.fo/0KkjH
// needs adjustment, what to do with clue scrolls? Numbing root, clinby mole and D2H add up to 100%, when is clue scroll rolled?
// unkown items are added to common
var (
	giantmoleUrl             = "https://runescape.wiki/images/Giant_Mole.png?6906f"
	giantmoleCommonDroptable = []core.Drop{
		{
			Name:   "Uncut diamond",
			Amount: 6,
		},
		{
			Name:    "Grimy irit",
			Amounts: []int{3, 6},
		},
		{
			Name:   "Limpwurt root",
			Amount: 10,
		},
	}
	giantmoleUncommonDroptable = []core.Drop{
		{
			Name:        "Coins",
			AmountRange: [2]int{13207, 23500},
		},
		{
			Name:        "Yew logs",
			AmountRange: [2]int{15, 30},
		},
		{
			Name:        "Red dragonhide",
			AmountRange: [2]int{1, 3},
		},
		{
			Name:   "Iron stone spirit",
			Amount: 40,
		},
		{
			Name:   "Coal stone spirit",
			Amount: 40,
		},
		{
			Name:        "Adamantite stone spirit",
			AmountRange: [2]int{2, 4},
		},
		{
			Name:    "Pure essence",
			Amounts: []int{300, 600},
		},
		{
			Name:        "Nature rune",
			AmountRange: [2]int{40, 60},
		},
		{
			Name:   "Blood rune",
			Amount: 60,
		},
		{
			Name:    "Grimy ranarr",
			Amounts: []int{3, 6},
		},
	}
	giantmoleRareDroptable = []core.Drop{
		{
			Name:        "Mystic cloth",
			AmountRange: [2]int{1, 3},
		},
	}
	giantmoleTertiaryDroptable = []core.Drop{
		{
			Name:        "Numbing root",
			AmountRange: [2]int{3, 6},
			Rate:        46.0 / 52.0,
			Bold:        true,
		},
		{
			Name:   "Clingy mole",
			Amount: 1,
			Rate:   5.0 / 52.0,
			Bold:   true,
		},
		{
			Name:   "Dragon 2h sword",
			Amount: 1,
			Rate:   1.0 / 52.0,
			Bold:   true,
		},
	}
	// Don't change order (we're removing the last element)
	giantmoleAlwaysDroptable = []core.Drop{
		{
			Name:   "Big bones",
			Amount: 1,
		},
		{
			Name:   "Mole claw",
			Amount: 1,
		},
		{
			Name:   "Mole skin",
			Amount: 3,
		},
		{
			Name:   "Mole nose",
			Amount: 1,
		},
	}
)

var GiantMoleCommand = &discordgo.ApplicationCommand{
	Name:        "giantmole",
	Description: "Emulate a Giant Mole drop. Warning: little is known about the droprate, it's not very reliable",
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
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "faladorshield-4",
			Description: "Set to false if you do not have Falador Shield 4 unlocked",
			Required:    false,
		},
	},
}

func GiantMole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var amount int64 = 1
	if amountOpt := core.GetOptionWithName(i.ApplicationCommandData().Options, "amount"); amountOpt.Name != "" {
		amount = amountOpt.IntValue()
	}

	// Replace this with max value later
	if amount > core.MAX_AMOUNT_ROLLS || amount < 1 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The amount has to be between 1 and %v", core.MAX_AMOUNT_ROLLS),
			},
		})
		return
	}

	drops := core.EmulateDropGiantMole(amount, giantmoleTertiaryDroptable, giantmoleRareDroptable, giantmoleUncommonDroptable, giantmoleCommonDroptable)

	// If option isn't there or it's true -> add guarantees
	if enableGuarantees := core.GetOptionWithName(i.ApplicationCommandData().Options, "enable-guarantees"); enableGuarantees.Name == "" || enableGuarantees.BoolValue() {
		// Adjust always table if faladorshield option is present and it's false
		if faladorshield4 := core.GetOptionWithName(i.ApplicationCommandData().Options, "faladorshield-4"); faladorshield4.Name != "" && !faladorshield4.BoolValue() {
			core.AddAlwaysDroptable(amount, &drops, core.RemoveDropFromTable(giantmoleAlwaysDroptable, len(giantmoleAlwaysDroptable)-1))
		} else {
			core.AddAlwaysDroptable(amount, &drops, giantmoleAlwaysDroptable)
		}
	}

	dropWithPrice, total, ok := core.AmountToPrice(drops)
	core.SortDrops(&dropWithPrice)
	content := core.MakeDropList(dropWithPrice, drops, total, ok)

	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: giantmoleUrl,
		},
		Title:       fmt.Sprintf("You killed Giant Mole %v times and you got:", amount),
		Description: content,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})
}
