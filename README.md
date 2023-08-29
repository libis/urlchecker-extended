> **Note**
> ## Repository Fork Details
> This repository has been forked from https://github.com/bythepixel/urlchecker for the following reason(s):
>
> 1. Each test failure sends a seperate message to slack/teams which can get messy if there are allot of tests for a single domain. This will be combined into a single message.  
>
> 2. The 'status' parameter in the URLs file is now called 'expected_statuses' and is an array (see example below).
>
> Special thanks to the people at bythepixel for the original GitHub Action that inspired this fork.

# urlchecker

## Description
This GitHub Action reads a JSON file containing the URLs to check, crawls the URLs, and checks the resposne status against an array of acceptable statuses.

Below is an example YAML file for this action.

```yaml
name: Check URLs

on:
  push:
    branches:
      - '*'

env:
  SLACK_WEBHOOK: ${{ secrets.SLACK_WEBHOOK }}

jobs:
  check-urls:
    runs-on: ubuntu-latest
    name: Checks URLs from JSON file
    steps:
      - name: Checkout
        uses: actions/checkout@v2
      - name: Check URLs
        uses: bythepixel/urlchecker@v0.2.0
        with:
          hostname: 'postman-echo.com'
          filename: ./urls.json
```

## JSON File

```json
[
    {
        "url": "/status/200",
        "expected_statuses": [200, 201]
    },
    {
        "url": "/status/200",
        "expected_statuses": [200, 201]
    },
    {
        "url": "/status/200",
        "expected_statuses": [200, 201]
        "regex": "200"
    },
    {
        "path": "/store-sitemap.xml",
        "expected_statuses": [200, 201]
        "xml_sitemap": true
    }
]
```
## Parameters

| Parameter   | Type     | Mandatory | Default | Description                                        |
|-------------|----------|-----------|---------|----------------------------------------------------|
| filename    | string   | Yes       | -       | JSON File with paths                              |
| hostname    | string   | Yes       | -       | Hostname of website                               |
| protocol    | string   | No        | https   | Protocol to use                                   |
| workers     | int      | No        | 5       | Number of concurrent workers                      |
| sleepFlag   | int      | No        | 0       | Number of seconds to sleep between requests      |


View the files in the [json](json) folder to see more examples. See the Golang
[regexp][1] package for additional information on supported regular expressions.

## Note
* The `SLACK_WEBHOOK` environment variable is required, this is the URL to send the message to when something goes wrong. Works with Teams also.

[1]: https://pkg.go.dev/regexp
