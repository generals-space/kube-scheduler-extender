当前目录下的`kube-scheduler-config.default.yaml`文件为, 通过`exec`进入到`kube-scheduler`核心调度器容器中, 执行`kube-scheduler --write-config-to config.default.yaml`得到的默认配置文件. 

不过需要将原生`kube-scheduler`先停止, 否则该命令无法执行, 这就需要手动修改`/etc/kubernetes/manifests/kube-scheduler.yaml`中`command`及`livenessProbe`字段, 启动Pod但不启动调度器进程.

本工程中并不会使用到此文件, 仅作为参考. 
