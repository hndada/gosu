// Deprecated.
// func NewHintGlowImage(skin Skin, noteImage draws.Image) {
// 	for i := range skin.HintSprites {
// 		const (
// 			padScale   = 1.1
// 			outerScale = 1.2
// 		)
// 		sw, sh := noteImage.Size()
// 		outer := draws.NewImageScaled(noteImage, outerScale)
// 		pad := draws.NewImageScaled(noteImage, padScale)
// 		inner := noteImage
// 		a := uint8(255 * FieldOpaque)
// 		img := ebiten.NewImage(outer.Size())
// 		{
// 			op := &ebiten.DrawImageOptions{}
// 			op.ColorM.ScaleWithColor(color.NRGBA{128, 128, 128, a})
// 			op.GeoM.Translate(0, 0)
// 			img.DrawImage(outer, op)
// 		}
// 		{
// 			op := &ebiten.DrawImageOptions{}
// 			op.ColorM.ScaleWithColor(color.NRGBA{255, 255, 0, a})
// 			if i == 0 { // Blank for idle, Yellow for highlight.
// 				op.CompositeMode = ebiten.CompositeModeDestinationOut
// 			}
// 			sd := outerScale - padScale // Size difference.
// 			op.GeoM.Translate(sd/2*float64(sw), sd/2*float64(sh))
// 			img.DrawImage(pad, op)
// 		}
// 		{
// 			op := &ebiten.DrawImageOptions{}
// 			op.ColorM.ScaleWithColor(color.NRGBA{60, 60, 60, a})
// 			sd := outerScale - 1 // Size difference.
// 			op.GeoM.Translate(sd/2*float64(sw), sd/2*float64(sh))
// 			img.DrawImage(inner, op)
// 		}
// 		sprite := draws.NewSpriteFromSource(img)
// 		sprite.SetScaleToH(1.2 * regularNoteHeight)
// 		sprite.Locate(HitPosition, FieldPosition, draws.CenterMiddle)
// 		skin.HintSprites[i] = sprite
// 	}
// }
