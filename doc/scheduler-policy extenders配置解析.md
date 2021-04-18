# scheduler-policy extenders配置解析

参考文章

1. [k8s调度器扩展方式](https://blog.csdn.net/yevvzi/article/details/79858170)
    - `ExtenderConfig`, 即为`scheduler-policy.json`中`extenders`数组的配置, 拥有详细的解释.
2. [Kubernetes 1.8抢占式调度Preemption源码分析](https://my.oschina.net/jxcdwangtao/blog/1563456)

```json
    "extenders": [
        {
            "urlPrefix.wrong": "http://kube-scheduler-extender:8080/scheduler",
            "urlPrefix": "http://kube-scheduler-extender.kube-system.svc:8080/scheduler",
            "filterVerb": "predicates",
            "prioritizeVerb": "prioritize",
            "weight": 1,
            "preemptVerb": "preemption",
            "bindVerb": "bind",
            "enableHttps": false,
            "nodeCacheCapable": false
        }
    ]
```

`urlPrefix`: 扩展调度器的访问地址, 核心 scheduler 调度器将向这个地址发送各种请求.

## filterVerb 预选

预选接口的访问uri, 实际的地址将为`urlPrefix`+`/`+`filterVerb`, 如`http://kube-scheduler-extender.kube-system.svc:8080/scheduler/predicates`.

核心 scheduler 调度器(在经过内置的预选方案后)将待调度的 Pod 对象, 与可供选择的 Node 列表对象发送给扩展调度器. 由后者按照自身的业务逻辑, 将一部分 Node 进行过滤并将过滤结果返回. 然后前者将在剩余的 Node 中进行后续阶段的比对.

## prioritizeVerb 优选

核心 scheduler 调度器(在经过内置的优选方案后)将待调度的 Pod 对象, 与可供选择的 Node 列表对象发送给扩展调度器. 由后者按照自身的业务逻辑, 对可供选择的 Node 对象进行打分, 并将结果返回. 得分越高的 Node, 将会优先与目标 Pod 进行匹配.

需要注意的是, 启用优选接口时, 必须要同时定义`weight(权重)`字段, 否则在核心 scheduler 调度器请求扩展调度器的优选接口时, 会出现如下错误.

```
couldn't create scheduler from policy: Priority for extender http://kube-scheduler-extender.kube-system.svc.cluster.local:8080/scheduler should have a positive weight applied to it
```

另外, 当 k8s 集群为单节点时, 优选接口是没有意义的. 当经过预选阶段后, 如果只剩下一个 Node 节点, 那么核心 scheduler 调度器也不会再进入优选阶段. 即不再请求扩展调度器的优选接口.

## preemptVerb 抢占

本工程中没有提到Preemption(抢占)阶段, 因为目前还没有研究过, 可见参考文章2.

## bindVerb 绑定

核心 scheduler 调度器(在经过内置的优选方案后)将待调度的 Pod 对象与其将要被调度上的 Node 对象发送给扩展调度器. 由后者按照自身的业务逻辑进行处理(比如写入数据库), 然后由后者自行创建`corev1.Binding{}`对象, 设置 Pod 与 Node 的绑定关系, 并返回绑定结果.

如果后者不设置 Pod 与 Node 的绑定关系, 那么虽然 Pod 会调度成功, 但是将会一直处于`Pending`状态, 不会出现在目标主机上.

```
webapp           0/1     Pending   0          4m15s   <none>          <none>          <none>           <none>
```

```
Events:
  Type    Reason     Age        From               Message
  ----    ------     ----       ----               -------
  Normal  Scheduled  <unknown>  default-scheduler  Successfully assigned kube-system/webapp-02 to k8s-worker-01
```
