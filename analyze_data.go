package main

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func getLifetimeEarnings(videoIDs []string, videoInfo map[string]*VideoDetails,
	cpm_under8mins float32, cpm_over8mins float32, channel *ChannelDetails) {
	// Calculates how much money this channel has earned after reaching 4000
	// watch hours within one year
	var earnings float64
	var totalEarnings float64 = 0
	monetized := videoInfo[channel.firstMonetized].uploadNum + 1

	for i := len(videoIDs) - monetized; i >= 0; i-- {
		if videoInfo[videoIDs[i]].durationSec < 480 {
			earnings = (float64(videoInfo[videoIDs[i]].views) / 1000) *
				float64(cpm_under8mins)
		} else {
			earnings = (float64(videoInfo[videoIDs[i]].views) / 1000) *
				float64(cpm_over8mins)
		}
		videoInfo[videoIDs[i]].earnings = earnings
		totalEarnings += earnings
	}
	channel.totalEarnings = totalEarnings
}

func getDuration(duration string) float64 {
	// Gets the length of a video in seconds
	var totalSeconds float64 = 0

	if strings.Contains(duration, "H") {
		hour, err := strconv.ParseFloat(duration[0:strings.IndexByte(duration,
			'H')],
			64)
		handleError(err, "Error converting min to float64")
		totalSeconds += hour * 3600
		duration = duration[strings.IndexByte(duration, 'H')+1:]
	}
	if strings.Contains(duration, "M") {
		min, err := strconv.ParseFloat(duration[0:strings.IndexByte(duration,
			'M')],
			64)
		handleError(err, "Error converting min to float64")
		totalSeconds += min * 60
		duration = duration[strings.IndexByte(duration, 'M')+1:]
	}
	if strings.Contains(duration, "S") {
		sec, err := strconv.ParseFloat(duration[0:strings.IndexByte(duration,
			'S')], 64)
		handleError(err, "Error converting sec to float64")
		totalSeconds += sec
	}
	return totalSeconds
}

func getWatchHoursandEngagement(videoIDs []string,
	videoInfo map[string]*VideoDetails, percent1 float64, percent2 float64) {
	// Calculates the watch hours for each video based on duration, views, and
	// engagement %. Engagement % is ((likes + dislieks + comments) / views) *
	// 100. It assumes that if x% of people engaged, then x% of people watched
	// percent1% of the video. It assumes the remaining people who viewed but
	// didn't interact watched percent2% of the video. These values can be
	// customized from the main menu.
	for i := len(videoIDs) - 1; i >= 0; i-- {
		videoInfo[videoIDs[i]].engagement =
			(float64(videoInfo[videoIDs[i]].interactions) /
				float64(videoInfo[videoIDs[i]].views)) * 100
		videoInfo[videoIDs[i]].watchHours =
			((float64(videoInfo[videoIDs[i]].interactions) *
				(videoInfo[videoIDs[i]].durationSec * (percent1 / 100))) +
				((float64(videoInfo[videoIDs[i]].views) -
					float64(videoInfo[videoIDs[i]].interactions)) *
					(videoInfo[videoIDs[i]].durationSec * (percent2 / 100)))) /
				3600
	}
}

