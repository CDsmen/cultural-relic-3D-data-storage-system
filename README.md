### 简介
  本项目用于简单搭建一个三维数据存储系统，通过grpc和docker保证了项目的可拓展性和可移植性。
### 使用方式
  在确保安装了docker与docker-compose的情况下，在项目目录中执行以下命令即可：
  ```bash
  docker-compose up
  ```
### 测试
  可通过./client_test文件夹下的测试文件进行简单测试：
  ```bash
  go run ./client_test/test.go
  ```
  也可以通过grpc利用./idl下的artifact_svr.proto文件进行不同语言的代码生成，进行接入测试。
### 拓展
  该项目完成了基础框架搭建，对后续接入远端文件系统、支持更多格式输入、提供更多信息存储等具体需要，欢迎为本项目提出贡献。
