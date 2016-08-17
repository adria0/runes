//go:generate rice embed-go

package web

import (
	"github.com/adriamb/gopad/server/instance"
	"github.com/adriamb/gopad/store"
	"github.com/adriamb/gopad/store/model"
	"github.com/adriamb/gopad/web/render"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	versionDateTimeDisplay = "2006-01-02 15:04:05"
)

type dtoMarkdownRender struct {
	Markdown string `json:"markdown" binding:"required"`
}

// Version of an entry
type Version struct {
	Description string
	URL         string
}

func doPOSTUpload(c *gin.Context) {

	ws := normalize(c.Param("ws"))
	id := normalize(c.Param("id"))

	file, fileHeader, err := c.Request.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := normalize(fileHeader.Filename)
	err = instance.Srv.Store.Entry.StoreFile(ws, id, filename, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	path := "/w/" + ws + "/e/" + id + "/f/" + filename
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

func doPOSTRender(c *gin.Context) {
	var json dtoMarkdownRender
	if c.BindJSON(&json) == nil {
		html := string(render.Render(json.Markdown))
		c.JSON(http.StatusOK, gin.H{"html": html})
	}
}

func doGETEntryEdit(c *gin.Context) {

	ws := normalize(c.Param("ws"))
	id := normalize(c.Param("id"))

	var entry *model.Entry
	var err error
	var editable = true
	versions := []Version{}

	if id != "new" {

		if strings.Contains(id, ".") {

			split := strings.SplitN(id, ".", 2)
			id = split[0]
			version := split[1]
			entry, err = instance.Srv.Store.Entry.GetEntry(ws, id, version)
			editable = false

		} else {

			entry, err = instance.Srv.Store.Entry.GetEntry(ws, id, "")

		}

		if err == nil {

			versions = append(versions, Version{
				Description: "Last",
				URL:         "/w/" + ws + "/e/" + id + "/edit",
			})

			var versionids []string
			versionids, err = instance.Srv.Store.Entry.GetEntryVersions(ws, id)
			sort.Sort(sort.Reverse(sort.StringSlice(versionids)))

			if len(versionids) > 15 {
				versionids = versionids[:15]
			}

			for _, versionid := range versionids {
				t, timeerr := time.Parse(store.DateTimeFormat, versionid)

				if timeerr == nil {
					versions = append(versions, Version{
						Description: t.Format(versionDateTimeDisplay),
						URL:         "/w/" + ws + "/e/" + id + "." + versionid + "/edit",
					})
				}
			}
		}

		if err != nil {
			dumpError(c, err)
			return
		}

	} else {

		entry = &model.Entry{
			ID: instance.Srv.Store.Entry.NewID(),
		}

	}

	c.HTML(http.StatusOK, "entry.tmpl", gin.H{
		"ws":       ws,
		"entry":    entry,
		"editable": editable,
		"versions": versions,
	})

}

func doPOSTEntryDelete(c *gin.Context) {

	ws := normalize(c.Param("ws"))
	id := normalize(c.Param("id"))

	err := instance.Srv.Store.Entry.DeleteEntry(ws, id)

	if err != nil {
		dumpError(c, err)
	}

	c.Redirect(http.StatusSeeOther, "/w/"+ws)

}

func doPOSTEntry(c *gin.Context) {

	ws := normalize(c.Param("ws"))
	id := normalize(c.Param("id"))

	entry := model.Entry{
		Workspace: ws,
		ID:        id,
		Markdown:  c.DefaultPostForm("Markdown", "undefined"),
	}

	err := instance.Srv.Store.Entry.StoreEntry(&entry)
	if err != nil {
		dumpError(c, err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/w/"+ws)
	return

}
