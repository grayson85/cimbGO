package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/fatih/color"
)

func main() {
	// Set up Chrome options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-logging", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-infobars", true),
	)

	// Create context
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Previous label value
	var prevRate float64

	// Function to fetch and print label content
	fetchAndPrintLabel := func() {
		var labelContent string

		// Run Chrome tasks
		err := chromedp.Run(ctx,
			chromedp.Navigate("https://www.cimbclicks.com.sg/sgd-to-myr"),
			chromedp.WaitVisible(`#rateStr`, chromedp.ByID),
			chromedp.Text(`#rateStr`, &labelContent, chromedp.ByID),
		)

		if err != nil {
			log.Println("Error fetching label:", err)
			return
		}

		// Extract rate value from the label content
		rateStr := strings.TrimPrefix(labelContent, "SGD 1.00 = MYR ")
		currentRate, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			log.Println("Error parsing rate:", err)
			return
		}

		// Get current timestamp
		currentTime := time.Now().Format("2006-01-02 15:04:05")

		// Determine color based on comparison with previous rate
		var colorFunc func(format string, a ...interface{}) string
		if prevRate == 0 {
			colorFunc = color.New(color.FgWhite).SprintfFunc()
		} else if currentRate < prevRate {
			colorFunc = color.New(color.FgRed).SprintfFunc()
		} else if currentRate > prevRate {
			colorFunc = color.New(color.FgGreen).SprintfFunc()
		} else {
			colorFunc = color.New(color.FgWhite).SprintfFunc()
		}

		// Print the label content with timestamp and color
		fmt.Println(colorFunc("%s : Rate : SGD 1.00 = MYR %.4f", currentTime, currentRate))

		// Update previous rate
		prevRate = currentRate
	}

	// Print application information in blue color
	colorFunc := color.New(color.FgBlue).SprintfFunc()
	fmt.Println(colorFunc("============================================================"))
	fmt.Println(colorFunc("CIMB Go - Check SGD to MYR Currency Rate"))
	fmt.Println(colorFunc("============================================================"))
	fmt.Println()
	fmt.Println(colorFunc("=== Version: 1.0 ==="))
	fmt.Println(colorFunc("=== Grayson Lee, July 2024 ==="))
	fmt.Println()
	colorFunc = color.New(color.FgRed).SprintfFunc()
	fmt.Println(colorFunc("***** Press CTRL+C to stop the program *****"))
	fmt.Println()

	// Loop to fetch and print label content every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		fetchAndPrintLabel()

		select {
		case <-ticker.C:
			// Continue fetching the label content
		case <-context.Background().Done():
			// Exit if the context is done
			return
		}
	}
}
