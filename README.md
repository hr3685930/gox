# GOX框架
## 简介
gox是轻量级web和rpc框架, 包含大量微服务框架和工具
1. HTTP/RPC
2. Rabbitmq/Kafka队列
3. Redis缓存
4. 事件驱动
5. 日志告警
6. 文件系统
7. 限流熔断降级
8. 链路监控
9. k8s/istio部署方便


## 设计理念
gox是一个Go语言实现的http和rpc框架, 保持简单为第一原则, 支持云原生架构, 更好的与容器、Kubernetes、DevOps、service mesh、serverless 等云原生关键技术融合到一起, 可配置的组件支持,  因此您可以根据喜好来选择库进行集成


## 快速开始
- CLI工具
```
1. 开启go module (goland开发工具需在编辑器里设置)
	go env -w GO111MODULE=on
	go env -w GOPROXY=https://goproxy.cn,direct
2. 安装gctl
	go install github.com/hr3685930/pkg/gctl
```	
- 初始化
```
1. 创建一个新项目
	mkdir demo && cd demo
2.1. 初始化一个http服务
	gctl new api
2.2. 初始化一个rpc服务
	gctl new rpc
3. 运行项目
	go run main.go
```


## 目录介绍
- 总览
```
├── api								api输出目录包含proto相关文件
│   └── proto
│       ├── pb
│       │   └── cloudevent.pb.go
│       └── v1
│           └── cloudevent
│               └── cloudevent.proto
├── config.yaml						项目配置文件
├── configs							各个配置文件目录
│   ├── app.go
│   ├── cache.go
│   ├── conf.go
│   ├── database.go
│   ├── queue.go
│   └── trace.go
├── go.mod
├── go.sum
├── init								项目启动时初始化组件目录
│   └── boot
│       ├── app.go
│       ├── cache.go
│       ├── command.go
│       ├── config.go
│       ├── database.go
│       ├── event.go
│       ├── governance.go
│       ├── http.go
│       ├── grpc.go
│       ├── log.go
│       ├── metric.go
│       ├── queue.go
│       ├── sentry.go
│       ├── signal.go
│       └── trace.go
├── internal							项目内部文件
│   ├── commands						一次性脚本,定时任务,消费程序
│   │   ├── command.go
│   │   ├── consumer.go
│   │   ├── event.go
│   │   └── migrate.go
│   ├── errs							自定义错误
│   │   ├── export				
│   │   │   ├── event.go
│   │   │   ├── goroutine.go
│   │   │   ├── http.go
│   │   │   ├── grpc.go
│   │   │   ├── queue.go
│   │   │   └── report.go
│   │   │── http.go					http错误定义(可选)
│   │   └── http.go					rpc错误定义(可选)
│   ├── events						事件监听目录
│   │   ├── event.go					事件定义及绑定
│   │   └── listener					监听者目录
│   │       └── example.go
│   ├── http							http入口 (可选)
│   │   └── handler
│   │       ├── event.go
│   │       └── router.go
│   ├── rpc							rpc入口 (可选)
│   │   └── event.go
│   ├── jobs							队列目录
│   │   └── example.go
│   ├── models						db的model
│   ├── repo							repository层
│   │   └── repo.go					repository绑定单例
│   ├── types							request,response定义及转换
│   └── utils							自定义工具包
│       ├── format
│       │   ├── datacodec.go
│       │   └── protobuf.go
│       └── kafka.go
├── main.go							入口文件
├── storage							存储目录
│   └── log
└── test								单元测试
    └── main_test.go
```


## 开发规范
	1. 分层定义
		- http/rpc层
	 	该层级主要实现请求和返回的转换, 参数验证及调用service或者repo
		- service (可选)
		类似 DDD 的 application 层, 聚合多个repo, 同时协作各类service
		- repository
		接口定义方式实现,提供业务数据访问, 微服务, db, cache, es等的调用
	2. 命名规范
		1. 文件
			- 全部小写
			- 除unit test外避免下划线(_)
	3. 编码规范
		1. context
		定义方法需要在第一个参数,如:func xxx(ctx context.context,…)
		2. 变量
		不可exported的首字母小写 驼峰命名


