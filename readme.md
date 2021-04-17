# kube-scheduler-extender

## 前言

本例实现了一个简单的扩展调度器, 没有任何作用, 只展示其工作流程和部署方式.

## 部署运行

首先部署扩展调度器服务

```
kubectl apply -f doc/deploy/kube-scheduler-extender.yaml
```

等待启动完成后, 执行如下命令, 将开启扩展调度器所需的配置文件, 拷贝到目标目录下

```
cp doc/deploy/kube-scheduler-config.yaml /etc/kubernetes/scheduler-config.yaml
cp doc/deploy/scheduler-policy.json /etc/kubernetes/scheduler-policy.json
```

> 其中, `kube-scheduler-config.default.yaml`为, 通过`exec`进入到`kube-scheduler`核心调度器容器中, 执行`kube-scheduler --write-config-to config.default.yaml`得到的默认配置文件. 
> 
> 不过需要将原生`kube-scheduler`先停止, 否则该命令无法执行, 这就需要手动修改`/etc/kubernetes/manifests/kube-scheduler.yaml`中`command`及`livenessProbe`字段, 启动Pod但不启动调度器进程.

然后修改`/etc/kubernetes/manifests/kube-scheduler.yaml`文件, 将这两个文件挂载进去, 并在`command`字段中指定其路径, 具体的修改步骤, 见[kube-scheduler.yaml](./doc/deploy/kube-scheduler.yaml)

------

部署扩展调度器, 需要同时修改`kube-scheduler`核心调度器的配置, 指定扩展调度器的访问地址. 这样, 在经过默认的核心调度器筛选后, 再由我们的扩展调度器进行进一步筛选.

但是必须要先部署扩展调度器本身, 再修改`kube-scheduler`核心调度器的配置, 顺序不能颠倒. 否则调度器本身在调度时, `kube-scheduler`就会尝试向扩展调度器服务发请求, 出现了鸡生蛋蛋生鸡的情况...

而且每次更新`kube-scheduler-extender`, 都要先把`kube-scheduler`的`--config`字段移除, 否则新的Pod无法创建成功.

```
go build -o kube-scheduler-extender
```

优选算法, 必须要同时定义`weight`字段.

```
couldn't create scheduler from policy: Priority for extender http://kube-scheduler-extender.kube-system.svc.cluster.local:8080/scheduler should have a positive weight applied to it
```


预选, 返回一个 NodeList 对象, 包含所有符合条件的 Node 对象.

