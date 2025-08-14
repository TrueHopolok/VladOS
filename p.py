suits = ["Hearts", "Diamonds", "Clubs", "Spades"]
values = ["Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"]

for val in values:
    for suit in suits:
        print("\t{Suit: %s, Value: %s}," % (suit, val))