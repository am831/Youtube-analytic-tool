# Youtube Analytic Tool
A command line app that allows the user to analyze youtube channels with the youtube data API

## The app provides the following information
1. Lifetime earnings
2. Earnings for specific videos
3. How often each tag is used
4. How many videos it may have taken for the channel to be monetized
5. Earnings by month on a bar chart

## Getting started
* Download the source code
* Get a free API key from google developer console for the youtube data API [here](https://console.cloud.google.com/marketplace/product/google/youtube.googleapis.com?q=search&referrer=search)
* Add a .env file to the app directory that contains a variable called api_key, then set it equal to your API key
* Get the channel ID for a channel you want to analyze
* Run or build all .go files

## How the information is calculated 
* Earnings for a video are calculated using views, CPM rate, and duration. By default, the CPM rate for videos under 8 minutes is 2 and the default for videos 8 minutes or longer is 5. The CPM rate is the dollars earned for every 1000 views. Shorter videos may have lower CPM, because a video can't contain midroll ads unless they are 8 minutes or longer. It is assumed that a video 8 minutes or longer contains midroll ads. All videos posted before the first monetized video do not contribute to earnings.
* CPM rates can be customized from the main menu
* How long it took for a channel to be monetized is calculated based on watch hours. A channel must reach 4000 watch hours within the last 12 months to be eligible for monetization. From the first video posted, earnings and the days passed between each post are totaled up. If 365 days is surpassed, then the earliest post within the 365 day window has its earnings removed from the total watch hours, and the days passed between it and the next are subtracted from the total days. Once 4000 watch hours within 365 days is reached, the next video posted is the first monetized video. 
* Watch hours are calculated based on duration, views and engagement. Engagement is considered to be a like, dislike, or a comment. All people who engaged with a video are assumed to have watched 100% of the video. All other people who viewed but didn't engage are assumed to have watched 25% of the video. 
* % of the video watched depending on engagement can be customized from the main menu
* Many details are inaccessible, such as when the video accumulated the views, what the exact watch hours are, and how many subscribers the channel had at the time. Money earned and date of monetization are estimates.

## Technologies
This project was created with go 1.19.5
