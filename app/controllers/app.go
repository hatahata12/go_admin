package controllers

import (
	"admin/app"
	"github.com/revel/revel"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Scraping() revel.Result {

	out, _ := exec.Command("sh", "-c", "ps -ef | grep batch.php | grep -v grep | head -n 3 | awk '{ print $2 }'").Output()

	pid := string(out)
	pid = strings.TrimRight(pid, "\n")
	log.Printf("pid: %s", pid)

	var count string
	app.DB.QueryRow("SELECT count(item_code) FROM `scraping`").Scan(&count)
	log.Printf("count: %s", count)

	return c.Render(pid, count)
}

func (c App) Start() revel.Result {
	_, err := app.DB.Exec("DELETE FROM scraping")
	if err != nil {
		log.Fatal(err)
	}
	fileName, _ := revel.Config.String("exe.file")
	command := "nohup php " + fileName + "sc.php > /dev/null &"
	exec.Command("sh", "-c", command).Start()
	log.Printf("file: %s", command)
	return c.Redirect(App.Scraping)
}

func (c App) Stop() revel.Result {
	pid := c.Params.Get("pid")
	var pidInt int
	pidInt, _ = strconv.Atoi(pid)
	process, _ := os.FindProcess(pidInt)
	log.Printf("pid: %d", pidInt)
	process.Kill()
	fileName, _ := revel.Config.String("exe.file")
	exec.Command("sh", "-c", "rm -rf " + fileName + "lock/*").Start()
	return c.Redirect(App.Scraping)
}

func init() {
	revel.InterceptFunc(checkAuthentifcation, revel.BEFORE, &App{})
}