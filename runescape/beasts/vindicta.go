package beasts

import (
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"
)

var (
	vindictaCommonDroptable = []Drop{
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

	vindictaUncommonDroptable = []Drop{
		{
			Name:        "Phasmatite stone spirit",
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
			Name:        "Large plated rune salvage",
			AmountRange: [2]int{8, 15},
		},
		{
			Name:        "Dragon bones",
			AmountRange: [2]int{150, 250},
		},
		{
			Name:        "Black dragonhide",
			AmountRange: [2]int{25, 44},
		},
	}

	vindictaRareDroptable = []Drop{
		{
			Rate: 1.0 / 256.0,
			Name: "Dormant anima core helm",
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Dormant anima core body",
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Dormant anima core legs",
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Dragon rider lance",
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Crest of Zaros",
		},
		{
			Rate: 1.0 / 256.0,
			Name: "Zarosian essence",
		},
	}

	vindictaAlwaysDroptable = []Drop{
		{
			Name:   "Dragon bones",
			Amount: 1,
		},
	}
)

var VindictaCommand = &discordgo.ApplicationCommand{
	Name:        "vindicta",
	Description: "Emulate a vindicta drop with full reputation",
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

func Vindicta(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var amount int64 = 1
	if lenOpt := len(i.ApplicationCommandData().Options); lenOpt >= 1 {
		amount = i.ApplicationCommandData().Options[0].IntValue()
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

	drops := EmulateDrop(commonRateWithoutRare, amount, vindictaRareDroptable, vindictaUncommonDroptable, vindictaCommonDroptable)

	if lenOpt := len(i.ApplicationCommandData().Options); lenOpt < 2 || i.ApplicationCommandData().Options[1].BoolValue() {
		AddAlwaysDroptable(amount, &drops, vindictaAlwaysDroptable)
	}

	dropWithPrice, ok := AmountToPrice(drops)
	if !ok {
		// Add content saying it's incomplete
	}

	sort.Slice(dropWithPrice, func(i, j int) bool {
		comp, _ := dropWithPrice[i].price.Compare(dropWithPrice[j].price)
		return comp == 1
	})

	content := "You got:\n"
	for _, d := range dropWithPrice {
		content += fmt.Sprintf("%v %v: %v gp\n", drops[d.name].Amount, d.name, d.price)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}