## 工具组件使用
- gctl
```
NAME:
   gctl new - 创建项目

USAGE:
   gctl new [command options] [arguments...]

OPTIONS:
   --name value    项目名称,和go mod同名
   --err value     错误上报,支持sentry
   --trace value   链路,支持jaeger
   --metric value  监控,支持prom
   --db value      数据库,支持mysql,postgre,clickhouse,sqlite
   --queue value   队列,支持kafka,rabbitmq
   --cache value   缓存,支持redis


NAME:
   gctl repo - 创建repo

USAGE:
   gctl repo [command options] [arguments...]

OPTIONS:
   --type value   repo类型 db,api,es
   --cache value  repo 增加cache层
   --dir value    repo生成的路径, 没有则创建, repo名称根据最后一层级来命名
   --model value  model名称, 需要放在models目录下, type为db时该字段生效
```
	
- http
```
该组件可选择gin或者echo来使用http服务(默认gin)
1. 在internal/http/handler下课定义你的meddlleware和route
	1.1 默认的中间件
		ErrHandler gin的错误处理
		TimeoutMiddleware  超时设置
		GovernanceMiddleware  限流熔断时的错误返回
		CustomRecovery  panic时的错误响应
	1.2 默认路由
		ping是为了更好的与健康检查相结合
		event则代表事件的监听
2. 在internal/http下创建route对应的func
3. 错误处理
	在internal/errs/http.go下简单封装了基于http状态码的错误处理,可自定义
	
```
	
- grpc
```
1. proto文件编写
syntax = "proto3";
package service.example.v1;
option go_package = "api/proto/pb;proto";

import "google/protobuf/timestamp.proto";

service Example {
  rpc ExampleInfo(ExampleInfoReq) returns (ExampleInfoRes);
}

message ExampleInfoReq {
  int64 id = 1;
}

message ExampleInfoRes {
  string msg = 1;
}

2. 生成pb文件
protoc --go_out=plugins=grpc:. api/proto/v1/*/*.proto
3. 注册rpc服务
位于/init/boot/grpc.go
err := grpcServer.Register(opts, func(s *grpc.Server) {
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, healthServer)
	grpc_prometheus.Register(s)
	proto.RegisterEventServer(s, rpcServer.NewEvent())
	proto.RegisterExampleServer(s, rpcServer.NewExample())
	reflection.Register(s)
})
	3.1 默认的中间件
		prometheus  监控
		UnaryTimeoutInterceptor  超时
		grpc_opentracing  trace
		CustomErrInterceptor  自定义error
		UnaryGovernanceServerInterceptor 熔断限流
	3.2 默认的服务
		health  健康检查
		prometheus  监控
		Event   事件
4. 在internal/rpc下创建rpc服务对应的func
5. 错误处理
	在internal/errs/grpc.go下简单封装了基于http状态码的错误处理,可自定义

```
	
- config 
```
微服务或者说云原生应用的配置最佳实践是将配置文件和应用代码分开管理, 不将配置文件放入代码仓库，也不打包进容器镜像，而是在服务运行时，把配置文件挂载进去加载。gox的config组件就是用来帮助应用从k8s configmap中加载配置。
项目中的配置文件存在于configs/conf.go文件
配置文件的优先级为config.yaml > env环境变量 > default
1. app
项目的基本配置
2. database
数据库的相关配置,支持mysql, sqlite, postgre, clickhouse, 可设置为多个数据库连接, default为默认的数据库连接
3. queue
队列相关配置, 支持rabbitmq, kafka, local
4. cache
缓存相关配置, 支持redis, sync
5. trace
外部链路的endpoint
6. 自定义配置
需定义配置的结构体,在conf.go增加对应的配置变量
7. 配置文件config.yaml
该配置文件为本地开发环境配置文件,其他环境请使用env环境变量或文件挂载方式
```
	
