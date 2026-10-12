package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"server/internal/constant"
	"server/internal/model"
	"server/internal/route"
	"server/internal/safemap"
	"server/internal/service"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// get path
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// get stat data
	tierLists := []model.TierList{}
	for _, tier := range constant.AllTier {
		tierList := model.TierList{}
		b, err := os.ReadFile(wd + "/asset/data/tier_" + strings.ToLower(string(tier)) + ".json")
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(b, &tierList); err != nil {
			panic(err)
		}
		tierLists = append(tierLists, tierList)
	}
	// get global counter
	b, err := os.ReadFile(wd + "/asset/data/global_counter.txt")
	if err != nil {
		panic(err)
	}
	counterStr := string(b)
	counter, err := strconv.ParseInt(counterStr, 10, 64)
	if err != nil {
		panic(err)
	}
	globalCounter := atomic.Int64{}
	globalCounter.Add(counter)
	safeMap := safemap.TierSafeMap{}

	for _, v := range tierLists {
		safeMap.SafeStore(&v)
	}

	svc := service.NewTierMakerService()
	t := route.NewTierMakerRoute(svc, &safeMap, &globalCounter)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Post("/tier-maker", t.MakeTier)
	app.Get("/tier-list", t.GetTiers)

	// Start Fiber in a goroutine.
	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Printf("Fiber server stopped: %v", err)
		}
	}()

	// Wait for SIGINT or SIGTERM.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop
	defer cancel()

	// store global counter
	a := strconv.FormatInt(globalCounter.Load(), 10)
	os.WriteFile(wd+"/asset/data/global_counter.txt", []byte(a), 0644)
	log.Printf("Received signal (%s), shutting down...", sig)

	// store stat data
	datas := safeMap.ToList()
	for _, data := range datas {
		b, err = json.MarshalIndent(&data, "", "  ")
		if err == nil {
			os.WriteFile(wd+"/asset/data/tier_"+strings.ToLower(string(data.Tier))+".json", b, 0644)
		}
	}
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	log.Println("Server shutdown successfully")
}
