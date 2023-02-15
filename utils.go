package main

import (
	"fmt"
	"log"
)

func printMenu() {
	fmt.Println("Pick an option from below:")
	fmt.Println("(1) Show lifetime earnings")
	fmt.Println("(2) Show earnings for a specific video")
	fmt.Println("(3) Show tags and how often they are used")
	fmt.Println("(4) Show how many videos it may have taken for this " +
		"channel to be monetized")
	fmt.Println("(5) Show earnings by month on a bar chart")
	fmt.Println("(6) Show basic channel info")
	fmt.Println("(7) Enter custom CPM rates")
	fmt.Println("(8) Enter custom engagement rates")
	fmt.Println("(9) Show information about each option")
	fmt.Println("(0) Quit")
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func printInfo() {
	fmt.Println("(1) Lifetime earnings are calculated using views, CPM rate, " +
		"and duration. By default, the CPM rate for videos under 8 minutes " +
		"is 2 and the default for videos 8 minutes or longer is 5. See (6) " +
		"for more details on CPM rates. All videos posted before the first " +
		"monetized video do not contribute to earnings.")
	fmt.Println()
	fmt.Println("(2) See (1) for details on how earnings are calculated. " +
		"This option shows the earnings for only the requested video.")
	fmt.Println()
	fmt.Println("(3) This option analyzes every video posted on the channel " +
		"and counts the occurence of every tag.")
	fmt.Println()
	fmt.Println("(4) A channel must reach 4000 watch hours within the last " +
		"12 months to be eligible for monetization. From the first video " +
		"posted, earnings and the days passed between each post are totaled " +
		"up. If 365 days is surpassed, then the earliest post within the 365 " +
		"day window has its earnings removed from the total watch hours, and " +
		"the days passed between it and the next are subtracted from the " +
		"total days. Once 4000 watch hours within 365 days is reached, the " +
		"next video posted is the first monetized video. Watch hours are " +
		"calculated based on duration, views and engagement. See (7) for " +
		"more details on engagement rate. Many details are inaccessible, " +
		"such as when the video accumulated the views, what the exact watch " +
		"hours are, and how many subscribers the channel had at the time.")
	fmt.Println()
	fmt.Println("(5) See (1) for details on how earnings are calculated. " +
		"The earnings for a month are determined by how much videos posted " +
		"in that month earned. If a video generated $2000 and was posted in " +
		"may, then that is considered revenue for may. It is unknown when " +
		"the views were accumulated, they could have been accumulated " +
		"outside of the month the video was posted. The purpose of this " +
		"option is to give an idea of which months might be good for posting " +
		"in this channel's niche.")
	fmt.Println()
	fmt.Println("(6) Shows channel name, total subscribers, total views, " +
		"total watch hours, and number of uploaded videos.")
	fmt.Println()
	fmt.Println("(7) The CPM rate is the dollars earned for every 1000 " +
		"views. The views are divided by 1000 then mulitplied by the CPM. " +
		"Default CPM is 2 for videos under 8 min, and 5 for videos 8 min or " +
		"more. Shorter videos may have lower CPM, because a video can't " +
		"contain midroll ads unless they are 8 minutes or longer. It is " +
		"assumed that a video 8 minutes or longer contains midroll ads. This " +
		"means that for every video under 8 min, $2 is earned for every " +
		"1000 views. CPM rates can be customized with options (6).")
	fmt.Println()
	fmt.Println("(8) Engagement is considered to be a like, dislike, or " +
		"a comment. All people who engaged with a video are assumed to have " +
		"watched 100% of the video. All other people who viewed but didn't " +
		"engage are assumed to have watched 25% of the video. Enagement is " +
		"used to calculate watch hours for videos. It affects the estimate " +
		"for when the channel became monetized and thus the earnigns. These " +
		"percentages can be customized with option (7).")
	fmt.Println()
}
