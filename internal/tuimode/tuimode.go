package tuimode

import (
	"code.rocketnine.space/tslocum/cview"
)

func App() {
	// Create a TUI app
	app := cview.NewApplication()
	defer app.HandlePanic()
	app.EnableMouse(true)

	// Generate flex grid
	flex := cview.NewFlex()
	flex.SetDirection(cview.FlexRow)

	tv := cview.NewTextView()
	tv.SetBorder(true)
	tv.SetTitle("Hello, world!")
	tv.SetText("Lorem ipsum dolor sit amet")
	flex.AddItem(tv, 0, 0, true)
	flex.AddItem(tv, 0, 1, false)

	app.SetRoot(flex, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
