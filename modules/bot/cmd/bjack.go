package cmd

import (
	"bytes"
	"cmp"
	"encoding/gob"
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
	"time"

	"github.com/TrueHopolok/VladOS/modules/db/dbconvo"
	"github.com/TrueHopolok/VladOS/modules/db/dbstats"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type (
	bjackCardSuit  int
	bjackCardValue int
	bjackCard      struct {
		Suit  bjackCardSuit
		Value bjackCardValue
	}
	bjackConvoStatus struct {
		DealerHand  []bjackCard
		PlayersHand []bjackCard
		CardsDeck   []bjackCard
		CardsRemain int
	}
)

const (
	bjackCardValueAce bjackCardValue = iota + 1
	bjackCardValueTwo
	bjackCardValueThree
	bjackCardValueFour
	bjackCardValueFive
	bjackCardValueSix
	bjackCardValueSeven
	bjackCardValueEight
	bjackCardValueNine
	bjackCardValueTen
	bjackCardValueJack
	bjackCardValueQueen
	bjackCardValueKing

	bjackCardSuitHearts bjackCardSuit = iota - 13
	bjackCardSuitDiamonds
	bjackCardSuitClubs
	bjackCardSuitSpades
)

var bjackCardDeck = []bjackCard{
	{Suit: bjackCardSuitHearts, Value: bjackCardValueAce},
	{Suit: bjackCardSuitDiamonds, Value: bjackCardValueAce},
	{Suit: bjackCardSuitClubs, Value: bjackCardValueAce},
	{Suit: bjackCardSuitSpades, Value: bjackCardValueAce},
	{Suit: bjackCardSuitHearts, Value: bjackCardValueTwo},
	{Suit: bjackCardSuitDiamonds, Value: bjackCardValueTwo},
	{Suit: bjackCardSuitClubs, Value: bjackCardValueTwo},
	{Suit: bjackCardSuitSpades, Value: bjackCardValueTwo},
	{Suit: bjackCardSuitHearts, Value: bjackCardValueThree},
	{Suit: bjackCardSuitDiamonds, Value: bjackCardValueThree},
	{Suit: bjackCardSuitClubs, Value: bjackCardValueThree},
	{Suit: bjackCardSuitSpades, Value: bjackCardValueThree},
	{Suit: bjackCardSuitHearts, Value: bjackCardValueFour},
	{Suit: bjackCardSuitDiamonds, Value: bjackCardValueFour},
	{Suit: bjackCardSuitClubs, Value: bjackCardValueFour},
	{Suit: bjackCardSuitSpades, Value: bjackCardValueFour},
	{Suit: bjackCardSuitHearts, Value: bjackCardValueFive},
	{Suit: bjackCardSuitDiamonds, Value: bjackCardValueFive},
	{Suit: bjackCardSuitClubs, Value: bjackCardValueFive},
	{Suit: bjackCardSuitSpades, Value: bjackCardValueFive},
	{Suit: bjackCardSuitHearts, Value: bjackCardValueSix},
	{Suit: bjackCardSuitDiamonds, Value: bjackCardValueSix},
	{Suit: bjackCardSuitClubs, Value: bjackCardValueSix},
	{Suit: bjackCardSuitSpades, Value: bjackCardValueSix},
	{Suit: bjackCardSuitHearts, Value: bjackCardValueSeven},
	{Suit: bjackCardSuitDiamonds, Value: bjackCardValueSeven},
	{Suit: bjackCardSuitClubs, Value: bjackCardValueSeven},
	{Suit: bjackCardSuitSpades, Value: bjackCardValueSeven},
	{Suit: bjackCardSuitHearts, Value: bjackCardValueEight},
	{Suit: bjackCardSuitDiamonds, Value: bjackCardValueEight},
	{Suit: bjackCardSuitClubs, Value: bjackCardValueEight},
	{Suit: bjackCardSuitSpades, Value: bjackCardValueEight},
	{Suit: bjackCardSuitHearts, Value: bjackCardValueNine},
	{Suit: bjackCardSuitDiamonds, Value: bjackCardValueNine},
	{Suit: bjackCardSuitClubs, Value: bjackCardValueNine},
	{Suit: bjackCardSuitSpades, Value: bjackCardValueNine},
	{Suit: bjackCardSuitHearts, Value: bjackCardValueTen},
	{Suit: bjackCardSuitDiamonds, Value: bjackCardValueTen},
	{Suit: bjackCardSuitClubs, Value: bjackCardValueTen},
	{Suit: bjackCardSuitSpades, Value: bjackCardValueTen},
	{Suit: bjackCardSuitHearts, Value: bjackCardValueJack},
	{Suit: bjackCardSuitDiamonds, Value: bjackCardValueJack},
	{Suit: bjackCardSuitClubs, Value: bjackCardValueJack},
	{Suit: bjackCardSuitSpades, Value: bjackCardValueJack},
	{Suit: bjackCardSuitHearts, Value: bjackCardValueQueen},
	{Suit: bjackCardSuitDiamonds, Value: bjackCardValueQueen},
	{Suit: bjackCardSuitClubs, Value: bjackCardValueQueen},
	{Suit: bjackCardSuitSpades, Value: bjackCardValueQueen},
	{Suit: bjackCardSuitHearts, Value: bjackCardValueKing},
	{Suit: bjackCardSuitDiamonds, Value: bjackCardValueKing},
	{Suit: bjackCardSuitClubs, Value: bjackCardValueKing},
	{Suit: bjackCardSuitSpades, Value: bjackCardValueKing},
}

var CommandBjack Command = Command{
	InfoBrief: "game of a blackjack",
	InfoFull: `
 /bjack
Play a blackjack against a dealer. 
No need to bet any money, since you will bet your score streak like in dice and slots.

Rules of the blackjack can be read here:
https://en.wikipedia.org/wiki/Blackjack#Player_decisions

This variation differences from the regular blackjack:
1) Only 1 card deck is in the game so counting would be viable.
2) Splitting your hand is removed since your only goal is to continue win streak and not to earn net positive.
3) Wining gives you 3 points to score streak.
4) Wining on Double Down move gives you 6 points to score streak. 
5) Draw gives you 1 point to score streak.

On losing score is reset.

Has a leaderboard to count largest score streak.`,
	handler: func(ctx *telegohandler.Context, update telego.Update) error {
		bot, chatID, _, valid, err := utilStart(ctx, update, "bjack", 0)
		if !valid {
			return err
		}

		status := bjackStart()

		_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID, append(
			bjackOutputHands(status, false),
			bjackOutputActions()...,
		)...))
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(status); err != nil {
			return fmt.Errorf("gob encoder: %w", err)
		}
		return dbconvo.Busy(update.Message.From.ID, "bjack", buf.Bytes())
	},
	conversation: func(ctx *telegohandler.Context, update telego.Update) error {
		slog.Debug("bot handler", "upd", update.UpdateID, "command", "bjack")
		bot := ctx.Bot()
		chatID := update.Message.Chat.ChatID()
		userID := update.Message.From.ID

		cs := ctx.Value(ctxValueConvoStatus{}).(dbconvo.Status)
		getbuf := bytes.NewBuffer(cs.Data)
		dec := gob.NewDecoder(getbuf)
		var status bjackConvoStatus
		if err := dec.Decode(&status); err != nil {
			return fmt.Errorf("gob decoder: %w", err)
		}

		addingCard := false
		standing := false
		doubled := false
		switch update.Message.Text {
		case "hit":
			addingCard = true
		case "doubledown":
			if len(status.PlayersHand) > 2 {
				_, err := bot.SendMessage(ctx, tu.MessageWithEntities(chatID, append([]tu.MessageEntityCollection{tu.Entity("You can't double down after the 1st move.")}, bjackOutputActions()...)...))
				return err
			}
			doubled = true
			addingCard = true
			standing = true
		case "stand":
			standing = true
		default:
			_, err := bot.SendMessage(ctx, tu.MessageWithEntities(chatID, append([]tu.MessageEntityCollection{tu.Entity("Non existing action.")}, bjackOutputActions()...)...))
			return err
		}

		if addingCard {
			status.CardsRemain--
			status.PlayersHand = append(status.PlayersHand, status.CardsDeck[status.CardsRemain])
			if bjackScore(status.PlayersHand) > 21 {
				if err := cmp.Or(dbstats.Update("bjack", userID, update.Message.From.FirstName, update.Message.From.Username, 0), dbconvo.Free(userID)); err != nil {
					return err
				}
				msgText, err := utilOutputDice("bjack", userID, false)
				if err != nil {
					return err
				}
				_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID, bjackOutputHands(status, true)...))
				if err != nil {
					return err
				}
				_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID, msgText...))
				return err
			}
		}

		if standing {
			for bjackScore(status.DealerHand) < 17 {
				status.CardsRemain--
				status.DealerHand = append(status.DealerHand, status.CardsDeck[status.CardsRemain])
			}
			hasWon := bjackScore(status.DealerHand) > 21 || bjackScore(status.DealerHand) <= bjackScore(status.PlayersHand)
			finScore := 0
			if hasWon {
				if bjackScore(status.DealerHand) == bjackScore(status.PlayersHand) {
					finScore = 1
				} else {
					finScore = 3
				}
			}
			if doubled {
				finScore *= 2
			}

			if err := cmp.Or(dbstats.Update("bjack", userID, update.Message.From.FirstName, update.Message.From.Username, finScore), dbconvo.Free(userID)); err != nil {
				return err
			}
			msgText, err := utilOutputDice("bjack", userID, hasWon)
			if err != nil {
				return err
			}
			_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID, bjackOutputHands(status, true)...))
			if err != nil {
				return err
			}
			_, err = bot.SendMessage(ctx, tu.MessageWithEntities(chatID, msgText...))
			return err
		}

		_, err := bot.SendMessage(ctx, tu.MessageWithEntities(chatID, append(
			bjackOutputHands(status, false),
			bjackOutputActions()...,
		)...))
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(status); err != nil {
			return fmt.Errorf("gob encoder: %w", err)
		}
		return dbconvo.Busy(update.Message.From.ID, "bjack", buf.Bytes())
	},
	cancelation: func(ctx *telegohandler.Context, update telego.Update) error {
		return dbstats.Update("bjack", update.Message.From.ID, update.Message.From.FirstName, update.Message.From.Username, 0)
	},
}

