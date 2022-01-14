package beasts

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

var (
	commonUncommonRatio float32 = 2.0 / 3.0

	commonDroptable = []Drop{
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
	}

	uncommonDroptable = []Drop{
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

	rareDroptable = []Drop{
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

	alwaysDrop = []Drop{
		{
			Name: "Dragon bones",
		},
	}
)

var VindictaCommand *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "vindicta",
	Description: "Emulate a vindicta drop with full reputation",
}

func Vindicta(s *discordgo.Session, i *discordgo.InteractionCreate) {
	drop := emulateDrop()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You got: %v %v", drop.Amount, drop.Name),
		},
	})
}

func emulateDrop() *Drop {
	var sum float32 = 0
	for _, d := range rareDroptable {
		sum += d.Rate
	}

	var drop Drop
	log.Debugf("sum: %v", sum)

	// Rare drop is chosen
	if rand.Float32() < sum {
		drop = rareDroptable[rand.Intn(len(rareDroptable))]
	} else if rand.Float32() > commonUncommonRatio {
		drop = uncommonDroptable[rand.Intn(len(uncommonDroptable))]
	} else {
		drop = commonDroptable[rand.Intn(len(commonDroptable))]
	}

	drop.SetAmount()
	return &drop
}
