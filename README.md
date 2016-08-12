# gopad: a personal markdown notepad [![Build Status](https://travis-ci.org/adriamb/gopad.svg?branch=master)](https://travis-ci.org/adriamb/gopad) [![Go Report Card](https://goreportcard.com/badge/github.com/amassanet/gopad)](https://goreportcard.com/report/github.com/amassanet/gopad)

[![](README.md.files/gopad.jpg)]()


GoPad is a web application, with the following features:

- **Markdown** support with github flavour
- **Realtime render**: All data written in the markdown are automatically rendered
- **Graphviz (dot)** support
- **Sequence diagrams** support via UMLlet
- **Attachments** support, just drag&drop
- **No database** , all data are written into the filesystem with .md and .json files. If you want to backup the data, just copy the folder, or create a git repo for it.
- **No data is overwritten**, attachments cannot be overwritten, page changes are versioned
- **Google OAuth2** support, if you want you can push your blog in a public space and log into with your google account.

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

Gopad uses the following configuration file, by default located in $HOME/.gopad but you can specify it with the `--config` command line parameter.

```yaml
port: <server port, e.g. 8080>
prefix: <server prefix, left it blank>
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

## Support for graphviz

Gopad uses graphviz using the external command `dot`, so it should be installed in your system.

The syntax for dot is (note that the `digraph {}` element are automatically added)

        ```dot
        func getTrue() bool {
            return true
        }
        ```

## Support for UMLet

Gopad uses umlet as an external command, and, since it's a java application needs java jre to be installed in order to work.

At this moment there's only support for sequence diagrams in the following form:

      ```umlet:sequence
      obj=Usr~usr
      obj=App~app
      
      usr->app
      app->app +:ACT=CreateKeyPair
      app->usr:show PIN=\nhash pbkACT
      ```



