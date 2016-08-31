# gopad: a personal markdown notepad [![Build Status](https://travis-ci.org/adriamb/gopad.svg?branch=master)](https://travis-ci.org/adriamb/gopad) [![Go Report Card](https://goreportcard.com/badge/github.com/amassanet/gopad)](https://goreportcard.com/report/github.com/amassanet/gopad)

[![](README.md.files/gopad.jpg)]()


GoPad is a web application, with the following features:

- **Markdown** support with github flavour
- **Realtime render**: All data written in the markdown are automatically rendered, with support for
  - **Graphviz** (see http://www.graphviz.org/)
  - **Goat**  (see https://github.com/blampe/goat)
  - **Js Sequence diagrams** (see https://bramp.github.io/js-sequence-diagrams/)
  - **Js Flowchart** (see http://flowchart.js.org/)
- **Attachments** support, just drag&drop
- **No database** , all data are written into the filesystem with .md and .json files. If you want to backup the data, just copy the folder, or create a git repo for it.
- **No data is overwritten**, attachments cannot be overwritten, page changes are versioned
- **Google OAuth2** support, if you want you can push your blog in a public space and log into with your google account.

## Installation

Install the gopad with

`go get https://github.com/adriamb/gopad`

create a minimal configuration file `$HOME/.gopad.yaml` with the following content:

```yaml
port: 8080
```

run the gopad 

`gopad`

go to a browser `http://localhost:8080`

## Fast test with docker

If you want to test the application, clone the repo

`git clone https://github.com/adriamb/gopad`

go to the docker folder

`cd docker`

run docker-compose

`docker-compose up`

go to a browser `http://localhost:8088`

all files generated will be kept in the `docker/data`, take a look

## Configuration file

Gopad uses the following configuration file, by default located in $HOME/.gopad.yaml but you can specify it with the `--config` command line parameter.

```yaml
port: <server port, e.g. 8080>
datadir: <where data is stored, by default is $HOME/.gopad>
tmpdir: <temp directory, by default is /tmp/gopad/tmp>
cachedir: <cache directory, by default is /tmp/gopad/tmp>
auth:
    type: <'none' or 'google'>
    googleclientid: <google oauth2 clientid>
    allowedemails:
        - <email of allowed oauth client>
        - <email of allowed oauth client>
```

## Dropping files

To add files, just drop the file into the markdown edit box. The file will be added to the
filesystem and a link will be created inside the markdown text.

You can access the file with the created link or using the Files menu item.



