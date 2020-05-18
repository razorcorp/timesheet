# Timesheet

A commandline tool written to help developer with logging work time to Atlassian Jira based projects.

This tool allows developers to add worklog to a specific ticket without having to browse and find issues
given the user knows the issue reference

#### This documentation is written for Unix like systems, please translate for other Operating Systems

## Table of Contents

1. [Usage](#usage)
2. [Requirements](#requirements)
3. [Environment configuration](#environment-configuration)
4. [Installation](#installation)
5. [Build](#build)
6. [License](#license)
7. [Authors](#authors)

## Usage
```
timesheet (-r -t [-d] [-m]] [[-h] [-e] [-d]) ([-remaining] [-history])
  -d string
        Default 2020-05-18. The date on which the worklog effort was started in full date (YYYY-MM-DD) or relative date (-N) format. eg: 2006-01-02 or -1.
  -e string
        HELP: Base64 encode the given credentials. Format: email:token;domain. e.g. example@example.com:abcThisIsFake;xyz.atlassian.net
  -h    HELP: This tool can be used to log time spent on a specific Jira ticket on a project.
  -history
        HELP: Print the timesheet of the day -d is also available to change the week
  -m string
        OPTIONAL: A comment about the worklog
  -month
        HELP: Print timesheet of the current month. -d is also available to change the week
  -r string
        REQUIRED: Jira ticket reference. E.g. DDSP-4
  -remaining
        HELP: Print how many hour can be book for the current day. -d is also available
  -t string
        REQUIRED: The time spent as days (#d), hours (#h), or minutes (#m or #). E.g. 8h
  -v    Print application version
  -week
        HELP: Print timesheet of the current week. -d is also available to change the week
Example:
        timesheet -r DDSP-XXXX -t 8h -m "Jenkins pipeline completed"
        timesheet -r DDSP-XXXX -t 1h -m "Investigated possible solutions" -d 2020-03-05
        timesheet -remaining
        timesheet -remaining -d 2020-03-05
        timesheet -history
        timesheet -history -d -1
```

## Requirements
1. Atlassian account
1. Atlassian personal access token. https://id.atlassian.com/manage/api-tokens

## Environment configuration
Following environment variable must be exported to the system's environment to work.

* Form the email, token, and atlasian-domain in following format. email`:`token`;`atlasian-domain
* Encode the above formed text in Base64.
```bash
$ echo "example@example.com:abcThisIsFake;xyz.atlassian.net" | base64

ZXhhbXBsZUBleGFtcGxlLmNvbTphYmNUaGlzSXNGYWtlO3h5ei5hdGxhc3NpYW4ubmV0Cg==
```
* Export the Base64 encoded value to `TIMESHEET` as an environment variable permanently
```bash
$ export TIMESHEET="ZXhhbXBsZUBleGFtcGxlLmNvbTphYmNUaGlzSXNGYWtlO3h5ei5hdGxhc3NpYW4ubmV0Cg=="
```
_add this to the `.bash_profile` to preserver it_
 
## Installation

1. Download the binary file from the repository's latest release.
```
https://github.com/praveenprem/timesheet/releases
```

2. Give application execution permission.
```bash
$ sudo chmod +x ./timesheet
```

3. Install on the system path.
```bash
install ./timesheet /usr/local/bin/timesheet
```

## Build

1. Clone the repository to your local GO path.
```bash
$ git clone git@github.com:praveenprem/timesheet.git
```

1. Run Makefile to build the application.
```bash
$ make build
```

## License

MIT License

Copyright (c) 2020 Praveen Premaratne

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Authors
   | <div><a href="https://github.com/praveenprem"><img width="200" src="https://avatars3.githubusercontent.com/u/23165760"/><p></p><p>Praveen Premaratne</p></a></div> |
   | :-------: |
