package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

type File interface {
	Write(*[][]int)
	Close()
	Save()
}

func NewFile(name string, dir string) File {
	return &_File{name: name, dir: fmt.Sprintf("%v/%v", dir, name)}
}

type _File struct {
	name string
	dir string
	file *excelize.File
	filledRows int
	filledSheets int
	filledFiles int
	sw *excelize.StreamWriter
}

func (f *_File) Write(results *[][]int) {
	f.start()

	// log.Printf("writing from %v %v %v", f.filledRows, f.filledSheets, f.filledFiles)
	for i, row := range *results {
		f.writeRow(row)

		if f.isRowsFilled() && i < len(*results)-1 {
			// log.Printf("stop writing on %v %v %v", f.filledRows, f.filledSheets, f.filledFiles)
			f.end()
			f.start()
			// log.Printf("writing from %v %v %v", f.filledRows, f.filledSheets, f.filledFiles)
		}
	}
	// log.Printf("stop writing on %v %v %v", f.filledRows, f.filledSheets, f.filledFiles)

	f.end()
}

func (f *_File) start() {
	if f.file == nil {
		f.newFile()
	}

	if f.sw == nil {
		f.newSheet()
	}
}

func (f *_File) end() {
	if f.isRowsFilled() {
		// log.Print("rows filled")
		f.closeSheet()
	}

	if f.isSheetsFilled() {
		// log.Print("sheets filled")
		f.closeFile()
	}
	
}

func (f *_File) Save() {
	// log.Print("Save")
	if err := f.sw.Flush(); err != nil {
		f.err("Error flushing StreamWriter:", err)
	}
	f.saveFile()
}

func (f *_File) saveFile() {
	err := os.MkdirAll(f.dir, os.ModePerm)
	if err != nil {
		f.err("Error creating directory: %v\n", err)
	}

	if err := f.file.SaveAs(fmt.Sprintf("%v/%v", f.dir, f.getFileName())); err != nil {
		f.err("Error saving Excel file:", err)
	}
}

func (f *_File) Close() {
	// log.Print("Close")
	var err error;
	if f.sw != nil {
		err = f.sw.Flush()
	}
	if err != nil {
		f.err("Error flushing StreamWriter", err)
	}
	f.sw = nil

	if f.file != nil {
		err = f.file.Close()
	}
	if err != nil {
		f.err("Error closing file", err)
	}
	f.file = nil
}

func (f *_File) newFile() {
	// log.Printf("newFile %v", f.getFileName())
	if f.file != nil {
		f.closeFile()
	}

	f.filledSheets = 0
	f.filledRows = 0
	f.file = excelize.NewFile()

	if err := f.file.SetSheetName(
		f.file.GetSheetName(f.filledSheets),
		f.getCurrentSheetName(),
	); err != nil {
		f.err("Error setting sheet name:",err)
	}

	if sw, err := f.file.NewStreamWriter(f.getCurrentSheetName()); err != nil {
		f.err("Error creating StreamWriter:",err)
	} else {
		f.sw = sw
	}
}

func (f *_File) closeFile() {
	// log.Printf("close file %v", f.getFileName())
	if f.file == nil {
		return
	}

	f.saveFile()

	if err := f.file.Close(); err != nil {
		f.err("Error closing Excel file:", err)
	}
	f.filledFiles++
	f.file = nil
}

func (f *_File) writeRow(row []int) {
	cell, err := excelize.CoordinatesToCellName(1, f.filledRows+1)
	if err != nil {
		f.err("Error getting cell name: ", err)
	}
	vals := make([]interface{}, len(row))
	for i, val := range row {
		vals[i] = val
	}
	
	if err := f.sw.SetRow(cell, vals); err != nil {
		f.err("Error writing to stream:", err)
	}
	f.filledRows++
}

func (f *_File) newSheet() {
	// log.Print("newSheet")
	if _, err := f.file.NewSheet(f.getCurrentSheetName()); err != nil {
		f.err("Error creating new sheet:", err)
	}

	if sw, err := f.file.NewStreamWriter(f.getCurrentSheetName()); err != nil {
		f.err("Error creating StreamWriter:",err)
	} else {
		f.sw = sw
	}

	f.filledRows = 0

}

func (f *_File) closeSheet() {
	// log.Print("close sheet")
	if f.sw == nil {
		return
	}

	if err := f.sw.Flush(); err != nil {
		f.err("Error flushing StreamWriter:", err)
	}
	f.sw = nil
	f.filledSheets++
}

func (f *_File) getFileName() string {
	return fmt.Sprintf("%v-%v.xlsx", f.name, f.filledFiles)
}

func (f *_File) isRowsFilled() bool {
	return f.filledRows >= ROWS
}

func (f *_File) isSheetsFilled() bool {
	return f.filledSheets >= SHEETS
}

func (f *_File) getCurrentSheetName() string {
	return fmt.Sprintf("Таблица %d", f.filledSheets+1)
}


func (f *_File) err(msg string, err error)  {
	if f.sw != nil {
		if err := f.sw.Flush(); err != nil {
			log.Println("Error flushing StreamWriter:", err)
		}
	}

	if f.file != nil {
		if err := f.file.Close(); err != nil {
			log.Println("Error closing Excel file:", err)
		}
	}
	
	log.Fatal(msg, err)
}