func getFirstMonetizedVid(videoIDs []string,
	videoInfo map[string]*VideoDetails, channel *ChannelDetails) (float64, int) {
	// Estimates when the channel may have met eligibility criteria for
	// monetization based on watch hours and engagement. A channel must reach
	// 4000 watch hours within the last 12 months to be eligible. If 365 days
	// is surpassed, then the video at the bottom of the queue and the watch
	// hours it contributed are subtrated from the total watch hours, and the
	// days passed between it and the next are subtracted from the total days.
	// Returns total hours and days passed since first post. Engagement rates
	// can be customized from the main menu
	var countDays int = 0
	var countHours float64 = 0
	var queue []string
	const (
		layout1 = "2006-01-02"
	)

	for i := len(videoIDs) - 1; i >= 0; i-- {
		queue = append(queue, videoIDs[i])
		countHours += videoInfo[videoIDs[i]].watchHours
		currDate, err := time.Parse(layout1,
			videoInfo[videoIDs[i]].published[0:10])
		handleError(err, "Error converting string to time")
		if i < len(videoIDs)-1 {
			prevDate, err := time.Parse(layout1,
				videoInfo[videoIDs[i+1]].published[0:10])
			handleError(err, "Error converting string to time")
			diff := currDate.Sub(prevDate)
			daysPassed := int(math.Floor(diff.Hours() / 24))
			countDays += daysPassed
		}
		if countHours >= 4000 && countDays <= 365 {
			channel.firstMonetized = videoIDs[i]
			break
		} else if countDays > 365 {
			for countDays > 365 {
				countHours -= videoInfo[queue[0]].watchHours
				nextDate, err := time.Parse(layout1,
					videoInfo[queue[1]].published[0:10])
				handleError(err, "Error converting string to time")
				prevDate, err := time.Parse(layout1,
					videoInfo[queue[0]].published[0:10])
				handleError(err, "Error converting string to time")
				diff := nextDate.Sub(prevDate)
				daysPassed := int(math.Floor(diff.Hours() / 24))
				countDays -= daysPassed
				queue = queue[1:]
			}
		}
	}
	currDate, err := time.Parse(layout1,
		videoInfo[channel.firstMonetized].published[0:10])
	handleError(err, "Error converting string to time")
	prevDate, err := time.Parse(layout1,
		videoInfo[channel.firstPublished].published[0:10])
	handleError(err, "Error converting string to time")
	diff := currDate.Sub(prevDate)
	daysPassed := int(math.Floor(diff.Hours() / 24))
	return countHours, daysPassed
}

func getTags(videoIDs []string, videoInfo map[string]*VideoDetails) ([]int,
	map[int][]string) {
	// Gets the frequencies of all tags used by the channel and sorts them by
	// frequency
	var frequenciesOrdered []int
	tagsToFrequencies := make(map[string]int)
	frequenciesToTags := make(map[int][]string)

	for i := 0; i < len(videoIDs); i++ {
		for _, item := range videoInfo[videoIDs[i]].tags {
			_, isPresent := tagsToFrequencies[item]
			if isPresent {
				tagsToFrequencies[item] += 1
			} else {
				tagsToFrequencies[item] = 1
			}
		}
	}
	for key, value := range tagsToFrequencies {
		_, isPresent := frequenciesToTags[value]
		if !isPresent {
			frequenciesOrdered = append(frequenciesOrdered, value)
		}
		frequenciesToTags[value] = append(frequenciesToTags[value], key)
	}
	sort.Ints(frequenciesOrdered)
	return frequenciesOrdered, frequenciesToTags
}

func createEarnignsGraph(videoIDs []string, videoInfo map[string]*VideoDetails,
	year int, channel ChannelDetails) {
	// creates a bar chart to show earnings by month for a given year
	earnings := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	err := ui.Init()
	handleError(err, "failed to initialize termui: %v")
	monetized := videoInfo[channel.firstMonetized].uploadNum + 1
	defer ui.Close()

	for i := len(videoIDs) - monetized; i >= 0; i-- {
		if videoInfo[videoIDs[i]].year == year {
			earnings[videoInfo[videoIDs[i]].month-1] +=
				videoInfo[videoIDs[i]].earnings
		} else if videoInfo[videoIDs[i]].year > year {
			break
		}
	}
	for i := 0; i < len(earnings); i++ {
		earnings[i] = math.Floor(earnings[i])
	}
	chart := widgets.NewBarChart()
	chart.Data = earnings
	chart.Labels = []string{"jan", "feb", "mar", "apr", "may", "jun", "jul",
		"aug", "sept", "oct", "nov", "dec"}
	chart.Title = "Bar Chart"
	chart.SetRect(5, 5, 145, 25)
	chart.BarWidth = 10
	chart.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	chart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	chart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorYellow)}

	ui.Render(chart)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
