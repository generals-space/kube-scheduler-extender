# kube-scheduler-extender

参考文章

1. [图解kubernetes调度器SchedulerExtender扩展](https://www.cnblogs.com/buyicoding/archive/2004/01/13/12250480.html)
2. [图解kubernetes调度器SchedulerExtender扩展](https://blog.csdn.net/qq_42952272/article/details/104138445)
    - 与参考文章1是同一篇, 作为备份
3. [kube-scheduler调度扩展](https://segmentfault.com/a/1190000019034882)
    - 小碗汤
4. [k8s-scheduler-extender-example](//github.com/everpeace/k8s-scheduler-extender-example.git)
    - 小碗汤的参考工程
5. [Huang-Wei/sample-scheduler-extender](https://github.com/Huang-Wei/sample-scheduler-extender)
    - 还没看
6. [Kubernetes集群调度器原理剖析及思考](http://blog.itpub.net/69908804/viewspace-2639995/)
7. [k8s调度器扩展方式](https://blog.csdn.net/yevvzi/article/details/79858170)
    - `ExtenderConfig`, 即为`scheduler-policy.json`中`extenders`数组的配置, 拥有详细的解释.
8. [Kubernetes 1.8抢占式调度Preemption源码分析](https://my.oschina.net/jxcdwangtao/blog/1563456)

## 前言

本例实现了一个简单的扩展调度器, 没有任何作用, 只展示其工作流程和部署方式.

一个简单的 scheduler-extender 只是一个普通的 http 服务器, 核心 scheduler 调度器通过 http API 与扩展调度器通信.

当然, 为了实现真正的调度器功能, 扩展调度器可能需要对集群中的资源进行比较判断, 所以还得拥有一个 kubernetes 的客户端.

```
       kube-scheduler                                kube-scheduler-extender(HttpAPI)
+------------|-------------+                          +--------------------------+
|            |             |                          |                          |
|    +-------↓--------+    |      Pod + []Node        |    +----------------+    |
|    |                ├───────────────────────────────────>|                |    |
|    |   predicates   |    |                          |    |   predicates   |    |
|    |                |<───────────────────────────────────┤                |    |
|    +-------┬--------+    |     []Node(filtered)     |    +----------------+    |
|            |             |                          |                          |
|            |             |                          |                          |
|    +-------↓--------+    |       Pod + []Node       |    +----------------+    |
|    |                ├───────────────────────────────────>|                |    |
|    |   prioritize   |    |                          |    |   prioritize   |    |
|    |                |<───────────────────────────────────┤                |    |
|    +-------┬--------+    |      []Node(scored)      |    +----------------+    |
|            |             |                          |                          |
|            |             |                          |                          |
|    +-------↓--------+    |        Pod + Node        |    +----------------+    |
|    |                ├───────────────────────────────────>|                |    |
|    |      bind      |    |                          |    |      bind      |    |
|    |                |<───────────────────────────────────┤                |    |
|    +-------┬--------+    |        Error/nil         |    +----------------+    |
|            |             |                          |                          |
|            |             |                          |                          |
+------------|-------------+                          +--------------------------+
             ↓                                                                    
```

> 上述过程中没有提到Preemption(抢占)阶段, 因为目前还没有研究过, 可见参考文章8.

## 环境准备

kubernetes: v1.17.2

go: 1.3+

## 部署运行

本工程在本地环境即可运行, 但是需要部署在集群内部才能生效.

### 编译可执行文件

```
go build -o kube-scheduler-extender
```

### 构建镜像

```
docker build -t registry.cn-hangzhou.aliyuncs.com/generals-kuber/kube-scheduler-extender:v0.0.1 .
```

### 部署扩展调度器服务

```
kubectl apply -f deploy/kube-scheduler-extender.yaml
```

等待启动完成后, 执行如下命令, 将开启扩展调度器所需的配置文件, 拷贝到目标目录下

```
cp deploy/kube-scheduler-config.yaml /etc/kubernetes/scheduler-config.yaml
cp deploy/scheduler-policy.json /etc/kubernetes/scheduler-policy.json
```

### 修改核心调度器配置

然后修改`/etc/kubernetes/manifests/kube-scheduler.yaml`文件, 将上面两个文件挂载进去, 并在`command`字段中指定其路径, 具体的修改步骤, 见[kube-scheduler.yaml](./deploy/kube-scheduler.yaml)

------

部署扩展调度器, 需要同时修改`kube-scheduler`核心调度器的配置, 指定扩展调度器的访问地址. 这样, 在经过默认的核心调度器筛选后, 再由我们的扩展调度器进行进一步筛选.

但是必须要先部署扩展调度器本身, 再修改`kube-scheduler`核心调度器的配置, 顺序不能颠倒. 否则调度器本身在调度时, `kube-scheduler`就会尝试向扩展调度器服务发请求, 出现了鸡生蛋蛋生鸡的情况...

而且每次更新`kube-scheduler-extender`, 都要先把`kube-scheduler`的`--config`字段移除, 否则新的Pod无法创建成功.
