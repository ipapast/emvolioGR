# emvolioGR

![github check](https://github.com/ipapast/emvolioGR/actions/workflows/go.yml/badge.svg) [![CodeQL](https://github.com/ipapast/emvolioGR/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/ipapast/emvolioGR/actions/workflows/codeql-analysis.yml) [![GitHub license](https://img.shields.io/github/license/Naereen/StrapDown.js.svg)](https://github.com/Naereen/StrapDown.js/blob/master/LICENSE)


_Work in progress._

This is a Go version of a covid vaccination tracker/tweetbot in Greece, tweeting here: https://twitter.com/emvolioGR.

This app gets data from the [Data Gov GR website](https://data.gov.gr/datasets/mdg_emvolio/), formats them and shares a daily percantage in Twitter upon trigger (to do: add a cron job to post automatically).

Results before formatting are saved in the `data/vaccinations_regions.json` file too.
# Run locally

- `go build`
- ```go run ./.```
