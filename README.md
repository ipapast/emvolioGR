# emvolioGR

![github check](https://github.com/ipapast/emvolioGR/actions/workflows/go.yml/badge.svg) [![CodeQL](https://github.com/ipapast/emvolioGR/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/ipapast/emvolioGR/actions/workflows/codeql-analysis.yml) [![GitHub license](https://img.shields.io/github/license/Naereen/StrapDown.js.svg)](https://github.com/Naereen/StrapDown.js/blob/master/LICENSE)


This is a Go version of a covid vaccination tracker/tweetbot in Greece, tweeting here: https://twitter.com/emvolioGR.

This app gets data from the [Data Gov GR website](https://data.gov.gr/datasets/mdg_emvolio/), formats them and shares a daily percantage in Twitter upon trigger (this happens via a Github Action every morning)

# Run locally

- ```go run .```
