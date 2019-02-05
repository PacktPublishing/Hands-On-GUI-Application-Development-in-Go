package main

import (
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/gesture"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/exp/shiny/widget"
	"golang.org/x/exp/shiny/widget/node"
	"golang.org/x/exp/shiny/widget/theme"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"image"
	"log"
)

const buttonPad = 4

type button struct {
	node.LeafEmbed

	label   string
	onClick func()
}

func newButton(label string, onClick func()) *button {
	b := &button{label: label, onClick: onClick}
	b.Wrapper = b

	return b
}

func (b *button) Measure(t *theme.Theme, widthHint, heightHint int) {
	face := t.AcquireFontFace(theme.FontFaceOptions{})
	defer t.ReleaseFontFace(theme.FontFaceOptions{}, face)

	b.MeasuredSize.X = font.MeasureString(face, b.label).Ceil() + 2*buttonPad
	b.MeasuredSize.Y = face.Metrics().Ascent.Ceil() + face.Metrics().Descent.Ceil() + 2*buttonPad
}

func (b *button) PaintBase(ctx *node.PaintBaseContext, origin image.Point) error {
	b.Marks.UnmarkNeedsPaintBase()
	face := ctx.Theme.AcquireFontFace(theme.FontFaceOptions{})
	defer ctx.Theme.ReleaseFontFace(theme.FontFaceOptions{}, face)

	draw.Draw(ctx.Dst, b.Rect.Add(origin).Inset(buttonPad), theme.Foreground.Uniform(ctx.Theme), image.Point{}, draw.Src)
	d := font.Drawer{
		Dst:  ctx.Dst,
		Src:  theme.Background.Uniform(ctx.Theme),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(b.Rect.Min.X + buttonPad), Y: fixed.I(b.Rect.Min.Y + face.Metrics().Ascent.Ceil() + buttonPad)},
	}
	d.DrawString(b.label)

	return nil
}

func (b *button) OnInputEvent(e interface{}, origin image.Point) node.EventHandled {
	if ev, ok := e.(gesture.Event); ok {
		if ev.Type == gesture.TypeTap && b.onClick != nil {
			b.onClick()
		}

		return node.Handled
	}

	return node.NotHandled
}

func main() {
	driver.Main(func(s screen.Screen) {
		label := widget.NewLabel("Hello World!")
		button := newButton("Quit",
			func() {
				log.Println("To quit close this window")
			})

		w := widget.NewFlow(widget.AxisVertical, label, button)
		sheet := widget.NewSheet(widget.NewUniform(theme.Neutral, w))

		w.Measure(theme.Default, 0, 0)
		if err := widget.RunWindow(s, sheet, &widget.RunWindowOptions{
			NewWindowOptions: screen.NewWindowOptions{
				Title:  "Hello",
				Width:  w.MeasuredSize.X,
				Height: w.MeasuredSize.Y,
			},
		}); err != nil {
			log.Fatal(err)
		}
	})
}
