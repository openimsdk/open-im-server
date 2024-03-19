# go doc 

Go语言项目十分重视代码的文档，在软件设计中，文档对于软件的可维护和易使用具有重大的影响。因此，文档必须是书写良好并准确的，与此同时它还需要易于书写和维护。


**OpenIM 支持的命令

Go语言注释
Go语言中注释一般分为两种，分别是单行注释和多行注释

单行注释是以 // 开头的注释，可以出现在任何地方。
多行注释也叫块注释，以 /* 开头，以 */ 结尾，不可以嵌套使用，一般用于包的文档描述或注释成块的代码片段。
每一个 package 都应该有相关注释，在 package 语句之前的注释内容将被默认认为是这个包的文档， package 的注释应该提供一些相关信息并对整体功能做简要的介绍。

在日常开发过程中，可以使用go doc和godoc命令生成代码的文档。



go doc
go doc 命令打印Go语言程序实体上的文档。可以使用参数来指定程序实体的标识符。

Go语言程序实体是指变量、常量、函数、结构体以及接口。

程序实体标识符就是程序实体的名称。

输出指定 package ，指定类型，指定方法的注释

$ go doc sync.WaitGroup.Add
输出指定 package ，指定类型的所有程序实体，包括未导出的

$ go doc -u -all sync.WaitGroup

输出指定 package 的所有程序实体（非所有详细注释）

$ go doc -u sync


godoc
godoc命令主要用于在无法联网的环境下，以web形式，查看Go语言标准库和项目依赖库的文档。

在 go 1.12 之后的版本中，godoc不再做为go编译器的一部分存在。依然可以通过go get命令安装：

go get -u -v golang.org/x/tools/cmd/godoc
国内的安装方法


Explain
mkdir -p $GOPATH/src/golang.org/x
cd $GOPATH/src/golang.org/x
git clone https://github.com/golang/tools.git
cd tools/cmd/godoc
go install 
ls -alh $GOPATH/bin
通过终端查看文档

go doc命令

$ go doc help
usage: go doc [-u] [-c] [package|[package.]symbol[.method]]
可以看到，go doc接受的参数，可以是包名，也可以是包里的结构、方法等，默认为显示当前目录下的文档。


通过网页查看文档

godoc命令

$ godoc -http=:6060
godoc会监听6060端口，通过网页访问 http://127.0.0.1:6060，godoc基于GOROOT和GOPATH路径下的代码生成文档的。打开首页如下，我们自己项目工程文档和通过go get的代码文档都在Packages中的Third party里面。



编写自己的文档

1、设计接口函数代码

创建documents/calc.go文件


编写文档规则

1、文档中显示的详细主体内容，大多是由用户注释部分提供，注释的方式有两种，单行注释"//"和代码块"/* */"注释。

2、在源码文件中，在package语句前做注释，在文档中看到的就是Overview部分， 注意：此注释必须紧挨package语句前一行，要作为Overview部分的，注释块中间不能有空行。

3、在函数、结构、变量等前做注释的，在文档中看到的就是该项详细描述。注释规则同上。

4、编写的Example程序，函数名必须以Example为前缀，可将测试的输出结果放在在函数尾部，以"// Output:"另起一行，然后将输出内容注释，并追加在后面。