- queue
```
队列的用法:
	项目中支持kafka, rabbitmq
	以kafka为例:
	1. Producer
	queue.NewProducer("topic", "", []byte(), 0)
	参数依次为topic, key(用于顺序消费), 消息体, 延迟消费时间
	2. Consumer
		2.1 绑定consumer命令internal/commands/command.go
		2.2 internal/commands/consumer
			queue.NewConsumer("example-topic", queue.Consumers})
			参数1为topic, 参数2为订阅的consumer, 可以为多个
				{
					Queue:   "example-queue",  kafka的group
					Job:     &jobs.Example{},  绑定执行的job实例
					Sleep:   0,					重试后sleep多久
					Retry:   0,					重试次数
					Timeout: 0,					超时时间, 0为不超时
				}
			job实例需实现的interface
			Handler() (queueErr *Error)
```
	
- event
```
事件处理是基于cloudevent 1.0协议实现, 通过eventbus来分发事件, 在项目中一般用于数据更新操作
事件发送
	1. channel
		进程内事件,异步
		event.NewChannelEvent(eventName)
	2. http
		跨进程事件,同步
		event.NewHttpEvent(endpoint, eventName)
	3. rpc
		跨进程事件,同步
		event.NewRpcEvent(endpoint, eventName)
	4. kafka
		跨进程事件,异步
		event.NewKafkaEvent(topic string, eventName string)
事件监听
	1. 位于internal/listener
	2. 实现interface
	Handler(ctx context.Context, event cloudevents.Event) error
	3. 绑定:  internal/event.go
		// Listeners Listeners
		var Listeners = map[string][]Listener{
			"com.example.create": {
				listener.NewExample(),
			},
		}
```
	
- cache
```
缓存功能, 支持sync/redis
使用方法: cache.cached.func
Contains check if a cached key exists
Delete remove the cached key
Fetch retrieve the cached key value
FetchMulti retrieve multiple cached keys value
Flush remove all cached keys
Save cache a value by key
若要使用reids更多功能需获取redis实例
```
	
- log
```
日志为uber zap的实现
使用方法:
	zap.Logger.Info("xxx")
```
	
- sentry
```
错误日志处理, 支持sentry, 可自定义对接阿里云, fluentd等
sentry例子:
官方sdk, 支持并发, 位于intelnal/errs/export/report.go
go func(localHub *sentry.Hub) {
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetExtras(option)
	})
	localHub.CaptureMessage(msg)
}(sentry.CurrentHub().Clone())

```
	
- trace
```
链路追踪 (可选, 若使用service mesh或者serverless平台可忽略此插件)
基于分布式追踪协议opentracing的jaeger实现
集成在项目中的组件有:
	http server/client
	grpc server/client
	redis
```
	
- metrics
```
监控指标 (可选, 若使用server mesh或者serverless平台可忽略此插件)
基于云原生、高度可扩展的指标协议OpenMetrics的Prometheus实现
集成在项目中的组件有:
	http server
	grpc server
	redis
	gorm
```
	
- sentinel
```
服务治理
该组件通过系统自适应流控从整体维度对应用入口流量进行控制，结合系统的 Load、CPU 使用率以及应用的入口 QPS、平均响应时间和并发量等几个维度的监控指标，通过自适应的流控策略，让系统的入口流量和系统的负载达到一个平衡，让系统尽可能跑在最大吞吐量的同时保证系统整体的稳定性。
系统保护规则是应用整体维度的，而不是单个调用维度的，并且仅对入口流量生效。
基于TCP BBR
// 自适应流控，启发因子为 load1 >= 8
_, err := system.LoadRules([]*system.SystemRule{
	{
		MetricType:system.Load,
		TriggerCount:8.0,
		Strategy:system.BBR,
	},
})

```
	