func bjackStart() bjackConvoStatus {
	status := bjackConvoStatus{
		DealerHand:  make([]bjackCard, 2),
		PlayersHand: make([]bjackCard, 2),
		CardsDeck:   make([]bjackCard, len(bjackCardDeck)),
		CardsRemain: len(bjackCardDeck) - 4,
	}
	copy(status.CardsDeck, bjackCardDeck)
	rand.New(rand.NewSource(time.Now().UnixNano())).Shuffle(len(bjackCardDeck), func(i, j int) {
		status.CardsDeck[i], status.CardsDeck[j] = status.CardsDeck[j], status.CardsDeck[i]
	})
	copy(status.PlayersHand, status.CardsDeck[status.CardsRemain+2:])
	copy(status.DealerHand, status.CardsDeck[status.CardsRemain:status.CardsRemain+2])
	return status
}

func bjackScore(hand []bjackCard) int {
	hasAce := false
	score := 0
	for _, card := range hand {
		switch card.Value {
		case bjackCardValueAce:
			score += 1
			hasAce = true
		case bjackCardValueJack, bjackCardValueQueen, bjackCardValueKing:
			score += 10
		default:
			score += int(card.Value)
		}
	}
	if hasAce {
		if score+10 <= 21 {
			score += 10
		}
	}
	return score
}

