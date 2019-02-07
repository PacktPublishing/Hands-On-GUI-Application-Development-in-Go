package main

import (
	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/client"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/golang-ui/nuklear/nk"
)

func init() {
	runtime.LockOSThread()
}

type mainUI struct {
	server  client.EmailServer
	current *client.EmailMessage

	compose *composeUI
}

func main() {
	glfw.Init()
	win, _ := glfw.CreateWindow(600, 400, "GoMail", nil, nil)
	win.SetSizeLimits(600, 400, glfw.DontCare, glfw.DontCare)
	win.MakeContextCurrent()
	gl.Init()

	ctx := nk.NkPlatformInit(win, nk.PlatformInstallCallbacks)
	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)
	font := nk.NkFontAtlasAddDefault(atlas, 14, nil)
	nk.NkFontStashEnd()
	nk.NkStyleSetFont(ctx, font.Handle())

	server := client.NewTestServer()
	mainUI := &mainUI{server: server, current: server.CurrentMessage()}

	quit := make(chan struct{}, 1)
	ticker := time.NewTicker(time.Second / 30)
	for {
		select {
		case <-quit:
			nk.NkPlatformShutdown()
			glfw.Terminate()
			ticker.Stop()
			return
		case <-ticker.C:
			if win.ShouldClose() {
				close(quit)
				continue
			}
			mainUI.draw(win, ctx)
			win.SwapBuffers()
			glfw.PollEvents()
		}
	}
}

func (ui *mainUI) drawCompose(ctx *nk.Context) {
	bounds := nk.NkRect(20, 20, 400, 320)
	update := nk.NkBegin(ctx, "Compose", bounds, nk.WindowNoScrollbar|nk.WindowBorder|nk.WindowTitle|nk.WindowMovable|nk.WindowMinimizable)

	if update > 0 {
		ui.compose.drawLayout(ctx, 296)
	}

	nk.NkEnd(ctx)
}

func (ui *mainUI) drawLayout(win *glfw.Window, ctx *nk.Context, height int) {
	nk.NkMenubarBegin(ctx)
	nk.NkLayoutRowBegin(ctx, nk.LayoutDynamicRow, 25, 3)
	nk.NkLayoutRowPush(ctx, 45)
	if nk.NkMenuBeginLabel(ctx, "File", nk.TextAlignLeft, nk.NkVec2(120, 200)) > 0 {
		nk.NkLayoutRowDynamic(ctx, 25, 1)
		if nk.NkMenuItemLabel(ctx, "New", nk.TextAlignLeft) > 0 {
			ui.compose = newComposeUI(ui)
		}
		if nk.NkMenuItemLabel(ctx, "Quit", nk.TextAlignLeft) > 0 {
			win.SetShouldClose(true)
		}

		nk.NkMenuEnd(ctx)
	}
	if nk.NkMenuBeginLabel(ctx, "Edit", nk.TextAlignLeft, nk.NkVec2(120, 200)) > 0 {
		nk.NkLayoutRowDynamic(ctx, 25, 1)
		nk.NkMenuItemLabel(ctx, "Delete", nk.TextAlignLeft)
		nk.NkMenuItemLabel(ctx, "Cut", nk.TextAlignLeft)
		nk.NkMenuItemLabel(ctx, "Copy", nk.TextAlignLeft)
		nk.NkMenuItemLabel(ctx, "Paste", nk.TextAlignLeft)

		nk.NkMenuEnd(ctx)
	}
	if nk.NkMenuBeginLabel(ctx, "Help", nk.TextAlignLeft, nk.NkVec2(120, 200)) > 0 {
		nk.NkLayoutRowDynamic(ctx, 25, 1)
		nk.NkMenuItemLabel(ctx, "About", nk.TextAlignLeft)

		nk.NkMenuEnd(ctx)
	}
	nk.NkMenubarEnd(ctx)

	toolbarHeight := float32(24)
	nk.NkLayoutRowStatic(ctx, toolbarHeight, 78, 7)
	if nk.NkButtonLabel(ctx, "New") > 0 {
		ui.compose = newComposeUI(ui)
	}
	nk.NkButtonLabel(ctx, "Reply")
	nk.NkButtonLabel(ctx, "Reply All")

	nk.NkButtonLabel(ctx, "Delete")

	nk.NkButtonLabel(ctx, "Cut")
	nk.NkButtonLabel(ctx, "Copy")
	nk.NkButtonLabel(ctx, "Paste")

	nk.NkLayoutRowTemplateBegin(ctx, float32(height)-toolbarHeight)
	nk.NkLayoutRowTemplatePushStatic(ctx, 80)
	nk.NkLayoutRowTemplatePushVariable(ctx, 320)
	nk.NkLayoutRowTemplateEnd(ctx)

	nk.NkGroupBegin(ctx, "Inbox", 1)
	nk.NkLayoutRowDynamic(ctx, 0, 1)
	for _, email := range ui.server.ListMessages() {
		var selected int32
		if email == ui.current {
			selected = 1
		}
		if nk.NkSelectableLabel(ctx, email.Subject, nk.TextAlignLeft, &selected) > 0 {
			ui.current = email
		}
	}
	nk.NkGroupEnd(ctx)

	nk.NkGroupBegin(ctx, "Content", 1)
	nk.NkLayoutRowDynamic(ctx, 0, 1)
	nk.NkLabel(ctx, ui.current.Subject, nk.TextAlignLeft)
	nk.NkLayoutRowTemplateBegin(ctx, 0)
	nk.NkLayoutRowTemplatePushStatic(ctx, 50)
	nk.NkLayoutRowTemplatePushVariable(ctx, 320)
	nk.NkLayoutRowTemplateEnd(ctx)
	nk.NkLabel(ctx, "From", nk.TextAlignRight)
	nk.NkLabel(ctx, string(ui.current.From), nk.TextAlignLeft)
	nk.NkLabel(ctx, "To", nk.TextAlignRight)
	nk.NkLabel(ctx, string(ui.current.To), nk.TextAlignLeft)
	nk.NkLabel(ctx, "Date", nk.TextAlignRight)
	nk.NkLabel(ctx, ui.current.DateString(), nk.TextAlignLeft)
	nk.NkLayoutRowDynamic(ctx, 0, 1)
	nk.NkLabel(ctx, ui.current.Content, nk.TextAlignLeft)
	nk.NkGroupEnd(ctx)
}

func (ui *mainUI) draw(win *glfw.Window, ctx *nk.Context) {
	nk.NkPlatformNewFrame()
	width, height := win.GetSize()
	bounds := nk.NkRect(0, 0, float32(width), float32(height))
	update := nk.NkBegin(ctx, "gomail", bounds, nk.WindowNoScrollbar)

	if update > 0 {
		ui.drawLayout(win, ctx, height)
	}
	nk.NkEnd(ctx)

	if ui.compose != nil {
		ui.drawCompose(ctx)
	}

	// Draw to viewport
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(0x10, 0x10, 0x10, 0xff)
	nk.NkPlatformRender(nk.AntiAliasingOn, 512*1024, 128*1024)
}
