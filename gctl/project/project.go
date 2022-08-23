package project

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const TplDir = ".goo"
const TplGitUrl = "https://gitee.com/geekers/goo-template"

type Opt struct {
	ProjectName string
	IsSentry    bool
	IsTrace     bool
	IsMetric    bool
	QueueDrive  string
	CacheDrive  string
	DBDrive     string
	ServiceType string
	IsDB        bool
}

func Create(c *cli.Context) {
	if c.NArg() != 1 {
		fmt.Print("必须存在 api 或 rpc 参数")
		return
	}
	serviceType := c.Args().Get(0)
	if serviceType != "api" && serviceType != "rpc" {
		fmt.Print("必须存在 api 或 rpc 参数")
		return
	}

	pwd, err := os.Getwd()
	TryErr(err)

	projectName := filepath.Base(pwd)
	if c.String("name") != "" {
		projectName = c.String("name")
	}

	//拉取git template
	CloneTpl()

	//生成项目
	opts := &Opt{}
	opts.ProjectName = projectName
	opts.ServiceType = serviceType
	opts.IsTrace = false
	opts.IsSentry = false
	opts.CacheDrive = "sync"
	opts.QueueDrive = "local"
	opts.IsDB = false

	if c.String("cache") != "" {
		opts.CacheDrive = c.String("cache")
	}

	if c.String("queue") != "" {
		opts.QueueDrive = c.String("queue")
	}

	if c.String("db") != "" {
		opts.DBDrive = c.String("db")
		opts.IsDB = true
	}

	if c.String("err") == "sentry" {
		opts.IsSentry = true
	}

	if c.String("trace") == "jaeger" {
		opts.IsTrace = true
	}

	if c.String("metric") == "prom" {
		opts.IsMetric = true
	}
	CreateProject(opts, pwd)
	TryErr(os.RemoveAll(".goo"))
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

func SimpleCreate(configFile, tempDir string, opt interface{}) {
	t, err := template.ParseFiles(tempDir)
	TryErr(err)
	f, err := os.Create(configFile)
	TryErr(err)
	TryErr(t.Execute(f, opt))
}

func osExecClone(workspace, url string) error {
	cmd := exec.Command("git", "clone", url, workspace)
	_, err := cmd.CombinedOutput()
	return err
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

func CheckFile(filepath string) (exist bool) {
	fileInfo, e := os.Stat(filepath)
	if fileInfo != nil && e == nil {
		exist = true
	} else if os.IsNotExist(e) {
		exist = false
	}
	return
}

func CloneTpl() {
	TryErr(os.Mkdir(TplDir, os.ModePerm))
	_ = os.RemoveAll(TplDir)
	err := osExecClone(TplDir, TplGitUrl)
	TryErr(err)
}

func CreateProject(opts *Opt, pwd string) {
	apiDir := pwd + "/api"
	TryErr(os.Mkdir(apiDir, os.ModePerm))
	TryErr(os.MkdirAll(apiDir+"/proto/v1/cloudevent", os.ModePerm))
	SimpleCreate(apiDir+"/proto/v1/cloudevent/cloudevent.proto", TplDir+"/api/proto/v1/cloudevent/cloudevent.proto", opts)
	if opts.ServiceType == "rpc" {
		TryErr(os.MkdirAll(apiDir+"/proto/v1/example", os.ModePerm))
		SimpleCreate(apiDir+"/proto/v1/example/example.proto", TplDir+"/api/proto/v1/example/example.proto", opts)
	}

	// config
	configDir := pwd + "/configs"
	TryErr(os.Mkdir(configDir, os.ModePerm))
	SimpleCreate(configDir+"/app.go", TplDir+"/configs/app.tpl", opts)
	SimpleCreate(configDir+"/cache.go", TplDir+"/configs/cache.tpl", opts)
	SimpleCreate(configDir+"/conf.go", TplDir+"/configs/conf.tpl", opts)
	if opts.IsDB {
		SimpleCreate(configDir+"/database.go", TplDir+"/configs/database.tpl", opts)
	}
	SimpleCreate(configDir+"/queue.go", TplDir+"/configs/queue.tpl", opts)
	if opts.IsTrace {
		SimpleCreate(configDir+"/trace.go", TplDir+"/configs/trace.tpl", opts)
	}

	// init boot
	initBootDir := pwd + "/init/boot"
	TryErr(os.MkdirAll(initBootDir, os.ModePerm))
	SimpleCreate(initBootDir+"/app.go", TplDir+"/init/boot/app.tpl", opts)
	SimpleCreate(initBootDir+"/cache.go", TplDir+"/init/boot/cache.tpl", opts)
	SimpleCreate(initBootDir+"/command.go", TplDir+"/init/boot/command.tpl", opts)
	SimpleCreate(initBootDir+"/config.go", TplDir+"/init/boot/config.tpl", opts)
	if opts.DBDrive != "" {
		SimpleCreate(initBootDir+"/database.go", TplDir+"/init/boot/database.tpl", opts)
	}
	SimpleCreate(initBootDir+"/governance.go", TplDir+"/init/boot/governance.tpl", opts)
	SimpleCreate(initBootDir+"/log.go", TplDir+"/init/boot/log.tpl", opts)
	SimpleCreate(initBootDir+"/queue.go", TplDir+"/init/boot/queue.tpl", opts)
	SimpleCreate(initBootDir+"/signal.go", TplDir+"/init/boot/signal.tpl", opts)
	SimpleCreate(initBootDir+"/event.go", TplDir+"/init/boot/event.tpl", opts)
	if opts.IsTrace {
		SimpleCreate(initBootDir+"/trace.go", TplDir+"/init/boot/trace.tpl", opts)
	}
	if opts.IsSentry {
		SimpleCreate(initBootDir+"/sentry.go", TplDir+"/init/boot/sentry.tpl", opts)
	}
	if opts.ServiceType == "api" {
		SimpleCreate(initBootDir+"/http.go", TplDir+"/init/boot/http.tpl", opts)
	} else {
		SimpleCreate(initBootDir+"/grpc.go", TplDir+"/init/boot/grpc.tpl", opts)
	}
	if opts.IsMetric {
		SimpleCreate(initBootDir+"/metric.go", TplDir+"/init/boot/metric.tpl", opts)
	}

	// commands
	commandsDir := pwd + "/internal/commands"
	TryErr(os.MkdirAll(commandsDir, os.ModePerm))
	SimpleCreate(commandsDir+"/command.go", TplDir+"/internal/commands/command.tpl", opts)
	SimpleCreate(commandsDir+"/consumer.go", TplDir+"/internal/commands/consumer.tpl", opts)
	SimpleCreate(commandsDir+"/migrate.go", TplDir+"/internal/commands/migrate.tpl", opts)
	SimpleCreate(commandsDir+"/event.go", TplDir+"/internal/commands/event.tpl", opts)

	//errsExportDir
	errsExportDir := pwd + "/internal/errs/export"
	TryErr(os.MkdirAll(errsExportDir, os.ModePerm))
	SimpleCreate(errsExportDir+"/goroutine.go", TplDir+"/internal/errs/export/goroutine.tpl", opts)
	SimpleCreate(errsExportDir+"/queue.go", TplDir+"/internal/errs/export/queue.tpl", opts)
	SimpleCreate(errsExportDir+"/report.go", TplDir+"/internal/errs/export/report.tpl", opts)
	SimpleCreate(errsExportDir+"/event.go", TplDir+"/internal/errs/export/event.tpl", opts)
	if opts.ServiceType == "api" {
		SimpleCreate(errsExportDir+"/http.go", TplDir+"/internal/errs/export/http.tpl", opts)
	} else {
		SimpleCreate(errsExportDir+"/grpc.go", TplDir+"/internal/errs/export/grpc.tpl", opts)
	}

	//errsDir
	errsDir := pwd + "/internal/errs"
	TryErr(os.MkdirAll(errsDir, os.ModePerm))
	if opts.ServiceType == "api" {
		SimpleCreate(errsDir+"/http.go", TplDir+"/internal/errs/http.tpl", opts)
	} else {
		SimpleCreate(errsDir+"/grpc.go", TplDir+"/internal/errs/grpc.tpl", opts)
	}

	//eventDir
	eventDir := pwd + "/internal/events"
	TryErr(os.MkdirAll(eventDir, os.ModePerm))
	SimpleCreate(eventDir+"/event.go", TplDir+"/internal/events/event.tpl", opts)

	//eventListenDir
	eventListenDir := pwd + "/internal/events/listener"
	TryErr(os.MkdirAll(eventListenDir, os.ModePerm))
	SimpleCreate(eventListenDir+"/example.go", TplDir+"/internal/events/listener/example.tpl", opts)

	if opts.ServiceType == "api" {
		//handlerDir
		handlerDir := pwd + "/internal/http/handler"
		TryErr(os.MkdirAll(handlerDir, os.ModePerm))
		SimpleCreate(handlerDir+"/router.go", TplDir+"/internal/http/handler/router.tpl", opts)
		SimpleCreate(handlerDir+"/event.go", TplDir+"/internal/http/handler/event.tpl", opts)
	} else {
		// rpcDir
		rpcDir := pwd + "/internal/rpc"
		TryErr(os.MkdirAll(rpcDir, os.ModePerm))
		SimpleCreate(rpcDir+"/event.go", TplDir+"/internal/rpc/event.tpl", opts)
	}

	// job
	jobDir := pwd + "/internal/jobs"
	TryErr(os.MkdirAll(jobDir, os.ModePerm))
	SimpleCreate(jobDir+"/example.go", TplDir+"/internal/jobs/example.tpl", opts)

	// models
	modelsDir := pwd + "/internal/models"
	TryErr(os.MkdirAll(modelsDir, os.ModePerm))

	//repo
	repoDir := pwd + "/internal/repo"
	TryErr(os.MkdirAll(repoDir, os.ModePerm))
	SimpleCreate(repoDir+"/repo.go", TplDir+"/internal/repo/repo.tpl", opts)

	// types
	typesDir := pwd + "/internal/types"
	TryErr(os.MkdirAll(typesDir, os.ModePerm))

	// utils
	utilsDir := pwd + "/internal/utils"
	TryErr(os.MkdirAll(utilsDir, os.ModePerm))
	SimpleCreate(utilsDir+"/kafka.go", TplDir+"/internal/utils/kafka.tpl", opts)
	// utils format
	utilsFormatDir := pwd + "/internal/utils/format"
	TryErr(os.MkdirAll(utilsFormatDir, os.ModePerm))
	SimpleCreate(utilsFormatDir+"/datacodec.go", TplDir+"/internal/utils/format/datacodec.tpl", opts)
	SimpleCreate(utilsFormatDir+"/protobuf.go", TplDir+"/internal/utils/format/protobuf.tpl", opts)

	// log
	storageDir := pwd + "/storage/log"
	TryErr(os.MkdirAll(storageDir, os.ModePerm))
	SimpleCreate(storageDir+"/.gitignore", TplDir+"/storage/log/.gitignore", opts)

	// test
	testDir := pwd + "/test"
	TryErr(os.MkdirAll(testDir, os.ModePerm))
	SimpleCreate(testDir+"/.gitignore", TplDir+"/test/.gitignore", opts)
	SimpleCreate(testDir+"/main_test.go", TplDir+"/test/main_test.tpl", opts)
	if opts.ServiceType != "api" {
		SimpleCreate(testDir+"/grpc.go", TplDir+"/test/grpc.tpl", opts)
	}

	// root
	SimpleCreate(pwd+"/config.yaml", TplDir+"/config.yaml", opts)
	SimpleCreate(pwd+"/main.go", TplDir+"/main.tpl", opts)

	isGitignore := CheckFile(".gitignore")
	if isGitignore {
		fmt.Println("you need add config.yaml\nsqlite.db\n to .gitignore")
	} else {
		fmt.Println("gitignore not found, creating gitignore...")
		SimpleCreate(pwd+"/.gitignore", TplDir+"/.gitignore.tpl", opts)
	}

	// 检查是否存在go.mod
	isMod := CheckFile("go.mod")
	if !isMod {
		fmt.Println("mod not found, creating mod.go...")
		TryErr(ExecShell("go mod init " + opts.ProjectName))
	}

	TryErr(ExecShell("protoc --go_out=plugins=grpc:. api/proto/v1/cloudevent/*.proto"))
	TryErr(ExecShell("go mod tidy -compat=1.17"))
}

func ExecShell(shell string) error {
	cmd := exec.Command("/bin/bash", "-c", shell)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return err
}
