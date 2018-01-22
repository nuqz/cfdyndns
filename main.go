package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/cloudflare/cloudflare-go"
	"golang.org/x/net/html"
)

var (
	cfgPath  = flag.String("config", "config.toml", "Path to config.toml file")
	interval = flag.Duration("interval", time.Minute, "Time interval between updates")
	quit     = make(chan os.Signal)
	cfg      map[string][]string
	api      *cloudflare.API
	err      error
	ticker   *time.Ticker
)

func init() {
	flag.Parse()

	_, err = toml.DecodeFile(*cfgPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	api, err = cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		log.Fatal(err)
	}

	ticker = time.NewTicker(*interval)

	signal.Notify(quit, syscall.SIGINT)

	log.SetPrefix("[CF DDNS] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	go func() {
	eternity:
		for {
			select {
			case <-ticker.C:
				ip, err := getExternalIP()
				if err != nil {
					log.Println(err)
					continue
				}

				for zone, records := range cfg {
					id, err := api.ZoneIDByName(zone)
					if err != nil {
						log.Println(err)
						continue
					}

					dns, err := api.DNSRecords(id, cloudflare.DNSRecord{})
					if err != nil {
						log.Println(err)
						continue
					}

					for _, recn := range records {
						for _, rec := range dns {
							if rec.Name == recn {
								rec.Content = ip
								if err = api.UpdateDNSRecord(id, rec.ID, rec); err != nil {
									log.Println(err)
									continue
								}
								log.Printf("Zone: %s\tRecord: %s\tIP: %s",
									zone, rec.Name, rec.Content)
							}
						}
					}

				}
			case <-quit:
				break eternity
			}
		}
	}()
	<-quit
	ticker.Stop()
}

func getExternalIP() (string, error) {
	resp, err := http.Get("http://checkip.dyndns.org/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Unable to connect to checkip.dyndns.org")
	}

	dec := html.NewTokenizer(resp.Body)
	var t []byte
	for string(t) != "body" {
		t, _ = dec.TagName()
		dec.Next()
	}

	return string(dec.Text()[20:]), nil
}
