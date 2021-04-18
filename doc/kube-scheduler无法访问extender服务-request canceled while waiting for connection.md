# kube-scheduler无法访问extender服务-request canceled while waiting for connection

```
E0414 14:51:51.645900       1 factory.go:469] Error scheduling kube-system/webapp: Post http://kube-scheduler-extender.kube-system.svc.cluster.local:8080/scheduler/predicates/always_true: net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers); retrying
I0414 14:51:51.645944       1 scheduler.go:769] Updating pod condition for kube-system/webapp to (PodScheduled==False, Reason=Unschedulable)
E0414 14:51:51.645959       1 scheduler.go:638] error selecting node for pod: Post http://kube-scheduler-extender.kube-system.svc.cluster.local:8080/scheduler/predicates/always_true: net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)
```

搜索[note-cloud](https://github.com/generals-space/note-cloud)中的相关文档.

关键字: dig resolv.conf dnsPolicy ClusterFirstWithHostNet
