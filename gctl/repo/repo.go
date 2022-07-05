package repo

import (
	"errors"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"text/template"
)

const TplDir = ".goo"
const TplGitUrl = "https://gitee.com/geekers/goo-template"

type Opt struct {
	Types       string
	Model       string
	ProjectName string
	FileName    string
}

func Create(c *cli.Context) {
	dir := c.String("dir")
	if dir == "" {
		TryErr(errors.New("dir 参数不存在"))
	}
	types := c.String("type")
	if types != "simple" && types != "db" {
		TryErr(errors.New("type 参数不存在"))
	}
	model := c.String("model")
	if types == "db" && model == "" {
		TryErr(errors.New("model 参数不存在"))
	}
	//拉取git template
	CloneTpl()

	pwd, err := os.Getwd()
	TryErr(err)

	projectName := filepath.Base(pwd)
	//生成项目
	opts := &Opt{}
	opts.Types = types
	opts.Model = model
	opts.ProjectName = projectName
	fileName := path.Base(dir)
	opts.FileName = fileName
	TryErr(os.Mkdir(dir, os.ModePerm))
	SimpleCreate(dir+"/"+fileName+".go", TplDir+"/internal/repo/impl.tpl", opts)
	SimpleCreate(dir+"/"+fileName+"_db.go", TplDir+"/internal/repo/impl_db.tpl", opts)
	TryErr(os.RemoveAll(".goo"))
}

func SimpleCreate(configFile, tempDir string, opt interface{}) {
	t, err := template.ParseFiles(tempDir)
	TryErr(err)
	f, err := os.Create(configFile)
	TryErr(err)
	TryErr(t.Execute(f, opt))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func TryErr(err error) {
	is, fileErr := PathExists(TplDir)
	if fileErr != nil {
		panic(err)
	}
	if is && err != nil {
		_ = os.RemoveAll(TplDir)
	}
	if err != nil {
		panic(err)
	}
}

func CloneTpl() {
	TryErr(os.Mkdir(TplDir, os.ModePerm))
	_ = os.RemoveAll(TplDir)
	err := osExecClone(TplDir, TplGitUrl)
	TryErr(err)
}

func osExecClone(workspace, url string) error {
	cmd := exec.Command("git", "clone", url, workspace)
	_, err := cmd.CombinedOutput()
	return err
}
