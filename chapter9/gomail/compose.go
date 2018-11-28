package main

import (
	"github.com/golang-ui/nuklear/nk"
)

type composeUI struct {
	main *mainUI

	composeSubject []byte
	composeEmail   []byte
	composeContent []byte
}

func newComposeUI(main *mainUI) *composeUI {
	c := &composeUI{main,
		make([]byte, 512, 512), make([]byte, 512, 512),
		make([]byte, 4096, 4096)}

	copy(c.composeSubject[:], "subject")
	copy(c.composeEmail[:], "email")
	copy(c.composeContent[:], "content")

	return c
}

func (c *composeUI) drawLayout(ctx *nk.Context, height int) {
	nk.NkLayoutRowDynamic(ctx, 0, 1)
	nk.NkEditStringZeroTerminated(ctx, nk.EditBox|nk.EditSelectable|nk.EditClipboard,
		c.composeSubject, int32(len(c.composeSubject)), nil)
	nk.NkLayoutRowTemplateBegin(ctx, 0)
	nk.NkLayoutRowTemplatePushStatic(ctx, 25)
	nk.NkLayoutRowTemplatePushVariable(ctx, 320)
	nk.NkLayoutRowTemplateEnd(ctx)
	nk.NkLabel(ctx, "To", nk.TextAlignRight)
	nk.NkEditStringZeroTerminated(ctx, nk.EditBox|nk.EditSelectable|nk.EditClipboard,
		c.composeEmail, int32(len(c.composeEmail)), nil)
	nk.NkLayoutRowDynamic(ctx, float32(height-114), 1)
	nk.NkEditStringZeroTerminated(ctx, nk.EditBox|nk.EditSelectable|nk.EditClipboard,
		c.composeContent, int32(len(c.composeContent)), nil)

	nk.NkLayoutRowTemplateBegin(ctx, 0)
	nk.NkLayoutRowTemplatePushVariable(ctx, 234)
	nk.NkLayoutRowTemplatePushStatic(ctx, 64)
	nk.NkLayoutRowTemplatePushStatic(ctx, 64)
	nk.NkLayoutRowTemplateEnd(ctx)
	nk.NkLabel(ctx, "", nk.TextAlignLeft)
	if nk.NkButtonLabel(ctx, "Cancel") > 0 {
		c.main.compose = nil
	}
	nk.NkButtonLabel(ctx, "Send")
}
