package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type Card struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type deck struct {
	deck      []Card
	sideboard []Card
}

type pool []Card

type sealedDeckResp struct {
	PoolId    string `json:"poolId"`
	Sideboard []Card `json:"sideboard"`
	Hidden    []Card `json:"hidden"`
	Deck      []Card `json:"deck"`
}

const (
	MAIN_DECK_MIN = 40
	SIDEBOARD_MAX = 15

	//lands
	PLAINS   = "Plains"
	ISLAND   = "Island"
	SWAMP    = "Swamp"
	MOUNTAIN = "Mountain"
	FOREST   = "Forest"

	SEALEDDECKTECH = "https://sealeddeck.tech/api/pools/"
)

var LANDS = []string{PLAINS, ISLAND, SWAMP, MOUNTAIN, FOREST}

func main() {
	deckFile := flag.String("deck-file", "", "realitive path of deck file")
	poolFile := flag.String("pool-file", "", "realitive path of pool file")
	deckId := flag.String("deck", "", "sealeddeck.tech id of deck")
	poolId := flag.String("pool", "", "sealeddeck.tech id of pool")
	flag.Parse()
	var d = deck{}
	var p = pool{}

	if *deckFile != "" {
		d = loadDeckFromFile(*deckFile)
	}
	if *poolFile != "" {
		p = loadPoolFromFile(*poolFile)
	}
	if *deckId != "" {
		d = loadDeckFromSealedDeck(*deckId)
	}
	if *poolId != "" {
		p = loadPoolFromSealedDeck(*poolId)
	}
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
	if numCards(d.deck) < MAIN_DECK_MIN {
		return fmt.Errorf("main Deck is under limit of %d", MAIN_DECK_MIN)
	}
	if numCards(d.sideboard) > SIDEBOARD_MAX {
		return fmt.Errorf("sideboard exceeds limit of %d", SIDEBOARD_MAX)
	}
	return nil
}

func checkDeckInPool(d deck, p pool) []error {
	err := []error{}
	for _, c := range d.deck {
		inPool := useCardFromPool(c, p)
		if !inPool {
			err = append(err, fmt.Errorf("main Deck Card %s not in pool or there are too many copies used", c.Name))
		}
	}
	for _, c := range d.sideboard {
		inPool := useCardFromPool(c, p)
		if !inPool {
			err = append(err, fmt.Errorf("sideboard Card %s not in pool or there are too many copies used", c.Name))
		}
	}
	return err
}

// returns true if the card was available to be used in the pool, false if the card is unavailable
func useCardFromPool(c Card, p pool) bool {
	//if its a land skip it all pools have land
	if slices.Contains(LANDS, c.Name) {
		return true
	}
	for i, pc := range p {
		if c.Name == pc.Name && pc.Count-c.Count >= 0 {
			//remove the number of cards used incase there are dupes sideboarded
			//update the actual val of the quantity cant use the range var
			(&p[i]).Count = pc.Count - c.Count
			return true
		}
	}
	return false
}

func numCards(cards []Card) int {
	num := 0
	for _, c := range cards {
		num += c.Count
	}
	return num
}

func loadDeckFromSealedDeck(id string) deck {
	deckResp := getCardsFromSealedDeck(id)
	d := deck{
		deck:      deckResp.Deck,
		sideboard: deckResp.Sideboard,
	}
	return d
}

func loadPoolFromSealedDeck(id string) pool {
	deckResp := getCardsFromSealedDeck(id)
	p := pool{}
	p = append(p, deckResp.Deck...)
	p = append(p, deckResp.Sideboard...)
	p = append(p, deckResp.Hidden...)
	return p
}

func getCardsFromSealedDeck(id string) sealedDeckResp {
	resp, err := http.Get(SEALEDDECKTECH + id)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var deckResp = sealedDeckResp{}
	if err := json.Unmarshal(body, &deckResp); err != nil {
		log.Fatal(err)
	}
	return deckResp
}

func loadDeckFromFile(f string) deck {
	d := deck{
		deck:      []Card{},
		sideboard: []Card{},
	}
	loadMainDeck := false

	file, err := os.Open(f)
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
			d.deck = append(d.deck, c)
		} else {
			d.sideboard = append(d.sideboard, c)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return d
}

func loadPoolFromFile(f string) pool {
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

func parseCard(s string) (Card, error) {
	quantity, name, found := strings.Cut(s, " ")
	if !found {
		return Card{}, fmt.Errorf("issue parsing card %v", s)
	}
	//strip off expansion like (ONE) 261
	name, _, _ = strings.Cut(name, " (")
	q, err := strconv.Atoi(quantity)
	if err != nil {
		return Card{}, fmt.Errorf("issue converting quantity %v", s)
	}

	return Card{
		Name:  name,
		Count: q,
	}, nil
}
