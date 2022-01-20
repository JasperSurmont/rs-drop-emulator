package main

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

var (
	nexUrl     = "https://runescape.wiki/images/Nex.png?67c8a"
	nexPotions = []Drop{
		{
			Name:   "Saradomin brew (4)",
			Amount: 10,
			Rate:   7.0 / 128.0,
			OtherDrops: []Drop{
				{
					Name:   "Super restore (4)",
					Amount: 30,
				},
			},
		},
		{
			Name:   "Super restore (4)",
			Amount: 10,
			Rate:   29.0 / 128.0,
			OtherDrops: []Drop{
				{
					Name:   "Saradomin brew (4)",
					Amount: 30,
				},
			},
		},
	}
	nexCommonDroptable = []Drop{
		{
			Name:   "Magic logs",
			Amount: 375,
		},
		{
			Name:   "Phasmatite stone spirit",
			Amount: 20,
		},
		{
			Name:   "Necrite stone spirit",
			Amount: 20,
		},
		{
			Name:   "Green dragonhide",
			Amount: 400,
		},
		{
			Name:   "Uncut dragonstone",
			Amount: 20,
		},
	}
	nexUncommonDroptable = []Drop{
		{
			Name:   "Onyx bolts (e)",
			Amount: 375,
		},
		{
			Name:   "Grimy avantoe",
			Amount: 75,
			OtherDrops: []Drop{
				{
					Name:   "Grimy dwarf weed",
					Amount: 75,
				},
			},
		},
		{
			Name:   "Grimy torstol",
			Amount: 40,
		},
		{
			Name:   "Torstol seed",
			Amount: 12,
		},
		{
			Name:   "Magic seed",
			Amount: 5,
		},
	}

	nexRareDroptable = []Drop{
		{
			Rate: 1.0 / 384.0,
			Name: "Torva full helm",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Torva platebody",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Torva platelegs",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Torva boots",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Torva gloves",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Pernix cowl",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Pernix body",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Pernix chaps",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Pernix boots",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Pernix gloves",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Virtus mask",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Virtus robe top",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Virtus robe legs",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Virtus boots",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Virtus gloves",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Virtus wand",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Virtus book",
			Bold: true,
		},
		{
			Rate: 1.0 / 384.0,
			Name: "Zaryte bow",
			Bold: true,
		},
	}

	nexAlwaysDroptable = []Drop{
		{
			Name:   "Big bones",
			Amount: 1,
		},
	}
)

var NexCommand = &discordgo.ApplicationCommand{
	Name:        "nex",
	Description: "Emulate a Nex drop with full reputation",
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

func Nex(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	drops := emulateDropNex(amount, nexRareDroptable, nexUncommonDroptable, nexCommonDroptable, nexPotions)
	if enableGuarantees := GetOptionWithName(i.ApplicationCommandData().Options, "enable-guarantees"); enableGuarantees.Name == "" || enableGuarantees.BoolValue() {
		AddAlwaysDroptable(amount, &drops, nexAlwaysDroptable)
	}

	dropWithPrice, total, ok := AmountToPrice(drops)
	SortDrops(&dropWithPrice)
	content := MakeDropList(dropWithPrice, drops, total, ok)

	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: nexUrl,
		},
		Title:       fmt.Sprintf("You killed Nex %v times and you got:", amount),
		Description: content,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})
	log.Info("command executed", "command", NexCommand.Name)

}

func emulateDropNex(amount int64, rareDroptable []Drop, uncommonDroptable []Drop, commonDroptable []Drop, nexPotions []Drop) map[string]*Drop {
	var drops map[string]*Drop = make(map[string]*Drop)

	for i := int64(0); i < amount; i++ {
		var sum float64 = 0
		for _, d := range rareDroptable {
			sum += d.Rate
		}

		var drop Drop

		roll := rand.Float64()

		// We split up the interval: [0, sum) = rare, [sum, sum+uncommonrate) = uncommon, [sum+uncommon, 1) = common
		if roll < sum {
			drop = rareDroptable[rand.Intn(len(rareDroptable))]
		} else if roll < sum+UncommonRateWithoutRare {
			drop = uncommonDroptable[rand.Intn(len(uncommonDroptable))]
		} else {
			// Compute if we have a potion drop, if not just common drop
			var sum2 float64 = 0
			for _, d := range nexPotions {
				sum2 += d.Rate
			}
			if roll2 := rand.Float64(); roll2 < sum2 {
				drop = determineDropWithRates(rand.Float64(), &nexPotions)
			} else {
				drop = commonDroptable[rand.Intn(len(commonDroptable))]
			}
		}

		drop.SetAmount()
		addDropValueToMap(drops, &drop)

		// Add drops that always go together
		for _, d := range drop.OtherDrops {
			d.SetAmount()
			addDropValueToMap(drops, &d)
		}
	}
	return drops
}
