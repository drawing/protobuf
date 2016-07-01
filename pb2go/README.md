
# 功能

protobuf 生成 Go 文件如果每个`.proto`中的包名都不同，并且互相依赖，生成`.go`文件后需要重建`import`和目录组织，使用此工具可以自动重建。

# 使用

首先需要安装此工具，通过简单的命令便可安装完成：

```
go get -u github.com/drawing/protobuf/protoc-gen-go
go get -u github.com/drawing/protobuf/pb2go
```

把所有proto文件放入同一个目录：

```
im@im:~/ $ ls -l pbfiles/
total 49
-rwxrwxrwx 1 root root 18558  6月 15  2015 1.proto
-rwxrwxrwx 1 root root 24146  6月 15  2015 2.proto
-rwxrwxrwx 1 root root   298  6月 15  2015 3.proto
```

执行生成命令，指定源路径和生成的目标路径：

```
pb2go -in pbfiles/ -out .
```

执行完成后，便在目标路径根据包名重建目录，使用时直接`import`对应包即可。
