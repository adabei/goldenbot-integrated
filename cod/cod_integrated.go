package cod

import (
	"github.com/adabei/goldenbot/events"
	"github.com/adabei/goldenbot/events/cod"
	"github.com/adabei/goldenbot/rcon"
	"log"
	"strconv"
	"strings"
)

var players = make([]string, 1)

type Integrated struct {
	requests chan rcon.RCONQuery
	events   chan interface{}
}

func NewIntegrated(requests chan rcon.RCONQuery, ea events.Aggregator) *Integrated {
	i := new(Integrated)
	i.requests = requests
	i.events = ea.Subscribe(i)
	return i
}

func (i *Integrated) Setup() error {
	return nil
}

func (i *Integrated) maxClients() int {
	maxch := make(chan []byte)
	i.requests <- rcon.RCONQuery{Command: "serverinfo", Response: maxch}
	res := <-maxch
	for _, line := range strings.Split(string(res), "\n") {
		if strings.HasPrefix(line, "sv_maxclients") {
			max, err := strconv.Atoi(strings.Fields(line)[1])
			if err != nil {
				return -1
			} else {
				return max
			}
		}
	}
	return -1
}

func (i *Integrated) Start() {
	maxClients := i.maxClients()
	for maxClients == -1 {
		log.Println("integrated: failed to fetch sv_maxclients. Trying again.")
		maxClients = i.maxClients()
	}
	players = make([]string, maxClients)

	statusch := make(chan []byte)
	i.requests <- rcon.RCONQuery{Command: "status", Response: statusch}
	resp := <-statusch
	if resp != nil {
		for _, line := range strings.Split(string(resp), "\n")[4:] {
			if line != "" {
				p := strings.Fields(line)
				num, err := strconv.Atoi(p[0])
				if err != nil {
					log.Println("integrated: could not parse status, num is no integer")
					continue
				}
				players[num] = p[3]
			}
		}
	}

	for {
		ev := <-i.events
		switch ev := ev.(type) {
		case cod.Join:
			if num := firstFreeNum(players); num != -1 {
				players[num] = ev.GUID
				log.Println("integrated: guid", ev.GUID, "is assigned num", num)
			} else {
				log.Fatal("integrated: more players than originally possible on server")
			}
		case cod.Quit:
			num, _ := Num(ev.GUID)
			players[num] = ""
		// flush players list after every round
		case cod.ExitLevel:
			for i := range players {
				players[i] = ""
			}
		}
	}
}

// Returns the first empty index.
// The server should prevent joins if already full.
func firstFreeNum(p []string) int {
	for i, v := range p {
		if v == "" {
			return i
		}
	}

	// turns out it is full
	return -1
}

// Num returns the num for a player.
// If that player is not in our list we try to predict it.
// If the server is full, -1 and false will be returned, as no num will match.
func Num(id string) (int, bool) {
	for i, v := range players {
		if v == id {
			return i, true
		}
	}

	// guid not in our players list, predict it
	num := firstFreeNum(players)
	if num != -1 {
		log.Println("integrated: predicting num", num, "for guid", id)
		return num, true
	}

	return -1, false
}
