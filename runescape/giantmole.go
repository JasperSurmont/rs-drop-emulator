package runescape

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

//ref: https://archive.fo/0KkjH
// needs adjustment, what to do with clue scrolls? Numbing root, clinby mole and D2H add up to 100%, when is clue scroll rolled?
// unkown items are added to common
var (
	giantmoleUrl             = "https://runescape.wiki/images/Giant_Mole.png?6906f"
	giantmoleCommonDroptable = []Drop{
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
	giantmoleUncommonDroptable = []Drop{
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
	giantmoleRareDroptable = []Drop{
		{
			Name:        "Mystic cloth",
			AmountRange: [2]int{1, 3},
		},
	}
	giantmoleTertiaryDroptable = []Drop{
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
	giantmoleAlwaysDroptable = []Drop{
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

	drops := EmulateDropGiantMole(amount, giantmoleTertiaryDroptable, giantmoleRareDroptable, giantmoleUncommonDroptable, giantmoleCommonDroptable)

	// If option isn't there or it's true -> add guarantees
	if enableGuarantees := GetOptionWithName(i.ApplicationCommandData().Options, "enable-guarantees"); enableGuarantees.Name == "" || enableGuarantees.BoolValue() {
		// Adjust always table if faladorshield option is present and it's false
		if faladorshield4 := GetOptionWithName(i.ApplicationCommandData().Options, "faladorshield-4"); faladorshield4.Name != "" && !faladorshield4.BoolValue() {
			AddAlwaysDroptable(amount, &drops, RemoveDropFromTable(giantmoleAlwaysDroptable, len(giantmoleAlwaysDroptable)-1))
		} else {
			AddAlwaysDroptable(amount, &drops, giantmoleAlwaysDroptable)
		}
	}

	dropWithPrice, total, ok := AmountToPrice(drops)
	SortDrops(&dropWithPrice)
	content := MakeDropList(dropWithPrice, drops, total, ok)

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
	log.Info("command executed", "command", GiantMoleCommand.Name)
}

// There are 8 uncommons, 3 commons and 1 rare
// let's say common is 50%, uncommon is 48%, then each uncommon has 6% chance, and rare has 2% chance
const (
	giantmoleRareRate     = 2.0 / 100.0
	giantmoleUncommonRate = 48.0 / 100
	giantmoleTertiaryRate = 1.0 / 10.0
)

func EmulateDropGiantMole(amount int64, tertiaryDroptable []Drop, rareDroptable []Drop, uncommonDroptable []Drop, commonDroptable []Drop) map[string]*Drop {
	var drops map[string]*Drop = make(map[string]*Drop)

	for i := int64(0); i < amount; i++ {
		var drop Drop

		roll := rand.Float64()

		// We split up the interval: [0, sum) = rare, [sum, sum+uncommonrate) = uncommon, [sum+uncommon, 1) = common
		if roll < giantmoleRareRate {
			drop = rareDroptable[rand.Intn(len(rareDroptable))]
		} else if roll < giantmoleRareRate+giantmoleUncommonRate {
			drop = uncommonDroptable[rand.Intn(len(uncommonDroptable))]
		} else {
			drop = commonDroptable[rand.Intn(len(commonDroptable))]
		}

		drop.SetAmount()
		addDropValueToMap(drops, &drop)

		// We roll for tertiary drop
		roll = rand.Float64()
		if roll < giantmoleTertiaryRate {
			roll = rand.Float64() // Roll again to decide which drop
			curr := 0.0

			for _, d := range tertiaryDroptable {
				curr += d.Rate
				if roll < curr {
					d.SetAmount()
					addDropValueToMap(drops, &d)
					break
				}
			}
		}

	}
	return drops
}
