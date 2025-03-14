package main

import (
	"github.com/robfig/cron"
	"github.com/youngking/gin-blog/models"
	"log"
)

func main() {
	log.Printf("Startng...")
	c := cron.New()
	c.AddFunc("* * * * *", func() {
		log.Println("Run models.CleanAllTag...")
		models.CleanAllTag()
	})

	c.AddFunc("* * * * *", func() {
		log.Println("Run models.CleanAllArticle...")
		models.CleanAllArticle()
	})

	c.Start()

	select {}
}
