# emvolioGR
[![Golang Version](https://img.shields.io/github/go-mod/go-version/ipapast/emvolioGR)](https://github.com/ipapast/emvolioGR)
[![Go Report Card](https://goreportcard.com/badge/github.com/ipapast/emvolioGr)](https://goreportcard.com/report/github.com/ipapast/emvolioGr)

![SNYK](https://github.com/ipapast/emvolioGR/actions/workflows/snyk-security.yml/badge.svg)
![github check](https://github.com/ipapast/emvolioGR/actions/workflows/go.yml/badge.svg) [![CodeQL](https://github.com/ipapast/emvolioGR/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/ipapast/emvolioGR/actions/workflows/codeql-analysis.yml) 

[![GitHub license](https://img.shields.io/github/license/Naereen/StrapDown.js.svg)](https://github.com/Naereen/StrapDown.js/blob/master/LICENSE)

### Update
The twitter account is now deleted as there is no use for it :) This was an experiment with government data and bots in Twitter. The code can be re-used for similar ideas, but it served its purpose so it's now ceased.

### Old Readme
This is a Go version of a covid vaccination tracker/tweetbot in Greece, tweeting here: https://twitter.com/emvolioGR.

This app gets data from the [Data Gov GR website](https://data.gov.gr/datasets/mdg_emvolio/), formats them and shares a daily percantage in Twitter upon trigger (this happens via a Github Action every morning)

# Run locally

- ```go run .```
