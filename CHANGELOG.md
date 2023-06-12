# Changelog

## [0.3.1] - 2023-06-12

### Change

- 不再通过将 `r.RemoteAddr` 设为空字符串来去除 `X-Forwarded-For` 请求头, 而是使用 omit 方式来去除. (omit 方式: 将 `X-Forwarded-For` 设为 nil)

## [0.3.0] - 2023-06-10

### Change

- 跟随 `xhe v0.0.6` 的配置模式

## [0.2.1] - 2023-06-07

### Fix

- 为浏览器默认加载的 wasm 链接添加版本号, 修复 v0.0.1/dist/xhe.umd.js 加载 v0.2.0/xhe-wc.wasm 的情况

## [0.2.0] - 2023-06-07

### Add

- 支持使用 `js.fetch:*` 头来设置 fetch options

## [0.1.2] - 2023-06-07

### Change

- 禁用反向代理的自动添加`X-Forwarded-For`请求头

## [0.1.1] - 2023-06-03

### Change

- 添加`application/json`响应头, 方便 http 客户端自动解码响应

## [0.1.0] - 2023-06-03

### Add

- 添加 HandleEval 方便远程操作

## [0.0.7] - 2023-06-02

### Fix Change

- 不再将 wasm 文件转成 data url base64 导入, 因为过于浪费内存. 具体是由 30M 内存飙升到 1G, 在浏览器端更是飙升至 3G, 因此不再转成 data url 导入

## [0.0.6] - 2023-06-02

### Add

- 添加 umd 打包, 方便直接编写 userscript 来连接 Xhe Wireguard

## [0.0.4] - 2023-05-31

### Change

- XheConnectInit 变更为函数, 支持设置 wasmUrl

## [0.0.2] - 2023-05-30

### Fix

- 修复 ListenTCP 带有端口参数时错误退出的问题

## [0.0.1] - 2023-05-30

### 好耶

- 第一个版本
