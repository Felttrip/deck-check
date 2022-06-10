package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type card struct {
	name     string
	quantity int
}

type deck struct {
	mainDeck  []card
	sideboard []card
}

type pool []card

const (
	MAIN_DECK_MIN = 40
	SIDEBOARD_MAX = 15

	//lands
	PLAINS   = "Plains"
	ISLAND   = "Island"
	SWAMP    = "Swamp"
	MOUNTAIN = "Mountain"
	FOREST   = "Forest"
)

var LANDS = []string{PLAINS, ISLAND, SWAMP, MOUNTAIN, FOREST}

func main() {
	deckFile := flag.String("deck", "deck.txt", "realitive path of deck file")
	poolFile := flag.String("pool", "pool.txt", "realitive path of pool file")
	flag.Parse()

	d := loadDeck(*deckFile)
	p := loadPool(*poolFile)
	err := checkDeckAndSideboardSize(d)
	if err != nil {
		fmt.Println(err)
		return
	}
	errs := checkDeckInPool(d, p)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)

		}
		return
	}

	fmt.Println("Deck is Valid")
}

func checkDeckAndSideboardSize(d deck) error {
	if numCards(d.mainDeck) < MAIN_DECK_MIN {
		return fmt.Errorf("Main Deck is under limit of %d", MAIN_DECK_MIN)
	}
	if numCards(d.sideboard) > SIDEBOARD_MAX {
		return fmt.Errorf("Sideboard exceeds limit of %d", SIDEBOARD_MAX)
	}
	return nil
}

func checkDeckInPool(d deck, p pool) []error {
	err := []error{}
	for _, c := range d.mainDeck {
		inPool := useCardFromPool(c, p)
		if !inPool {
			err = append(err, fmt.Errorf("Main Deck Card %s not in pool or there are too many copies used", c.name))
		}
	}
	for _, c := range d.sideboard {
		inPool := useCardFromPool(c, p)
		if !inPool {
			err = append(err, fmt.Errorf("Sideboard Card %s not in pool or there are too many copies used", c.name))
		}
	}
	return err
}

//returns true if the card was available to be used in the pool, false if the card is unavailable
func useCardFromPool(c card, p pool) bool {
	//if its a land skip it all pools have land
	if slices.Contains(LANDS, c.name) {
		return true
	}
	for i, pc := range p {
		if c.name == pc.name && pc.quantity-c.quantity >= 0 {
			//remove the number of cards used incase there are dupes sideboarded
			//update the actual val of the quantity cant use the range var
			(&p[i]).quantity = pc.quantity - c.quantity
			return true
		}
	}
	return false
}

func numCards(cards []card) int {
	num := 0
	for _, c := range cards {
		num += c.quantity
	}
	return num
}

func loadDeck(p string) deck {
	d := deck{
		mainDeck:  []card{},
		sideboard: []card{},
	}
	loadMainDeck := false

	file, err := os.Open(p)
	if err != nil {
		fmt.Println("No Deck Submitted")
		os.Exit(0)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		raw := scanner.Text()
		if len(raw) == 0 {
			continue
		}
		if raw == "Deck" || raw == "Sideboard" {
			loadMainDeck = raw == "Deck"
			continue
		}
		c, err := parseCard(raw)
		if err != nil {
			log.Fatal(err)
		}
		if loadMainDeck {
			d.mainDeck = append(d.mainDeck, c)
		} else {
			d.sideboard = append(d.sideboard, c)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return d
}

func loadPool(f string) pool {
	p := pool{}

	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		raw := scanner.Text()
		if len(raw) == 0 {
			continue
		}
		c, err := parseCard(raw)
		if err != nil {
			log.Fatal(err)
		}
		p = append(p, c)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return p
}

func parseCard(s string) (card, error) {
	quantity, name, found := strings.Cut(s, " ")
	if !found {
		return card{}, fmt.Errorf("issue parsing card %v", s)
	}
	q, err := strconv.Atoi(quantity)
	if err != nil {
		return card{}, fmt.Errorf("issue converting quantity %v", s)
	}

	return card{
		name:     name,
		quantity: q,
	}, nil
}
