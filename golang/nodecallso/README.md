## nodejs 调用动态链接库

### 生成动态链接库
go build -buildmode=c-shared -o hello.so hello.go   
此步骤，会生成对应的头文件`hello.h`

Mac下动态链接库需要命名为`hello.dylib`

### 安装ffi(需要python2)
```
npm install -g ffi
```
### 调用动态链接库
```
node test.js
```
