//go:generate rice embed-go

package web

import (
	"net/http"
	"os"
	"strings"
    "sort"
	"time"
    "github.com/adriamb/gopad/store"
    "github.com/adriamb/gopad/model"
	"github.com/adriamb/gopad/server"
	"github.com/adriamb/gopad/web/render"
	"github.com/gin-gonic/gin"
)

const (
    versionDateTimeDisplay = "2006-01-02 15:04:05"
)

type dtoMarkdownRender struct {
	Markdown string `json:"markdown" binding:"required"`
}

func doPOSTUpload(c *gin.Context) {

	id := c.Param("id")

	if existsStaticMd(id) {
		return
	}

	file, fileHeader, err := c.Request.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename, err := server.Srv.Store.File.Write(fileHeader.Filename, id, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path := server.Srv.Config.Prefix + "/files/" + filename
	split := strings.Split(filename, ".")
	ext := split[len(split)-1]
	ico := ""

	if ext == "png" || ext == "jpg" || ext == "jpeg" || ext == "gif" {

	} else if _, err := os.Stat("web/static/ico/" + ext + ".png"); err == nil {
		ico = "/static/ico/" + ext + ".png"
	} else {
		ico = "/static/ico/default.png"
	}

	c.JSON(http.StatusOK, gin.H{"name": fileHeader.Filename, "path": path, "ico": ico})
}

func doPOSTMarkdown(c *gin.Context) {
	var json dtoMarkdownRender
	if c.BindJSON(&json) == nil {
		html := string(render.Render(json.Markdown, server.Srv.Dict))
		c.JSON(http.StatusOK, gin.H{"html": html})
	}
}

type Version struct {
    Description string
    URL         string
}

func doGETEntry(c *gin.Context) {


    id := c.Param("id")

	var entry *model.Entry
	var err error
	var editable = true
    versions := []Version{}

	if id != "new" {

		if existsStaticMd(id) {

			entry, err = getStaticMdEntry(id)
			editable = false

		} else {

            if strings.Contains(id,".") {

                split := strings.SplitN(id,".",2)
                id = split[0]
                version := split[1]
			    entry, err = server.Srv.Store.Entry.GetVersion(id,version)
                editable = false

            } else {

                entry, err = server.Srv.Store.Entry.Get(id)

            }

            if err == nil {

                versions = append ( versions , Version {
                    Description: "Last",
                    URL: "/entries/"+id+"/edit",
                })

                var versionids []string
                versionids, err = server.Srv.Store.Entry.GetVersions(id)
                sort.Sort(sort.Reverse(sort.StringSlice(versionids)))

                if len(versionids) > 15 {
                    versionids = versionids[:15]
                }

                for _, versionid:= range versionids {
                    t, timeerr := time.Parse(store.DateTimeFormat,versionid)

                    if timeerr == nil {
                        versions = append ( versions , Version {
                            Description: t.Format( versionDateTimeDisplay ),
                            URL: "/entries/"+id+"."+versionid+"/edit",
                        } )
                    }
                }
            }
        }

		if err != nil {
			dumpError(c, err)
			return
		}

	} else {

        entry = &model.Entry{
			ID: server.Srv.Store.Entry.NewID(),
		}

    }

	c.HTML(http.StatusOK, "entry.tmpl", gin.H{
		"prefix":   server.Srv.Config.Prefix,
		"entry":    entry,
		"editable": editable,
        "versions": versions,
	})

}

func doPOSTEntry(c *gin.Context) {

	id := c.Param("id")

	if existsStaticMd(id) {
		return
	}

	entry := model.Entry{
		ID:       id,
		Title:    c.DefaultPostForm("Title", "undefined"),
		Markdown: c.DefaultPostForm("Markdown", "undefined"),
	}

	err := server.Srv.Store.Entry.Store(&entry)
	if err != nil {
		err = server.Srv.Dict.Rebuild()
	}
	if err != nil {
		dumpError(c, err)
		return
	}
	c.Redirect(http.StatusSeeOther, server.Srv.Config.Prefix+"/entries")
	return

}
