## 构建命令

- 构建Mac平台
```shell
go build -buildmode=pie -v -ldflags '-w -s ' -o ./build/callback_metersphere callback_metersphere.go
```

- 构建linux平台
```shell
CGO_ENABLED=0 GOOS=linux  GOARCH=amd64  CC=x86_64-linux-musl-gcc  CXX=x86_64-linux-musl-g++  go build -buildmode=pie -o ./build/callback_metersphere -v -a -ldflags "-s -w" callback_metersphere.go
```
- Windows支持的构建
```shell
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc fyne package -os windows -icon static/icon/logo.png --id org.cyz.tools.callback_metersphere --release
```

- 构建Windows平台
```shell
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -installsuffix cgo -buildmode=pie -v -ldflags '-w -s -H windowsgui' -o ./build/callback_metersphere.exe callback_metersphere.go
```

**有图标的**
```shell
go build -installsuffix cgo -ldflags="-H windowsgui -w -s" -o ./build/callback_metersphere.exe callback_metersphere.go
```

- 构建Android平台包
```shell
fyne package -os android -appID org.cyz.tools.callback_metersphere --appVersion=0.0.1 -icon ./static/icon/logo.png
```

- 构建IOS平台
 ```shell
fyne package -os ios - appID org.cyz.tools.callback_metersphere -icon ./assets/device.png
```