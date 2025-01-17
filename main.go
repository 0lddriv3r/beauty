// MIT License

// Copyright (c) 2017 FLYING

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/yang-f/beauty/controllers"

	"github.com/yang-f/beauty/router"
	"github.com/yang-f/beauty/settings"
	"github.com/yang-f/beauty/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app      = kingpin.New("beauty", "A command-line tools of beauty.")
	demo     = app.Command("demo", "Demo of web server.")
	generate = app.Command("generate", "Generate a new app.")
	name     = generate.Arg("name", "AppName for app.").Required().String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case generate.FullCommand():
		goPathEnv := os.Getenv("GOPATH")
		goPaths := strings.Split(goPathEnv, ":")
		targetGoPath := ""
		targetZipPath := ""
		tailZipPath := "github.com/yang-f/beauty@v0.0.6/etc/demo.zip"
		for _, goPath := range goPaths {
			tempZipPath := fmt.Sprintf("%s/pkg/mod/%s", goPath, tailZipPath)
			_, err := os.Stat(tempZipPath)
			if err == nil {
				targetGoPath = goPath
				targetZipPath = tempZipPath
				break
			}
		}
		if targetZipPath == "" {
			log.Fatal("Fatal: can not find " + tailZipPath + " in your gopath.")
		}
		appPath := fmt.Sprintf("%s/src/%s", targetGoPath, *name)
		dst := fmt.Sprintf("%s.zip", appPath)
		_, err := utils.CopyFile(dst, targetZipPath)
		if err != nil {
			log.Fatal("Err_Fatal: " + err.Error())
		}
		utils.Unzip(dst, appPath)
		os.RemoveAll(dst)
		helper := utils.ReplaceHelper{
			Root:    appPath,
			OldText: "{appName}",
			NewText: *name,
		}
		helper.DoWork()
		log.Printf("Generate %s success.", *name)
	case demo.FullCommand():
		log.Printf("Start server on port %s", settings.Listen)
		r := router.New()
		r.GET("/", controllers.Config().ContentJSON())
		r.GET("/demo1", controllers.Config().ContentJSON().Verify())
		log.Fatal(http.ListenAndServe(settings.Listen, r))
	}
}
