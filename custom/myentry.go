package custom

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type MultilineEdit struct {
	widget.Entry
	OnEnter func()
}

func (t *MultilineEdit) TypedKey(key *fyne.KeyEvent) {
	if key.Name == "Return" {

		if t.OnEnter != nil {
			t.OnEnter()
		}

	} else {

		t.Entry.TypedKey(key)
	}
}

func NewMultilineEdit() *MultilineEdit {
	edit := &MultilineEdit{}
	edit.Entry.MultiLine = true
	edit.Entry.Wrapping = fyne.TextTruncate
	//edit := &Entry{MultiLine: true, Wrapping: fyne.TextTruncate}
	edit.SetMinRowsVisible(1)
	edit.ExtendBaseWidget(edit)
	return edit
}