- command
```
命令行模式基于github.com/urfave/cli包实现
var Commands = []cli.Command{
	{
		Name:    "db",
		Usage:   "db操作",
		Subcommands: []cli.Command{
			{
				Name:   "migrate",
				Usage:  "迁移数据表",
				Action: Migrate,
			},
		},
	},
	{
        Name:    "queue",
        Usage:   "队列job",
        Subcommands: []cli.Command{
            {
                Name:   "example",
                Usage:  "example消费程序",
                Action: Example,
            },
        },
    },
    {
        Name:    "event",
        Usage:   "kafka事件监听",
        Subcommands: []cli.Command{
            {
                Name:   "listen",
                Usage:  "example topic事件监听",
                Action: Event,
            },
        },
    },
}
```
	
- database
```
db基于gorm 2.0支持mysql, sqlite, clickhouse, postgresql
使用方法:
	db.Orm来获取orm默认DB实例
	db.GetConnect获取其他DB实例
 
其他使用请移步: https://gorm.io/zh_CN/docs/index.html
```

## 最佳实践

- 错误日志处理
```
在日志中分为错误日志,请求日志,查询日志....
比如请求日志则直接在网关层上报,错误日志在项目中上报, 查询日志在各类云厂商集成上报
处理日志的方式有很多种:
log文件, 错误上报, 控制台抓取等等...无论那种方式我们都要定义我的抛出的错误信息,环境,类型,堆栈,请求,返回等,在gox里面封装了grpc/http/queue/goroutine等错误,并且支持错误上报至各个平台,在internal/errs/export文件夹下可以看到错误处理时需要的信息以及上报的平台
自定义err接口
type Error interface {
	error
	GetStack() string
}
如何使用?
	http服务:
		使用errs.InternalError("失败") 返回自定义错误
		gin返回错误:_ = c.Error(err)
	grpc服务:
		使用errs.InternalError("失败") 返回自定义错误
	queue:
		queue.Err("xxx")

在项目中避免使用其他error组件
		
```

- 服务部署
```
1. 安装组件
	gitlab
	gitlab runner
	k8s
	helm v3
2. 定义Dockerfile
FROM golang:1.17
ENV GO111MODULE=on 
    GOPROXY=https://goproxy.cn,direct 
WORKDIR /app
COPY . .
RUN go build .
EXPOSE 80
ENTRYPOINT ["./main]
3. 定义CI文件
推荐使用kaniko build和 helm 来部署
4. helm部署
增加deployment,job,cornjob,consumer,service类型
```

- testing
```
单元测试
github.com/stretchr/testify  断言相关操作
github.com/brianvoe/gofakeit  模拟数据
github.com/golang/mock/gomock  mock
http: 自带一个ping的测试
func TestPing(t *testing.T) {
	hs := gin.NewHTTPServer(configs.ENV.App.Debug)
	_ = hs.LoadRoute(handler.Route)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping",  nil)
	hs.G.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

grpc: 自带一个Health的测试
	
func TestHealth(t *testing.T) {
	conn,err := GrpcClient()
	defer conn.Close()
	client := healthpb.NewHealthClient(conn)
	resp, err := client.Check(context.Background(), &healthpb.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("health failed: %v", err)
	}
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.Status)
}
	
```

- 事件驱动伸缩
```
在k8s环境下, 服务扩容基本上是通过cpu, memery的阈值来进行扩容, 自从serverless框架knative出现, 以事件驱动来实现业务, 为了更精确的扩容, 于是出现了通过获取k8s服务的事件来进行扩容, 比如在kafka队列中, 我们可以通过后滞消息数来进行扩容, 当然不止是kafka, 事件驱动扩容适用于各种供应商、数据库、消息系统、遥测系统、CI/CD 等的开箱即用的伸缩器, 这里推荐一款云原生事件驱动自动伸缩 KEDA
```

