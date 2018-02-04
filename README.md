# goRemoveSame
搜索当前目录以及子目录所有文件，通过计算hash值判断文件是否相同，通过web界面预览重复文件，选择需要删除的文件，然后可以提交删除。
> jquery通过ajax与web服务器通讯，vuejs绑定数据操作。  
> 界面没有美化，比较丑，能用就可以了。

# 安装
```shell
$ go get -v github.com/ohko/goRemoveSame
$ go install github.com/ohko/goRemoveSame
```

# 使用
用过命令行启动程序，然后通过浏览器访问 [http://127.0.0.1:8080]()。
```
$ goRemoveSame
2018/02/04 22:04:09 开始分析文件...
2018/02/04 22:04:09 [1/6] 分析: ./bindata_assetfs.go
2018/02/04 22:04:09 [2/6] 分析: ./install.sh
2018/02/04 22:04:09 [3/6] 分析: ./main.go
2018/02/04 22:04:09 [4/6] 分析: ./static/index.html
2018/02/04 22:04:09 [5/6] 分析: ./static/jquery.min.js
2018/02/04 22:04:09 [6/6] 分析: ./static/vue.min.js
2018/02/04 22:04:09 分析完成，可刷新浏览器了，文件总数: 6 个 / 重复文件： 0 组
2018/02/04 22:04:09 Listen http:// :8080
```

# 功能
- 搜索相同文件
- 通过浏览器预览
- 执行删除选中文件