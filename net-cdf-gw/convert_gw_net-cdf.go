package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/batchatco/go-native-netcdf/netcdf"
	"github.com/batchatco/go-native-netcdf/netcdf/api"
	"github.com/mazznoer/colorgrad"
)

func main() {

	//inputFile := flag.String("input", "OCEANIA_WTD_monthlymeans.nc", "Input file")
	inputFile := flag.String("input", "EURASIA_WTD_annualmean.nc", "Input file")
	startTime := flag.String("start", "2000-01-01", "Start time")
	endTime := flag.String("end", "2010-12-31", "End time")
	eumask := flag.String("eumask", "stu_eu_layers.csv", "EU mask")

	flag.Parse()

	// Open the file
	nc, err := netcdf.Open(*inputFile)
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	// part 1: get a base overview of the file
	for _, key := range nc.Attributes().Keys() {
		if val, ok := nc.Attributes().Get(key); ok {
			fmt.Println(key, ":")
			fmt.Println(val)
		}
	}

	fmt.Println(nc.ListVariables())
	fmt.Println(nc.ListSubgroups())
	fmt.Println(nc.ListTypes())
	fmt.Println(nc.ListDimensions())
	for _, dim := range nc.ListDimensions() {
		val, has := nc.GetDimension(dim)
		if has {
			fmt.Println(dim, val)
		}
	}

	// // Read the NetCDF variable from the file
	// vr, _ := nc.GetVariable("lat")
	// if vr == nil {
	// 	panic("lat variable not found")
	// }

	// // Cast the data into a Go type we can use
	// lats, has := vr.Values.([]float32)
	// if !has {
	// 	panic("latitude data not found")
	// }
	// for i, lat := range lats {
	// 	fmt.Println(i, lat)
	// }

	for _, attr := range nc.ListVariables() {

		fmt.Println(attr)
		// vr, _ := nc.GetVariable(attr)
		// if vr == nil {
		// 	fmt.Printf("%s variable not found\n", attr)
		// 	continue
		// }
		// fmt.Println("dimensions:", vr.Dimensions)
		// fmt.Println("Attributes:")
		// for _, key := range vr.Attributes.Keys() {
		// 	if val, ok := vr.Attributes.Get(key); ok {
		// 		fmt.Println("  ", key, ":", val)

		// 	}
		// }
		getVar, err := nc.GetVarGetter(attr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Len:", getVar.Len())
		fmt.Println("Type:", getVar.Type())
		gotype := getVar.GoType()
		fmt.Println("GoType:", gotype)
		lenDim := len(getVar.Dimensions())
		fmt.Println("Dimensions:", getVar.Dimensions(), lenDim)

		fmt.Println("Attributes:")
		for _, key := range getVar.Attributes().Keys() {
			if val, ok := getVar.Attributes().Get(key); ok {
				fmt.Println("  ", key, ":", val)
			}
		}
		// vals, err := getVar.Values()
		// if err != nil {
		// 	fmt.Println(err)
		// 	continue
		// }
		// switch gotype {
		// case "float32":

		// 	if lenDim == 0 {
		// 		fmt.Println(vals.(float32))
		// 	} else if lenDim == 1 {
		// 		for i, val := range vals.([]float32) {
		// 			fmt.Println(i, val)
		// 		}
		// 	} else if lenDim == 2 {
		// 		for i, val := range vals.([][]float32) {
		// 			fmt.Println(i, val)
		// 		}
		// 	} else if lenDim == 3 {
		// 		for i, val := range vals.([][][]float32) {
		// 			fmt.Println(i, val)
		// 		}
		// 	}
		// case "int8":
		// 	if lenDim == 0 {
		// 		fmt.Println(vals.(int8))
		// 	} else if lenDim == 1 {
		// 		for i, val := range vals.([]int8) {
		// 			fmt.Println(i, val)
		// 		}
		// 	} else if lenDim == 2 {
		// 		for i, val := range vals.([][]int8) {
		// 			fmt.Println(i, val)
		// 		}
		// 	} else if lenDim == 3 {
		// 		for i, val := range vals.([][][]int8) {
		// 			fmt.Println(i, val)
		// 		}
		// 	}
		// case "int16":
		// 	if lenDim == 0 {
		// 		fmt.Println(vals.(int16))
		// 	} else if lenDim == 1 {
		// 		for i, val := range vals.([]int16) {
		// 			fmt.Println(i, val)
		// 		}
		// 	} else if lenDim == 2 {
		// 		for i, val := range vals.([][]int16) {
		// 			fmt.Println(i, val)
		// 		}
		// 	} else if lenDim == 3 {
		// 		for i, val := range vals.([][][]int16) {
		// 			fmt.Println(i, val)
		// 		}
		// 	}
		// }

		// if len(vr.Dimensions) == 2 {
		// 	fmt.Println("2D")
		// 	// Cast the data into a Go type we can use
		// 	data, has := vr.Values.([][]float32)
		// 	if !has {
		// 		fmt.Printf("%s data not found\n", attr)
		// 		continue
		// 	}
		// 	for i, val := range data {
		// 		fmt.Println(i, val)
		// 	}
		// } else if len(vr.Dimensions) == 1 {
		// 	vals, has := vr.Values.([]float32)
		// 	if !has {
		// 		fmt.Printf("%s data not found\n", attr)
		// 		continue
		// 	}
		// 	for i, lat := range vals {
		// 		fmt.Println(i, lat)
		// 	}
		// }
	}

	//Part 2: transform the data into..
	start, _ := time.Parse("2006-01-02", *startTime)
	if err != nil {
		log.Fatal(err)
	}
	end, _ := time.Parse("2006-01-02", *endTime)
	if err != nil {
		log.Fatal(err)
	}
	createGWTimeSeries(&nc, start, end, *inputFile+".png", *eumask)
}

// a Hermes ground water time series
// create mapping csv file (lat, lon, gwId)
// ground water time series file (gwId, date, value)
// requires input time range  (start, end)
func createGWTimeSeries(nc *api.Group, start, end time.Time, imgFileName, maskFile string) {

	// time
	timeVar, err := (*nc).GetVarGetter("time")
	if err != nil {
		log.Fatal(err)

	}
	valsTime, err := timeVar.Values()
	if err != nil {
		log.Fatal(err)
	}
	var timeValues []float32
	switch timeVar.GoType() {
	case "float32":
		timeValues = valsTime.([]float32)
	case "int8":
		timeValuesInt := valsTime.([]int8)
		timeValues = make([]float32, 0, len(timeValuesInt))
		for _, val := range timeValuesInt {
			timeValues = append(timeValues, float32(val))
		}
	}

	// latitudes
	latVar, err := (*nc).GetVarGetter("lat")
	if err != nil {
		log.Fatal(err)
	}
	lenLat := latVar.Len()
	valsLat, err := latVar.Values()
	if err != nil {
		log.Fatal(err)
	}
	valLat := valsLat.([]float32)
	// longitude
	lonVar, err := (*nc).GetVarGetter("lon")
	if err != nil {
		log.Fatal(err)
	}
	lenLon := lonVar.Len()
	valsLon, err := lonVar.Values()
	if err != nil {
		log.Fatal(err)
	}
	valLon := valsLon.([]float32)
	// ground water
	WTDVar, err := (*nc).GetVarGetter("WTD")
	if err != nil {
		log.Fatal(err)
	}
	valsWTD, err := WTDVar.Values()
	if err != nil {
		log.Fatal(err)
	}
	valWTD := valsWTD.([][][]int16)
	scaleFactor := 1.0
	if val, ok := WTDVar.Attributes().Get("scale_factor"); ok {
		scaleFactor = val.(float64)
	}
	var add_offset float64 = 0.0
	if val, ok := WTDVar.Attributes().Get("add_offset"); ok {
		add_offset = val.(float64)
	}

	// mask for valid data
	maskVar, err := (*nc).GetVarGetter("mask")
	if err != nil {
		log.Fatal(err)
	}
	valsMask, err := maskVar.Values()
	if err != nil {
		log.Fatal(err)
	}
	// // out map for mapping csv file (lat, lon, gwId)
	// latLonMappings := make([]gwMapping, 0, lenLat*lenLon)

	// // out map for time series csv file (gwId, date, value)
	// uniqueGWId := make(map[string][]float64)

	gwValues := make([][]float64, lenLat)
	min, max := 0.0, 0.0
	init := false
	// counter := 0
	// loop through lat, lon, mask, time and get WTD
	for iLat := int64(0); iLat < lenLat; iLat++ {
		gwValues[iLat] = make([]float64, lenLon)
		for iLon := int64(0); iLon < lenLon; iLon++ {
			// check against mask (1 = valid, 0 = invalid)
			if valsMask.([][]int8)[iLat][iLon] == 1 {
				// create gwId with 12 digits for wtd values
				//gwIdValues := make([]string, len(timeValues))
				timeSlice := make([]float64, len(timeValues))
				// valid := true
				for iTime := 0; iTime < len(timeValues); iTime++ {
					value := valWTD[iTime][iLat][iLon]
					timeSlice[iTime] = math.Ceil((float64(value)*scaleFactor + add_offset))
					//gwIdValues[iTime] = fmt.Sprintf("%01.2f", timeSlice[iTime])

					if !init {
						min = timeSlice[iTime]
						max = timeSlice[iTime]
						init = true
					}
					if timeSlice[iTime] < min {
						min = timeSlice[iTime]
					}
					if timeSlice[iTime] > max {
						max = timeSlice[iTime]
					}
					gwValues[iLat][iLon] = timeSlice[iTime]

					//img.Set(int(iLon), int(lenLat-iLat), ToColor(timeSlice[iTime]))
					// if timeSlice[iTime] > 100 {
					// 	// 	fmt.Println("value", value)
					// 	valid = false
					// 	break
					// }
				}
				// if !valid {
				// 	continue
				// }

				// gwId := strings.Join(gwIdValues, "")
				// if _, ok := uniqueGWId[gwId]; !ok {
				// 	uniqueGWId[gwId] = timeSlice
				// 	counter++
				// }
				// latLonMappings = append(latLonMappings, gwMapping{gwId: len(uniqueGWId), lat: valsLat.([]float32)[iLat], lon: valsLon.([]float32)[iLon]})
			} else {
				gwValues[iLat][iLon] = 2 // invalid
				//img.Set(int(iLon), int(lenLat-iLat), color.RGBA{188, 190, 198, 0xff})
			}
		}
	}

	// find start index (0) of lon and lat (may be out of range)
	startLat, startLon := findStartLatLon(valLat, valLon, lenLat, lenLon)
	stepLatSize := math.Abs(float64(valLat[0] - valLat[1]))
	stepLonSize := math.Abs(float64(valLon[0] - valLon[1]))

	// open file maskFile
	sourcefile, err := os.Open(maskFile)
	if err == nil {
		defer sourcefile.Close()
		// write mapping csv file
		gwOutPath := strings.SplitN(filepath.Base(maskFile), ".", 2)[0] + "_groundwater.csv"
		gwOut, err := os.Create(gwOutPath)
		if err != nil {
			log.Fatal(err)
		}
		defer gwOut.Close()
		w := bufio.NewWriter(gwOut)
		w.WriteString("Column,Row,latitude,longitude,groundwater\n")

		lookup := make(map[GridCoord]float64)
		firstLine := true
		colID := -1
		rowID := -1
		latID := -1
		lonID := -1

		colExt := 0
		rowExt := 0
		scanner := bufio.NewScanner(sourcefile)
		for scanner.Scan() {
			line := scanner.Text()
			tokens := strings.Split(line, ",")
			if firstLine {
				firstLine = false
				// Column,Row,latitude,longitude,groundwater
				for index, token := range tokens {
					if token == "col" {
						colID = index
					}
					if token == "row" {
						rowID = index
					}
					if token == "latitude" {
						latID = index
					}
					if token == "longitude" {
						lonID = index
					}
				}
			} else {
				col, _ := strconv.ParseInt(tokens[colID], 10, 64)
				row, _ := strconv.ParseInt(tokens[rowID], 10, 64)
				lat, _ := strconv.ParseFloat(tokens[latID], 64)
				lon, _ := strconv.ParseFloat(tokens[lonID], 64)

				// calculate index of lat and lon in gwValues
				iLat := ((startLat) + int64(lat/stepLatSize))
				iLon := ((startLon) + int64(lon/stepLonSize))

				// check if lat and lon of a neibor is closer, to correct rounding errors
				iLat = isNeiborCloser(valLat, iLat, lat)
				iLon = isNeiborCloser(valLon, iLon, lon)

				gw := gwValues[iLat][iLon] * -1
				if gw < 0 {
					gw = 0
				}
				lookup[GridCoord{int(row), int(col)}] = gw

				// find image extension
				if int(col) > colExt {
					colExt = int(col)
				}
				if int(row) > rowExt {
					rowExt = int(row)
				}

				// write mapping to csv file
				w.WriteString(fmt.Sprintf("%s,%s,%f,%f,%1.0f\n", tokens[colID], tokens[rowID], lat, lon, gw))

				// print warnings if lat or lon is not the most accurate
				prevLat := math.Abs(float64(valLat[iLat])-lat) > math.Abs(float64(valLat[iLat-1])-lat)
				nextLat := math.Abs(float64(valLat[iLat])-lat) > math.Abs(float64(valLat[iLat+1])-lat)

				if prevLat || nextLat {
					whatOffset := 0
					if prevLat {
						whatOffset = -1
					} else {
						whatOffset = 1
					}
					fmt.Printf("Row_Col %d_%d - Lat Is: %v Wants: %v Prev: %v Next: %v Offset %d \n", row, col, valLat[iLat], lat, valLat[iLat-1], valLat[iLat+1], whatOffset)
				}
				preLon := math.Abs(float64(valLon[iLon])-lon) > math.Abs(float64(valLon[iLon-1])-lon)
				nextLon := math.Abs(float64(valLon[iLon])-lon) > math.Abs(float64(valLon[iLon+1])-lon)

				if preLon || nextLon {
					whatOffset := 0
					if preLon {
						whatOffset = -1
					} else {
						whatOffset = 1
					}
					fmt.Printf("Row_Col %d_%d - Lon Is: %v Wants: %v Prev: %v Next: %v Offset %d\n", row, col, valLon[iLon], lon, valLon[iLon-1], valLon[iLon+1], whatOffset)
				}
			}
		}
		if err := w.Flush(); err != nil {
			log.Fatal(err)
		}

		img := generatePic(colExt, rowExt)
		gradup := colorgrad.Viridis()
		for row := 0; row < rowExt; row++ {
			for col := 0; col < colExt; col++ {
				if val, ok := lookup[GridCoord{int(row), int(col)}]; ok {
					img.Set(int(col), int(row), ToColor(val, -9999, 0, 10, &gradup))
				} else {
					img.Set(int(col), int(row), ToColor(-9999, -9999, 0, 10, &gradup))
				}
			}
		}
		saveImg(img, imgFileName)
	} else {
		// render imge without mask
		img := generatePic(int(lenLon), int(lenLat))
		grad := colorgrad.Viridis()
		for iLat := int64(0); iLat < lenLat; iLat++ {
			for iLon := int64(0); iLon < lenLon; iLon++ {
				img.Set(int(iLon), int(lenLat-iLat), ToColor(gwValues[iLat][iLon]*(-1), -2, 0, 10, &grad))
			}
		}
		saveImg(img, imgFileName)
	}
}

type GridCoord struct {
	row int
	col int
}

func isNeiborCloser(arr []float32, idx int64, val float64) int64 {

	if idx+1 < int64(len(arr)) {
		next := math.Abs(float64(arr[idx])-val) > math.Abs(float64(arr[idx+1])-val)
		if next {
			return idx + 1
		}
	}
	if idx > 0 {
		prev := math.Abs(float64(arr[idx])-val) > math.Abs(float64(arr[idx-1])-val)
		if prev {
			return idx - 1
		}
	}
	return idx
}

func generatePic(width, height int) *image.RGBA {

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	return img
}

func saveImg(img *image.RGBA, imgName string) {
	// Encode as PNG.
	f, _ := os.Create(imgName)
	png.Encode(f, img)
}
func ToColor(val, nodata, minValue, maxValue float64, grad *colorgrad.Gradient) color.RGBA {
	if val == nodata {
		return color.RGBA{188, 190, 198, 0xff} // blank
	}
	valRange := maxValue - minValue
	if val > maxValue {
		val = maxValue
	}
	if val < minValue {
		val = minValue
	}
	perc := math.Abs((val - minValue) / valRange)
	r, g, b := grad.At(perc).RGB255()
	texture := color.RGBA{r, g, b, 0xff} // blank

	return texture

}

func findStartLatLon(valLat, valLon []float32, lenLat, lenLon int64) (int64, int64) {

	nearestToMiddle := func(val []float32, len int64, grads float64) int64 {
		min := math.Abs(float64(val[0]))
		currentIdx := int64(0)
		for i := int64(1); i < len; i++ {
			currVal := math.Abs(float64(val[i]))
			if currVal < min {
				min = currVal
				currentIdx = i
			}
		}
		// calculate step size

		if math.Abs(float64(valLat[currentIdx])) > 0.0001 {
			otherIdx := currentIdx + 1
			if otherIdx >= len {
				otherIdx = currentIdx - 1
			}
			stepSize := math.Abs(float64(valLat[currentIdx] - valLat[otherIdx]))
			numSteps := grads / stepSize

			currentIdx = currentIdx - int64(math.Round(float64(valLat[currentIdx])/numSteps))
		}

		return currentIdx
	}
	idxLat := nearestToMiddle(valLat, lenLat, 180.0)
	idxLon := nearestToMiddle(valLon, lenLon, 365.0)

	return idxLat, idxLon
}
