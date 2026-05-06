package main

import (
	_ "embed" // 引入 embed 模块打包静态资源
	"net/url"
	"strings"
	"time"
	"unsafe"

	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
	"golang.org/x/sys/windows"
)

// --- 将图标文件嵌入到程序二进制中 ---
//
//go:embed icon.ico
var iconData []byte

// ------------------------------------------

var (
	moduser32                    = windows.NewLazySystemDLL("user32.dll")
	procGetClipboardOwner        = moduser32.NewProc("GetClipboardOwner")
	procGetWindowThreadProcessId = moduser32.NewProc("GetWindowThreadProcessId")
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	// --- 托盘图标 ---
	systray.SetIcon(iconData) // 直接使用嵌入的图标数据
	// ---------------------------

	systray.SetTitle("Chrome 链接净化器")
	systray.SetTooltip("正在监控 Chrome 地址栏复制...")

	mQuit := systray.AddMenuItem("退出程序", "停止运行并关闭")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	go clipboardLoop()
}

func onExit() {
	// 程序退出时的清理工作（此处可为空）
}

func clipboardLoop() {
	var lastText string
	for {
		text, err := clipboard.ReadAll()
		if err == nil && text != lastText && text != "" {
			// A. 识别来源是否为 Chrome
			if getClipboardOwnerProcessName() == "chrome.exe" {
				// B. 核心处理逻辑
				newText, changed := processURL(text)
				if changed {
					clipboard.WriteAll(newText)
					lastText = newText
				} else {
					lastText = text
				}
			} else {
				lastText = text
			}
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

// 核心逻辑：识别网址是否有后缀
func processURL(rawUrl string) (string, bool) {
	trimmed := strings.TrimSpace(rawUrl)
	// 必须以协议头开头
	if !strings.HasPrefix(trimmed, "http://") && !strings.HasPrefix(trimmed, "https://") {
		return rawUrl, false
	}

	// 使用解析器分析 URL
	u, err := url.Parse(trimmed)
	if err != nil {
		return rawUrl, false
	}

	// 规则判定：
	// 1. Path 必须为空，或者仅仅是一个 "/"
	// 2. Query (参数，即 ? 后的内容) 必须为空
	// 3. Fragment (锚点，即 # 后的内容) 必须为空
	if (u.Path == "" || u.Path == "/") && u.RawQuery == "" && u.Fragment == "" {
		// 满足条件：这是一个纯域名地址，去掉前缀
		// 返回 Host (如 www.baidu.com)
		return u.Host, true
	}

	// 不满足条件：说明有后缀（如 /user/123 或 ?id=1），保留原样
	return rawUrl, false
}

// 获取剪贴板所有者进程名（保留之前的稳定逻辑）
func getClipboardOwnerProcessName() string {
	hwnd, _, _ := procGetClipboardOwner.Call()
	if hwnd == 0 {
		return ""
	}
	var pid uint32
	procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(handle)
	var n uint32 = 1024
	buf := make([]uint16, n)
	_ = windows.QueryFullProcessImageName(handle, 0, &buf[0], &n)
	fullPath := windows.UTF16ToString(buf[:n])
	parts := strings.Split(fullPath, "\\")
	return strings.ToLower(parts[len(parts)-1])
}
