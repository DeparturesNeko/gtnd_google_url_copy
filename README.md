# Chrome Link Cleaner | Chrome 链接净化器

![Platform](https://img.shields.io/badge/Platform-Windows-blue) ![Language](https://img.shields.io/badge/Language-Golang-00ADD8) ![License](https://img.shields.io/badge/License-GPL-green)

Chrome 浏览器在复制地址栏网址时强制添加 `http://` 或 `https://` 前缀，当你需要复制链接作为ping或者ssh地址的时候不想手动删？那就直接把他干掉吧。

---

## ✨ 核心特性

*   **识别环境**：仅当复制动作发生在 `chrome.exe` (Google Chrome) 内部时触发，不干扰记事本、IDE 其他软件。
*   **触发逻辑**：
    *   **纯域名**（如 `https://www.google.com/`）→ 自动净化为 `www.google.com`。
    *   **带路径链接**（如 `https://github.com/trending`）→ 会自动保留前缀，尽量不杀错。
*   **轻量化**：只包含任务栏图标，没有任何前台窗口，小工具就不整花里胡哨的了。

---

## 🛠️ 开发与编译指南

基于 **Golang** 开发，推荐使用 **VSCode** 作为开发环境。

### 1. 前置准备
*   安装 [Go 编程语言](https://golang.google.cn/dl/) (推荐 1.16 或以上版本)。
*   在 VSCode 中安装 **Go** 官方插件。
*   下载本仓库代码，并确保目录下有 `icon.ico` 图标文件。

### 2. 安装必要工具
在 VSCode 集成终端中运行以下命令，用于处理 Windows 资源文件：
```bash
go install [github.com/akavel/rsrc@latest](https://github.com/akavel/rsrc@latest)
go get [github.com/atotto/clipboard](https://github.com/atotto/clipboard)
go get [github.com/getlantern/systray](https://github.com/getlantern/systray)
go get golang.org/x/sys/windows
```

### 3. 执行编译
```bash
rsrc -arch amd64 -ico icon.ico -o rsrc.syso
go build -ldflags "-H windowsgui" -o ChromeLinkCleaner.exe main.go
```

### 4. RUN IT
执行exe，任务栏托盘会出现程序图标，不需要的时候右键关闭就可以咯