func bjackOutputActions() []tu.MessageEntityCollection {
	return []tu.MessageEntityCollection{
		tu.Entity("\n\nYou can use any of those actions:\n"),
		tu.Entity("doubledown").Code(),
		tu.Entity("- hit+stand at 1st move;\n"),
		tu.Entity("hit").Code(),
		tu.Entity("- add 1 card to your hand;\n"),
		tu.Entity("stand").Code(),
		tu.Entity(" - finish the game;\n"),
		tu.Entity("/cancel - surrender and cancel the game."),
	}
}

func bjackOutputHands(status bjackConvoStatus, showDealer bool) []tu.MessageEntityCollection {
	var msgText []tu.MessageEntityCollection

	msgText = append(msgText,
		tu.Entity("Dealer's hand:\n").Bold(),
	)

	if showDealer {
		for _, card := range status.DealerHand {
			msgText = append(msgText, tu.Entityf("%s\n", bjackCardOutput(card)))
		}
	} else {
		msgText = append(msgText,
			tu.Entityf("%s\n???\n", bjackCardOutput(status.DealerHand[0])),
		)
	}

	msgText = append(msgText,
		tu.Entity("\nYour's hand:\n").Bold(),
	)
	for _, card := range status.PlayersHand {
		msgText = append(msgText, tu.Entityf("%s\n", bjackCardOutput(card)))
	}

	msgText = append(msgText,
		tu.Entity("\nYour hand score: "),
		tu.Entityf("%d\n", bjackScore(status.PlayersHand)).Bold(),
	)

	if showDealer {
		msgText = append(msgText,
			tu.Entity("\nDealer's hand score: "),
			tu.Entityf("%d\n", bjackScore(status.DealerHand)).Bold(),
		)
	}

	return msgText
}

func bjackCardOutput(card bjackCard) string {
	var result string

	switch card.Value {
	case bjackCardValueAce:
		result = "A"
	case bjackCardValueJack:
		result = "J"
	case bjackCardValueQueen:
		result = "Q"
	case bjackCardValueKing:
		result = "K"
	default:
		result = strconv.Itoa(int(card.Value))
	}

	switch card.Suit {
	case bjackCardSuitHearts:
		result += "♥️"
	case bjackCardSuitClubs:
		result += "♣️"
	case bjackCardSuitDiamonds:
		result += "♦️"
	case bjackCardSuitSpades:
		result += "♠️"
	}

	return result
}
