package snippetlib

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"strconv"
)

const NONEVALUE = -9999

func WriteAGridHeader(name string, nCol, nRow int) (fout Fout) {
	cornerX := 0.0
	cornery := 0.0
	novalue := NONEVALUE
	cellsize := 1.0
	// create an ascii file, which contains the header
	MakeDir(name)
	file, err := os.OpenFile(name+".gz", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	gfile := gzip.NewWriter(file)
	fwriter := bufio.NewWriter(gfile)
	fout = Fout{file, gfile, fwriter}

	fout.Write(fmt.Sprintf("ncols %d\n", nCol))
	fout.Write(fmt.Sprintf("nrows %d\n", nRow))
	fout.Write(fmt.Sprintf("xllcorner     %f\n", cornerX))
	fout.Write(fmt.Sprintf("yllcorner     %f\n", cornery))
	fout.Write(fmt.Sprintf("cellsize      %f\n", cellsize))
	fout.Write(fmt.Sprintf("NODATA_value  %d\n", novalue))

	return fout
}

func WriteMetaFile(gridFilePath, title, labeltext, colormap, colorlistType string, colorlist []string, cbarLabel []string, ticklist []float64, factor float64, maxValue, minValue int, minColor string) {
	metaFilePath := gridFilePath + ".meta"
	MakeDir(metaFilePath)
	file, err := os.OpenFile(metaFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("title: '%s'\n", title))
	file.WriteString("yTitle: 0.88\n")
	file.WriteString("xTitle: 0.05\n")
	file.WriteString("removeEmptyColumns: True\n")
	file.WriteString(fmt.Sprintf("labeltext: '%s'\n", labeltext))
	if colormap != "" {
		file.WriteString(fmt.Sprintf("colormap: '%s'\n", colormap))
	}
	if colorlist != nil {
		file.WriteString("colorlist: \n")
		for _, item := range colorlist {
			file.WriteString(fmt.Sprintf(" - '%s'\n", item))
		}
	}
	if cbarLabel != nil {
		file.WriteString("cbarLabel: \n")
		for _, cbarItem := range cbarLabel {
			file.WriteString(fmt.Sprintf(" - '%s'\n", cbarItem))
		}
	}
	if ticklist != nil {
		file.WriteString("ticklist: \n")
		for _, tick := range ticklist {
			file.WriteString(fmt.Sprintf(" - %f\n", tick))
		}
	}
	if len(colorlistType) > 0 {
		file.WriteString(fmt.Sprintf("colorlisttype: %s\n", colorlistType))
	}
	file.WriteString(fmt.Sprintf("factor: %f\n", factor))
	if maxValue != NONEVALUE {
		file.WriteString(fmt.Sprintf("maxValue: %d\n", maxValue))
	}
	if minValue != NONEVALUE {
		file.WriteString(fmt.Sprintf("minValue: %d\n", minValue))
	}
	if len(minColor) > 0 {
		file.WriteString(fmt.Sprintf("minColor: %s\n", minColor))
	}
}

func WriteRows(fout Fout, extRow, extCol int, simGrid []int, gridSourceLookup [][]int) {
	for row := 0; row < extRow; row++ {

		for col := 0; col < extCol; col++ {
			refID := gridSourceLookup[row][col]
			if refID > 0 {
				fout.Write(strconv.Itoa(simGrid[refID-1]))
				fout.Write(" ")
			} else {
				fout.Write("-9999 ")
			}
		}
		fout.Write("\n")
	}
}
