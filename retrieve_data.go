package main

import (
	"strconv"
	"time"

	"google.golang.org/api/youtube/v3"
)

type ChannelDetails struct {
	title           string
	subs            uint64
	views           uint64
	videoCount      uint64
	playlistId      string
	totalEarnings   float64
	firstPublished  string
	firstMonetized  string
	totalWatchHours float64
}

type VideoDetails struct {
	views        uint64
	title        string
	categoryID   int
	published    string
	durationSec  float64
	interactions uint64
	tags         []string
	earnings     float64
	watchHours   float64
	engagement   float64
	uploadNum    int
	year         int
	month        int
}

func getChannelInfo(service *youtube.Service, ID string) ChannelDetails {
	// Takes channel ID as argument and gets the channel's name, and total
	// number of subscribers, views and uploads playlist id
	var details ChannelDetails
	part := []string{"snippet", "contentDetails", "statistics"}

	call := service.Channels.List(part).Id(ID)
	response, err := call.Do()
	handleError(err, "")
	details.title = response.Items[0].Snippet.Title
	details.subs = response.Items[0].Statistics.SubscriberCount
	details.views = response.Items[0].Statistics.ViewCount
	details.videoCount = response.Items[0].Statistics.VideoCount
	details.playlistId =
		response.Items[0].ContentDetails.RelatedPlaylists.Uploads
	return details
}

func getAllVideos(service *youtube.Service, playlistID string, videoCount int) []string {
	// Takes playlist ID as argument, and gets all video ID's from the playlist
	var videoIDs []string
	part := []string{"snippet", "contentDetails"}

	call :=
		service.PlaylistItems.List(part).PlaylistId(playlistID).MaxResults(50)
	response, err := call.Do()
	handleError(err, "")
	for _, item := range response.Items {
		videoIDs = append(videoIDs, item.ContentDetails.VideoId)
	}
	nextPageToken := response.NextPageToken
	for len(videoIDs) != videoCount {
		call :=
			service.PlaylistItems.List(part).PlaylistId(playlistID).MaxResults(50)
		call = call.PageToken(nextPageToken)
		response, err := call.Do()
		handleError(err, "")
		for _, item := range response.Items {
			videoIDs = append(videoIDs, item.ContentDetails.VideoId)
		}
		nextPageToken = response.NextPageToken
	}
	return videoIDs
}

func getVideoInfo(service *youtube.Service, videoIDs []string) map[string]*VideoDetails {
	// Gets information for each video the channel has published
	part := []string{"snippet", "contentDetails", "statistics"}
	videoInfo := make(map[string]*VideoDetails)
	uploads := len(videoIDs)
	const (
		layout1 = "2006-01-02"
	)

	for i := 0; i < len(videoIDs); i++ {
		call := service.Videos.List(part).Id(videoIDs[i])
		response, err := call.Do()
		handleError(err, "")
		entry := new(VideoDetails)
		entry.views = response.Items[0].Statistics.ViewCount
		entry.title = response.Items[0].Snippet.Title
		entry.categoryID, err =
			strconv.Atoi(response.Items[0].Snippet.CategoryId)
		handleError(err, "Error converting string to int")
		entry.published = response.Items[0].Snippet.PublishedAt
		entry.durationSec =
			getDuration(response.Items[0].ContentDetails.Duration[2:])
		entry.interactions = (response.Items[0].Statistics.DislikeCount +
			response.Items[0].Statistics.LikeCount +
			response.Items[0].Statistics.CommentCount)
		entry.tags = response.Items[0].Snippet.Tags
		entry.uploadNum = uploads
		date, err := time.Parse(layout1, entry.published[0:10])
		handleError(err, "Error converting string to time")
		entry.year = date.Year()
		entry.month = int(date.Month())
		videoInfo[videoIDs[i]] = entry
		uploads -= 1
	}
	return videoInfo
}
