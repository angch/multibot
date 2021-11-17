package kulll

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/angch/discordbot/pkg/bothandler"
)

var triggers = []string{
	"m0in", "moin", "morning",
	"Selamat pagi!", "ohaiyo",
	"早安", "boin", "yawn",
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

func init() {
	bothandler.RegisterCatchallHandler(MoinHandler)
	load()
	rand.Seed(time.Now().Unix())
	// math.Rand()
}

const savefile string = "kulll.js"

func load() {
	lock.Lock()
	defer lock.Unlock()

	history = make(map[string]History)
	f, err := os.Open(savefile)
	if err == nil {
		b, err := ioutil.ReadAll(f)
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

func MoinHandler(request bothandler.Request) string {
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
		if strings.Contains(i, v) {
			count++
		}
	}

	uncount := 0
	for _, v := range handles {
		if strings.Contains(i, v) {
			uncount++
		}
	}

	if count >= 1 && uncount == 0 {
		pick := rand.Intn(len(triggers))
		lock.Lock()
		history[key] = History{i}
		lock.Unlock()
		go save()
		return triggers[pick]
	}

	return ""
}
