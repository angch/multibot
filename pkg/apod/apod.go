package apod

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/angch/multibot/pkg/bothandler"
)

var posts map[string]ApodPost

func init() {
	// go Tick()

	posts = make(map[string]ApodPost)
	f, err := os.Open("posts.js")
	if err == nil {
		b, err := io.ReadAll(f)
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

func GetMessagePlatforms() []bothandler.MessagePlatform {
	return bothandler.ActiveMessagePlatforms
}

// For testing.
func Tick() {
	for {
		time.Sleep(10 * time.Second)
		MessagePlatforms := GetMessagePlatforms()
		for _, v := range MessagePlatforms {
			v.Send("tick")
		}
	}
}

type ApodPost struct {
	Text        string
	ImageURL    string
	Description string
}

func ParseApod(body string) (*ApodPost, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	imgUrl := ""
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		img := s.AttrOr("href", "")
		lowerImg := strings.ToLower(img)
		if !strings.HasPrefix(lowerImg, "http") && (strings.HasSuffix(lowerImg, ".jpg") || strings.HasSuffix(lowerImg, ".png") && strings.HasPrefix(img, "image")) {
			// We usually only want the *first* one.
			if imgUrl == "" {
				imgUrl = "https://apod.nasa.gov/apod/" + img
			} else { // nolint
				// 2nd link to image which probably is wrong
			}
		}
	})

	description := ""

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
		return &ApodPost{
			Text:        title,
			ImageURL:    imgUrl,
			Description: description,
		}, nil
	}
	return nil, fmt.Errorf("No title or url")
}

func doYMD(y, m, d int) *ApodPost {
	if y > 2000 {
		// Perlism
		y -= 2000
	}

	url := fmt.Sprintf("https://apod.nasa.gov/apod/ap%02d%02d%02d.html", y, m, d)
	resp, err := http.Get(url)
	if err != nil || resp == nil || resp.Body == nil {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	resp.Body.Close()
	myApod, err := ParseApod(string(body))
	if err != nil {
		log.Println(err)
		return nil
	}
	return myApod
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
			MessagePlatforms := GetMessagePlatforms()
			for _, v := range MessagePlatforms {
				v.SendWithOptions(text, bothandler.SendOptions{Silent: true})
			}
		}

		log.Println("Sleeping 5 minutes")
		time.Sleep(5 * time.Minute)
	}
}
