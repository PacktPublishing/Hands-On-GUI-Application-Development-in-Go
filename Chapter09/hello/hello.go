package main

import (
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/golang-ui/nuklear/nk"
)

const pad = 8

func init() {
	runtime.LockOSThread()
}

func main() {
	glfw.Init()
	win, _ := glfw.CreateWindow(120, 80, "Hello World", nil, nil)
	win.MakeContextCurrent()
	gl.Init()

	ctx := nk.NkPlatformInit(win, nk.PlatformInstallCallbacks)
	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)
	font := nk.NkFontAtlasAddDefault(atlas, 14, nil)
	nk.NkFontStashEnd()
	nk.NkStyleSetFont(ctx, font.Handle())

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
			draw(win, ctx)
			win.SwapBuffers()
			glfw.PollEvents()
		}
	}
}

func draw(win *glfw.Window, ctx *nk.Context) {
	// Define GUI
	nk.NkPlatformNewFrame()
	width, height := win.GetSize()
	bounds := nk.NkRect(0, 0, float32(width), float32(height))
	update := nk.NkBegin(ctx, "", bounds, nk.WindowNoScrollbar)

	if update > 0 {
		cellWidth := int32(width-pad*2)
		cellHeight := float32(height-pad*2) / 2.0
		nk.NkLayoutRowStatic(ctx, cellHeight, cellWidth, 1)
		{
			nk.NkLabel(ctx, "Hello World!", nk.TextCentered)
		}
		nk.NkLayoutRowStatic(ctx, cellHeight, cellWidth, 1)
		{
			if nk.NkButtonLabel(ctx, "Quit") > 0 {
				win.SetShouldClose(true)
			}
		}
	}
	nk.NkEnd(ctx)

	// Draw to viewport
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(0x10, 0x10, 0x10, 0xff)
	nk.NkPlatformRender(nk.AntiAliasingOn, 4096, 1024)
}
