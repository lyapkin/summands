package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const ROWS = 1000000
const BATCH = 1000000
const SHEETS = 5

func main() {
	app := app.NewWithID("fund_summands_of_given_sum")
	window := app.NewWindow("Поиск комбинаций суммы")
	window.Resize(fyne.NewSize(700, 300))
	window.SetFixedSize(true)

	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	curPath := filepath.Dir(ex)
	log.Print(curPath)
	
	sumLabel := widget.NewLabel("Введите сумму")
	sumInput := newNumericalEntry()
	sumInput.Validator = func(val string) error {
		num, _ := strconv.Atoi(val)
		if num <= 0 {
			return errors.New("число должно быть больше чем 0")
		}
		return nil
	}

	summandsNumLabel := widget.NewLabel("Введите кол-во слагаемых")
	summandsNumInput := newNumericalEntry()
	summandsNumInput.Validator = func(val string) error {
		num, _ := strconv.Atoi(val)
		if num <= 0 {
			return errors.New("число должно быть больше чем 0")
		}
		if num > 7 {
			return errors.New("число должно быть НЕ больше чем 7")
		}
		return nil
	}

	ubLabel := widget.NewLabel("Введите значение макс. слагаемого")
	ubInput := newNumericalEntry()

	progress := widget.NewProgressBarInfinite()
	progress.Hide()
	progressText := widget.NewLabel("Найдено комбинаций: 0")

	pathLabel := widget.NewLabel(curPath)
	pathPicker := widget.NewButton("Выбрать папку", func() {
		onChosen := func(uri fyne.ListableURI, err error) {
			if err != nil {
				log.Print(err)
				return
			}
			if uri == nil {
				log.Print("uri is nill")
				return
			}
			pathLabel.SetText(uri.Path())
			curPath = uri.Path()
		}
		dialog.NewFolderOpen(onChosen, window).Show()
	})

	btn := widget.NewButton("Найти комбинации", nil)
	btn.Disable()
	btn.OnTapped = func() {
		progress.Show()
		progress.Start()
		progress.Refresh()
		progressText.SetText("Найдено комбинаций: 0")

		sumInput.Disable()
		summandsNumInput.Disable()
		ubInput.Disable()
		btn.Disable()
		pathPicker.Disable()

		sum, _ := strconv.Atoi(sumInput.Text)
		num, _ := strconv.Atoi(summandsNumInput.Text)
		ub, _ := strconv.Atoi(ubInput.Text)
		if ub > sum || ub <= 0 {
			ub = 0
		}

		go func() {
			comb := NewCombs(sum, num, ub, curPath, func(val int) {
				fyne.Do(func() {progressText.SetText(fmt.Sprintf("Найдено комбинаций: %v", val))})
			})
			start := time.Now()
			_, numCombs := comb.FindCombs()
			elapsed := time.Since(start)
			fyne.Do(func() {
				progressText.SetText(fmt.Sprintf("Найдено комбинаций: %d; Время выполнения: %d секунд", numCombs, int(elapsed.Seconds())))
				sumInput.Enable()
				summandsNumInput.Enable()
				ubInput.Enable()
				btn.Enable()
				pathPicker.Enable()
				progress.Stop()
				progress.Hide()
				progress.Refresh()
			})
			comb = nil
			runtime.GC()
		}()
	}

	validateForm := func() bool {
		if err := sumInput.Validate(); err != nil {
			return false
		}
		if err := summandsNumInput.Validate(); err != nil {
			return false
		}
		return true
	}

	updateBtnState := func() {
		if validateForm() {
			btn.Enable()
		} else {
			btn.Disable()
		}
	}

	sumInput.OnChanged = func(_ string) {
		updateBtnState()
	}
	summandsNumInput.OnChanged = func(_ string) {
		updateBtnState()
	}

	
	
	pathContainer := container.New(layout.NewBorderLayout(nil, nil, nil, pathPicker), pathPicker, pathLabel)

	inputContainer := container.New(layout.NewVBoxLayout(),
		sumLabel, sumInput,
		summandsNumLabel, summandsNumInput,
		ubLabel, ubInput,
	)
	content := container.New(layout.NewCustomPaddedVBoxLayout(32),
		inputContainer,
		pathContainer,
		container.NewGridWithRows(2,
			progressText,
			progress,
		),
		btn,
	)
	
	window.SetContent(content)
	window.ShowAndRun()
	
}

type numericalEntry struct {
	widget.Entry
}

func (e *numericalEntry) TypedRune(r rune) {
	if r >= '0' && r <= '9' {
		e.Entry.TypedRune(r)
	}
}

func newNumericalEntry() *numericalEntry {
	entry := &numericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}