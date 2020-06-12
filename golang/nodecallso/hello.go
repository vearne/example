package main // 这个文件一定要在main包下面

import "C" // 这个 import 也是必须的，有了这个才能生成 .h 文件
// 下面这一行不是注释，是导出为SO库的标准写法，注意 export前面不能有空格！！！
//export hello
func hello(value string)*C.char { // 如果函数有返回值，则要将返回值转换为C语言对应的类型
    return C.CString("hello " + value)
}
func main(){
    // 此处一定要有main函数，有main函数才能让cgo编译器去把包编译成C的库
}
