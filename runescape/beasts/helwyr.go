package beasts

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	helwyrCommonDroptable = []Drop{
		{
			Name:        "Drakolith stone spirit",
			AmountRange: [2]int{15, 25},
		},
		{
			Name:        "Oricalcite stone spirit",
			AmountRange: [2]int{15, 25},
		},
		{
			Name:        "Uncut diamond",
			AmountRange: [2]int{20, 30},
		},
		{
			Name:        "Grimy dwarf weed",
			AmountRange: [2]int{20, 30},
		},
		{
			Name:        "Raw shark",
			AmountRange: [2]int{45, 60},
		},
	}

	helwyrUncommonDroptable = []Drop{
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
			AmountRange: [2]int{175, 350},
		},
		{
			Name:        "Large bladed rune salvage",
			AmountRange: [2]int{10, 20},
		},
		{
			Name:        "Crystal key",
			AmountRange: [2]int{2, 4},
		},
		{
			Name:        "Grimy lantadyme",
			AmountRange: [2]int{90, 120},
		},
	}

	helwyrRareDroptable = []Drop{
		{
			Rate: 1.0 / 179.0,
			Name: "Dormant anima core helm",
		},
		{
			Rate: 1.0 / 179.0,
			Name: "Dormant anima core body",
		},
		{
			Rate: 1.0 / 179.0,
			Name: "Dormant anima core legs",
		},
		{
			Rate: 1.0 / 179.0,
			Name: "Orb of the Cywir elders",
		},
		{
			Rate: 1.0 / 179.0,
			Name: "Crest of Seren",
		},
		{
			Rate: 1.0 / 179.0,
			Name: "Wand of the Cywir elders",
		},
		{
			Rate: 1.0 / 179.0,
			Name: "Serenic essence",
		},
	}

	helwyrAlwaysDroptable = []Drop{
		{
			Name:   "Bones",
			Amount: 1,
		},
	}
)

var HelwyrCommand = &discordgo.ApplicationCommand{
	Name:        "helwyr",
	Description: "Emulate a helwyr drop with full reputation",
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

func Helwyr(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var amount int64 = 1
	if opt0 := i.ApplicationCommandData().Options[0]; opt0 != nil {
		amount = opt0.IntValue()
	}

	// Replace this with max value later
	if amount > MAX_AMOUNT_ROLLS || amount < 1 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The maximum amount has to be between 1 and %v", MAX_AMOUNT_ROLLS),
			},
		})
		return
	}

	drops := EmulateDrop(commonRateWithoutRare, amount, helwyrRareDroptable, helwyrUncommonDroptable, helwyrCommonDroptable)

	if opt1 := i.ApplicationCommandData().Options[1]; opt1 == nil || opt1.BoolValue() {
		AddAlwaysDroptable(amount, &drops, helwyrAlwaysDroptable)
	}

	content := "You got:\n"
	for _, d := range drops {
		content += fmt.Sprintf("%v %v\n", d.Amount, d.Name)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(content),
		},
	})
}
