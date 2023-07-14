package kulll

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/angch/multibot/pkg/bothandler"
)

var triggers = []string{
	"m0in", "moin", "morning",
	"Selamat pagi", "ohayo",
	"早安", "boin", "yawn",
	"おはよう", "Guten Morgen",
	"zzzz", "Magandang umaga",
	"Goedemorgen", "Goeiemorgen",
	"καλημέρα", "Kaliméra", "kalimera",
	"မင်္ဂလာပါ", "mingalaba", "mingalabar",
}

var handles = []string{
	// "faz",
	// "tech_tarik",
}

type History struct {
	Input string
}

var lock = sync.Mutex{}
var history map[string]History
var myrand = rand.New(rand.NewSource(time.Now().UnixNano()))

func init() {
	bothandler.RegisterCatchallHandler(KulllHandler)
	load()
	// math.Rand()
}

const savefile string = "kulll.js"

func load() {
	lock.Lock()
	defer lock.Unlock()

	history = make(map[string]History)
	f, err := os.Open(savefile)
	if err == nil {
		b, err := io.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(b, &history)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
	}
}

func save() {
	lock.Lock()
	defer lock.Unlock()

	f, err := os.OpenFile(savefile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	b, _ := json.Marshal(history)
	_, err = f.Write(b)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
}

func KulllHandler(request bothandler.Request) string {
	input := request.Content
	// Jan 2 15:04:05 2006 MST
	today := time.Now().Local().Format("20060102")
	key := fmt.Sprintf("%s/%s/%s", request.Platform, request.Channel, today)
	lock.Lock()
	_, ok := history[key]
	lock.Unlock()
	if ok {
		return ""
	}

	i := strings.ToLower(input)

	count := 0
	for _, v := range triggers {
		// Trigger is also the response, so some of them are in caps.
		if strings.Contains(i, strings.ToLower(v)) {
			count++
		}
	}

	uncount := 0
	for _, v := range handles {
		// Trigger is also the response, so some of them are in caps.
		if strings.Contains(i, strings.ToLower(v)) {
			uncount++
		}
	}

	if count >= 1 && uncount == 0 {
		pick := myrand.Intn(len(triggers))
		lock.Lock()
		history[key] = History{i}
		lock.Unlock()
		go save()
		return triggers[pick]
	}

	return ""
}
