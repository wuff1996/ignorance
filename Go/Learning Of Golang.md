# Learning Of Golang



##		What for:	I used to write diary about my daily learning of Golang,but as the time goes by,the problem of looking for what I want to review has become larger. So I decide to restruct what I have learned of Golang,or any other knowledge e.g. Linux, Git, Docker... By the way,the reason of writing it in English is that I think English is more primitive.



### Struct:

![image-20200929144915037](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20200929144915037.png)

---------------



#### HelloWorld:

```go
 //声明这个文件属于哪一个包，而main包用来定义一个可独立执行的函数，即main包中的main函数为程序入口。
package main 

//import=bring in 即导入 " fmt" 包， fmt包用来格式化输出和扫描输入
import "fmt"  

 //main函数为程序入口
func main () { 
    //Println函数为fmt包中的基本输出函数，它会在末尾自动加上换行符。如需打印多个参数，用 "," 隔开即可.
	fmt.Println("Hello World !") 
}
```



------------------------------



#### Go mod:

[URL]: https://blog.golang.org/using-go-modules	"Go mod"

Go mod是**Golangv1.12**版本推出的,用来管理go **module(各种go package的集合体)**。它允许我们在**非GOPATH目录**下随意创建项目。增加了依赖的灵活性。其使用方法为：

在**非GOPATH目录**下

```shell
go mod init  github.com/wuff1996/ignorance/Golang/goMod
touch main.go
```

值得注意的是，我们需要确实GO111MODULE变量是否为off

```
go env
windows: set GO111MODULE=auto || linux: export GO111MODULE=auto
或者
vim ~/.bashrc
在末尾加上：
export GO111MODULE=auto
然后wq
source ~/.bashrc 执行初始化文件。
```

##### Problem:

https://github.com/wuff1996/ignorance/tree/master/Go/problem

-----------------------



#### 程序结构：

 		大程序都是从小的基本组件构建而来的：变量存储值；简单表达式通过加和减合并成大的；基本类型通过数组和结构体进行聚合；表达式通过if、else if、else、for、switch、select控制语句来决定执行顺序；语句被组织成函数用于隔离和复用；函数被组织成源文件按和包。

1.**名称**严格区分大小写。//**[a-zA-Z]\w***     

2.**关键字**只能在语法中使用，**不能作为名称**，关键字如下：**break,default,func,interface,select,case,defer,go,map,struct,const,fallthrough,if,range,type,continue,for,import,return,var**.     

3.**程序实体的第一个字母大小写决定其可见性是否跨包**，即是否可被外界可见，可访问，可引用。但包本身由小写字母表示。 

4.**声明**：给一个程序实体命名，设定其部分或全部属性。。有四个主要的声明：var, const, type, func.

 5.**var**声明通常时为了那些初始化表达式类型不一致的局部变量。**短变量声明不需要声明所有在左边的变量**,对于已经在词法块中声明的变量，短边量的行为等同于赋值，如

```go
package main

import "fmt"

func main() {
	fmt.Println(R())
}

func R()(s string){
	s = "-----------"
	i,s := R2()
	fmt.Println(i)
	return
}

func R2()(int,string){
	return 0,"+++++++++++"
}
/*result: 
0
+++++++++++
*/

```

6.**指针**：指针的值是一个变量地址，指向变量所在的内存空间，而变量存储值，变量出现在左边是赋值，右边是读取值。每一个聚合类型变量的组成（结构体的成员或数组中的元素）都是变量，所以也有一个地址。  

**指针和引用的区别**：https://www.tutorialspoint.com/what-is-difference-between-a-pointer-and-reference-parameter-in-cplusplus

7.**new**:new(T)分配一块类型为T的内存空间，并返回其地址，使得变量不需要变量名，指针可以 初始化任何类型的零值。

**make**： make(T)，创建一个未命名的T类型，并进行初始化，T类型只能为slice、map、channel，并返回其引用。

8.**变量的声明周期**：变量如果不可访问意味着它不会影响其他的计算过程则将被回收，编译器可以选择使用堆或栈来分配变量的空间，而不是基于关键字，如

```go
var global *int
func f () {
	var x = 1
    //由于可以通过全局变量 global 访问到局部变量x，所以x从f中逃逸，每一次的变量逃逸都需要一次额外的内存分配过程。
	global = &x
}
```

