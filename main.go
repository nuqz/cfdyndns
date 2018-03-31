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
	api      *cloudflare.API
	cfg      map[string][]string
	err      error
	quit     = make(chan os.Signal, 1)
	cfgPath  = flag.String("config", "config.toml", "Path to config.toml file")
	interval = flag.Duration("interval", time.Minute, "Time interval between updates")
)

func init() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.Ltime)

	_, err = toml.DecodeFile(*cfgPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	api, err = cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		log.Fatal(err)
	}

	signal.Notify(quit, syscall.SIGINT)
}

func main() {
	t := time.Tick(*interval)
eternity:
	for {
		select {
		case <-t:
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
					found := false
					for _, rec := range dns {
						log.Println(rec.Name)
						if rec.Name == recn {
							rec.Content = ip
							if err = api.UpdateDNSRecord(id, rec.ID, rec); err != nil {
								log.Println(err)
							}
							log.Printf("Zone: %s\tRecord: %s\tIP: %s",
								zone, rec.Name, rec.Content)
							found = true
							break
						}
					}
					if !found {
						if _, err := api.CreateDNSRecord(id, cloudflare.DNSRecord{
							Type:    "A",
							Name:    recn,
							Content: ip,
						}); err != nil {
							log.Println(err)
						}
					}
				}

			}
		case <-quit:
			break eternity
		}
	}
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
