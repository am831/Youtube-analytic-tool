package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	var id string
	var year int
	var input rune
	var quit bool = true
	var cpm_under8mins float32 = 2
	var cpm_over8mins float32 = 5
	var percent1, percent2 float64

	err := godotenv.Load(".env")
	handleError(err, "Error loading .env file")
	api_key := os.Getenv("api_key")

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(api_key))
	handleError(err, "Error creating YouTube client")

	fmt.Println("Enter channel id: ")
	fmt.Scanln(&id)

	channel := getChannelInfo(service, id)
	videoIDs := getAllVideos(service, channel.playlistId,
		int(channel.videoCount))
	videoInfo := getVideoInfo(service, videoIDs)
	channel.firstPublished = videoIDs[len(videoIDs)-1]
	getWatchHoursandEngagement(&channel, videoIDs, videoInfo, 100, 25)
	fmt.Printf("This channel's name is %s. It has %d subscribers, %d total "+
		"views, %s total watch hours, and has uploaded a total of %d videos.",
		channel.title, channel.subs, channel.views,
		strconv.FormatFloat(channel.totalWatchHours, 'f', 2, 64),
		channel.videoCount)
	fmt.Println()
	fmt.Println()
	printMenu()

	for quit {
		fmt.Scanf("%c\n", &input)
		if input == '1' {
			if channel.totalWatchHours < 4000 {
				fmt.Println("This channel has not been monetized yet. " +
					"Lifetme earnings are 0.")
			} else {
				if channel.totalEarnings == 0 {
					getFirstMonetizedVid(videoIDs, videoInfo, &channel)
					getLifetimeEarnings(videoIDs, videoInfo, cpm_under8mins,
						cpm_over8mins, &channel)
				}
				monetized := videoIDs[len(videoIDs)-
					(videoInfo[channel.firstMonetized].uploadNum)-1]
				fmt.Printf("Since posting their first monetized video on %s, "+
					"this channel has earned ",
					videoInfo[monetized].published[0:10])
				fmt.Printf("%s ", strconv.FormatFloat(channel.totalEarnings, 'f',
					2, 64))
				fmt.Print("in USD" + "\n")
			}
			fmt.Println()
			printMenu()
		} else if input == '2' {
			if channel.totalWatchHours < 4000 {
				fmt.Println("This channel has not been monetized yet. " +
					"Lifetme earnings are 0")
			} else {
				if channel.totalEarnings == 0 {
					getFirstMonetizedVid(videoIDs, videoInfo, &channel)
					getLifetimeEarnings(videoIDs, videoInfo, cpm_under8mins,
						cpm_over8mins, &channel)
				}
				fmt.Println("Enter video ID (found in the video url):")
				fmt.Scanln(&id)
				_, ok := videoInfo[id]
				if ok {
					fmt.Printf("The video titled %s earned %s",
						videoInfo[id].title,
						strconv.FormatFloat(videoInfo[id].earnings, 'f', 2, 64))
				}
				fmt.Println()
			}
			fmt.Println()
			printMenu()
		} else if input == '3' {
			frequenciesOrdered, tags := getTags(videoIDs, videoInfo)
			for i := 0; i < len(frequenciesOrdered)-1; i++ {
				for _, item := range tags[frequenciesOrdered[i]] {
					fmt.Println("\"", item, "\"", "used ",
						frequenciesOrdered[i], " times")
				}
			}
			fmt.Println()
			printMenu()
		} else if input == '4' {
			if channel.totalWatchHours < 4000 {
				fmt.Println("This channel has not been monetized yet.")
			} else {
				countHours, daysPassed := getFirstMonetizedVid(videoIDs, videoInfo,
					&channel)
				fmt.Printf("The channel first surpassed 4000 watch hours within "+
					"12 months with the video titled %s",
					videoInfo[channel.firstMonetized].title)
				fmt.Println()
				fmt.Printf("This video resulted in %s total watch hours for the "+
					"channel", strconv.FormatFloat(countHours, 'f', 2, 64))
				fmt.Println()
				fmt.Printf("This video was posted on %s, %d days after their "+
					"first post",
					videoInfo[channel.firstMonetized].published[0:10], daysPassed)
				fmt.Println()
				fmt.Printf("It took %d videos for this channel to meet this "+
					"criteria for monetization",
					videoInfo[channel.firstMonetized].uploadNum)
				fmt.Println()
			}
			fmt.Println()
			printMenu()
		} else if input == '5' {
			if channel.totalWatchHours < 4000 {
				fmt.Println("This channel has not been monetized yet. " +
					"Lifetme earnings are 0.")
			} else {
				if channel.totalEarnings == 0 {
					getFirstMonetizedVid(videoIDs, videoInfo, &channel)
					getLifetimeEarnings(videoIDs, videoInfo, cpm_under8mins,
						cpm_over8mins, &channel)
				}
				fmt.Println("Show graph for which year? Enter in YYYY format:")
				fmt.Scanln(&year)
				createEarnignsGraph(videoIDs, videoInfo, year, channel)
			}
			fmt.Println()
			printMenu()
		} else if input == '6' {
			fmt.Printf("This channel's name is %s. It has %d subscribers, "+
				"%d total views, %s total watch hours, and has uploaded a total "+
				"of %d videos.", channel.title, channel.subs, channel.views,
				strconv.FormatFloat(channel.totalWatchHours, 'f', 2, 64),
				channel.videoCount)
			fmt.Println()
			fmt.Println()
			printMenu()
		} else if input == '7' {
			fmt.Println("Enter the CPM for videos under 8 mins:")
			fmt.Scanln(&cpm_under8mins)
			fmt.Println("Enter the CPM for videos 8 mins or longer:")
			fmt.Scanln(&cpm_over8mins)
			getLifetimeEarnings(videoIDs, videoInfo, cpm_under8mins,
				cpm_over8mins, &channel)
			printMenu()
		} else if input == '8' {
			fmt.Println("Enter percent of the video that is watched by " +
				"people who interact with the video:")
			fmt.Scanln(&percent1)
			fmt.Println("Enter percent of the video that is watched by " +
				"people who did not interact with the video:")
			fmt.Scanln(&percent2)
			getWatchHoursandEngagement(&channel, videoIDs, videoInfo, percent1,
				percent2)
			fmt.Println()
			printMenu()
		} else if input == '9' {
			printInfo()
			printMenu()
		} else if input == '0' {
			quit = false
		}
	}
}
