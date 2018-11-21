# Baseball Scorecard Generator

Uses play-by-play data files from [retrosheet](http://www.retrosheet.org) to generate json data that will eventually display a baseball scorecard.

This project currently ships with the 1986 New York Mets home games event file to test the data conversion on a smaller scale. The steps below will run the conversion for this data and produce a collection of transformed json files for each game under the `data/out` directory of the project.

## Getting Started
```
# cd into your go workspace
git clone https://github.com/bricemason/go-baseball-scorecard.git
cd go-baseball-scorecard
go run *.go
```

The information used here was obtained free of
charge from and is copyrighted by Retrosheet.  Interested
parties may contact Retrosheet at 20 Sunset Rd.,
Newark, DE 19711.