package apod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/angch/discordbot/pkg/bothandler"
)

var posts map[string]ApodPost

func init() {
	bothandler.RegisterPlatformRegisteration(AddMessagePlatform)

	// go Tick()

	posts = make(map[string]ApodPost)
	f, err := os.Open("posts.js")
	if err == nil {
		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(b, &posts)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
	}
	// log.Printf("%+v\n", posts)

	go Apod()
}

var MessagePlatforms = []bothandler.MessagePlatform{}

func AddMessagePlatform(m bothandler.MessagePlatform) {
	// log.Printf("Registering module apod in %+v\n", m)

	MessagePlatforms = append(MessagePlatforms, m)
}

// For testing.
func Tick() {
	for {
		time.Sleep(10 * time.Second)
		for _, v := range MessagePlatforms {
			v.Send("tick")
		}
	}
}

type ApodPost struct {
	Text     string
	ImageURL string
}

func findlinks() string {

}

func doYMD(y, m, d int) *ApodPost {
	if y > 2000 {
		// Perlism
		y -= 2000
	}

	url := fmt.Sprintf("https://apod.nasa.gov/apod/ap%02d%02d%02d.html", y, m, d)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil
	}
	if resp.StatusCode != 200 {
		log.Println("Not ready yet ", resp.Status)
	}
	if resp.Body == nil {
		log.Println("Body is empty")
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	resp.Body.Close()
	imgUrl := ""
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		img := s.AttrOr("href", "")
		lowerImg := strings.ToLower(img)
		if !strings.HasPrefix(lowerImg, "http") && (strings.HasSuffix(lowerImg, ".jpg") || strings.HasSuffix(lowerImg, ".png") && strings.HasPrefix(img, "image")) {
			// We usually only want the *first* one.
			if imgUrl == "" {
				imgUrl = "https://apod.nasa.gov/apod/" + img
			} else {
				// 2nd link to image which probably is wrong
			}
		}
	})
	title := ""
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		title = strings.TrimSpace(s.Text())
	})

	if imgUrl == "" {
		// Is it a video?
		doc.Find("iframe").Each(func(i int, s *goquery.Selection) {
			img := s.AttrOr("src", "")
			lowerImg := strings.ToLower(img)
			if strings.HasPrefix(lowerImg, "https://www.youtube.com") {
				imgUrl = img

				if strings.Contains(lowerImg, "/embed/") {
					// FIXME: Lazy.
					// https://www.youtube.com/embed/zIqG42AD4Gw?rel=0
					// to
					// https://www.youtube.com/watch?v=zIqG42AD4Gw
					imgUrl = strings.ReplaceAll(imgUrl, "/embed/", "/watch?v=")
				}
			}
		})
	}

	if title != "" && imgUrl != "" {
		post := ApodPost{
			Text:     title,
			ImageURL: imgUrl,
		}
		return &post
	} else {
		log.Println("No title", title, "or url", imgUrl, "in", url)
	}
	return nil
}

func Apod() {
	// Let all the platforms get initialized first
	time.Sleep(5 * time.Second)

	for {
		today := time.Now() // Yes, I know. timezone.
		d := today.Day()
		m := int(today.Month())
		y := today.Year()
		key := fmt.Sprintf("%04d%02d%02d", y, m, d)
		_, exists := posts[key]
		if exists {
			log.Println("Done for today")
			time.Sleep(1 * time.Hour)
			continue
		}

		log.Printf("Doing %d %d %d\n", y, m, d)

		p := doYMD(y, m, d)
		if p != nil {
			posts[key] = *p
			f, err := os.OpenFile("posts.js", os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				log.Fatal(err)
			}
			b, _ := json.Marshal(posts)
			_, err = f.Write(b)
			if err != nil {
				log.Fatal(err)
			}
			f.Close()

			log.Printf("%+v\n", p)

			text := fmt.Sprintf("%s %s", p.Text, p.ImageURL)
			for _, v := range MessagePlatforms {
				v.SendWithOptions(text, bothandler.SendOptions{Silent: true})
			}
		}

		log.Println("Sleeping 5 minutes")
		time.Sleep(5 * time.Minute)
	}
}