- 本地环境
```
微服务下的本地开发环境怎么做?
这个问题一直困扰了很久, 服务有依赖怎么调试, 数据库, 队列, 网关下怎么调试?
和前端联调怎么调试? 网络怎么打通?
于是乎衍生出了telepresence和kt connect等开发工具来实现与k8s环境的打通
这里我推荐阿里的kt connect
KtConnect提供了本地和测试环境集群的双向互联能力。 
具体使用: https://alibaba.github.io/kt-connect/#/zh-cn/guide/quickstart

若你的项目已经上了service mesh还可以通过KtEnv的隔离域来实现联调
具体使用: https://alibaba.github.io/virtual-environment
```

- Protocol Buffers IDL 
```
我们在开发多个rpc服务中, 经常会遇到proto文件更新, 调用方则会去拷贝新的proto文件, 多次交互会影响开发效率, 于是乎protodep出现, 解决了proto文件的依赖问题
使用:
1. 安装
go get github.com/hr3685930/protodep
2. 定义
	  proto_outdir = "./api/proto"
    [[dependencies]]
      target = "git.kid17.com/tiny/library/service/proto/auth"
      branch = "master"
      path = "auth"
      protocol = "https"
    [[dependencies]]
      target = "git.kid17.com/tiny/library/service/proto/feed"
      branch = "master"
      path = "feed"
      protocol = "https"
3. 拉取
protodep up -f --basic-auth-username $(PROTOUSER)  --basic-auth-password $(PROTOPWD)
```

- 并发原语
```
在golang中有很多并发原语,提升我们的开发效率,gox封装了部分原语

无需等待的执行goroutine,增加了error处理
func GO(fns AsyncFunc) 

需等待的执行goroutine
方法1: 无序返回
g := goo.NewGroup(num int)
g.One(ctx context.Context, fn SyncFunc) 
g.Wait()

方法2: 顺序返回
func All(ctx context.Context, fns ...SyncFunc) ([]interface{}, []error)

其他原语:
1. singleflight (请求合并)
可解决缓存更新问题
它的作用是，在处理多个goroutine 同时调用同一个函数的时候，只让一个 goroutine 去调用这个函数，等到这个goroutine 返回结果的时候，再把结果返回给这几个同时调用的 goroutine，这样可以减少并发调用的数量
2. CyclicBarrier (循环栅栏)
允许一组 goroutine 彼此等待，到达一个共同的执行点。同时，因为它可以被重复使用，所以叫循环栅栏。具体的机制是，大家都在栅栏前等待，等全部都到齐了，就抬起栅栏放行。
3. ErrGroup
我们经常会碰到需要将一个通用的父任务拆成几个小任务并发执行的场景，其实，将一个大的任务拆成几个小任务并发执行
4. gollback
用来处理一组子任务的执行的，不过它解决了 ErrGroup 收集子任务返回结果的痛点。使用 ErrGroup 时，如果你要收到子任务的结果和错误，你需要定义额外的变量收集执行结果和错误，但是这个库可以提供更便利的方式。
5. Hunch
提供的功能和 gollback 类似，不过它提供的方法更多，而且它提供的和gollback 相应的方法，也有一些不同。
5.1. Waterfall 方法
func Waterfall(parentCtx context.Context, execs …ExecutableInSequence) (I interface{}, e error)
它其实是一个 pipeline 的处理方式，所有的子任务都是串行执行的，前一个子任务的执行结果会被当作参数传给下一个子任务，直到所有的任务都完成，返回最后的执行结果。一旦一个子任务出现错误，它就会返回错误信息，执行结果（第一个返回参数）为 nil。
6. schedgroup
可以指定任务在某个时间或者某个时间之后执行
type Group
func New(ctx context.Context) *Group
func(g *Group) Delay(delay time.Duration, fn func())
func(g *Group) Schedule(when time.Time, fn func())
func(g *Group) Wait()error
6.1. Delay 和 Schedule
它们的功能其实是一样的，都是用来指定在某个时间或者之后执行一个函数。只不过，Delay 传入的是一个 time.Duration 参数，它会在 time.Now()+delay 之后执行函数，而Schedule 可以指定明确的某个时间执行。
6.2. Wait
这个方法调用会阻塞调用者，直到之前安排的所有子任务都执行完才返回。如果 Context被取消，那么，Wait 方法会返回这个 cancel error。
注意点:
第一点是，如果调用了 Wait 方法，你就不能再调用它的 Delay 和 Schedule 方法，否则会 panic。第二点是，Wait 方法只能调用一次，如果多次调用的话，就会 panic。
你可能认为，简单地使用 timer 就可以实现这个功能。其实，如果只有几个子任务，使用timer 不是问题，但一旦有大量的子任务，而且还要能够 cancel，那么，使用 timer 的话，CPU 资源消耗就比较大了。所以，schedgroup 在实现的时候，就使用container/heap，按照子任务的执行时间进行排序，这样可以避免使用大量的 timer，从而提高性能。
```

-  分布式事务
```
dtm 目前支持最多
支持saga， tcc， xa， 二阶段， at， 子事物屏障, workflow...
https://www.dtm.pub
例子： SAGA模式
整个SAGA事务的逻辑是：
执行转出成功=>执行转入成功=>全局事务完成
如果在中间发生错误，例如转入B发生错误，则会调用已执行分支的补偿操作，即：
执行转出成功=>执行转入失败=>执行转入补偿成功=>执行转出补偿成功=>全局事务回滚完成

Saga并发执行（默认是顺序执行）
saga := dtmcli.NewSaga(DtmServer, dtmcli.MustGenGid(DtmServer)).
			Add(Busi+"/CanRollback1", Busi+"/CanRollback1Revert", req).
			Add(Busi+"/CanRollback2", Busi+"/CanRollback2Revert", req).
			Add(Busi+"/UnRollback1", "", req).
			Add(Busi+"/UnRollback2", "", req).
			EnableConcurrent().
			AddBranchOrder(2, []int{0, 1}). // 指定step 2，需要在0，1完成后执行
			AddBranchOrder(3, []int{0, 1}) // 指定step 3，需要在0，1完成后执行

```

- serverless
```
在笔者目前的工作实践中，我们在架构设计中严格遵循着微服务的设计理念。每块业务上独立的领域作为一个微服务，每个微服务有自己单独的数据库。任何一个微服务只能读写自己的数据库，而绝对不能干涉其他微服务的数据库，所有微服务之间获取信息都是通过http/rpc进行通信。在业务初期时，由于产品处于雏形中，所有服务的接口在逻辑上都比较简单，大部分接口都可以看做各自领域模型上的增删改查形态，服务与服务之间的调用也并不密集频繁。随着业务的持续发展与产品的迭代，大大小小、许许多多的功能需求交织重叠在一起，相应的各个服务直接服务于需求的接口也逐渐变多，每个接口越来越“需求相关”。并且我们逐渐发现随着需求的逐渐复杂，每个接口涉及的服务也越来越多，很难有一个强有力的理由去确定这个接口就必须放在某个服务里。这个时候将这个接口放在哪个服务里面，往往取决于做这个需求的开发对哪个服务的掌控力更强，或者说这个接口看上去更倾向于哪块业务领域。同时由于大大小小的“需求相关”的接口越来越多的堆积在各个微服务内以后，整个服务仓库的代码量逐渐增大，代码质量也逐渐下降、微服务仿佛变得不再那么“微”。
Api聚合层，第一次听到这个名词依旧也是在今年前端领域的技术分享中。相应的、API聚合层在前端领域中也有个专门的名词叫做BFF，即Backend For Frontend。在这一块，国内已经有过不少BFF与Serverless结合的的分享。可以说BFF与Serverless的结合即解决了在前端领域中，多端适配、又或者是UI模型与后端API数据的转化这一系列问题，同时也没有引入多余的维护服务稳定性以及服务治理等一系列额外的运维负担。
使用knative来作为BFF层
knative 
推荐使用阿里云ASK集群来使用knative, 简单易操作, 用保留实例功能解决了冷启动问题, 费用低.
具体: https://help.aliyun.com/document_detail/184831.html
```
