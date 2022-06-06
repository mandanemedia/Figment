package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const referenceKeyword = "reference"
const deviceType1 = "thermometer"
const deviceType2 = "humidity"
const filePath = "sample.log"

type Thermometer struct {
	Name            string
	Quality         string
	Sum             float64
	N               int64
	SumSquareRoot   float64
	Mean            float64
	SD              float64
	CheckMeanWithin bool
}

type Humidity struct {
	Name      string
	Quality   string
	Discarded bool
}

func detectReferences(line string, referenceKeyword string, refernces *[]float64) {
	if !strings.Contains(line, referenceKeyword) {
		log.Fatal("Log file is not in proper format on the first line")
	}
	// remove reference keyword from first line
	var firstLine = strings.Replace(line, referenceKeyword, "", 1)

	// split them base on space and grab index 1 and 2
	var array = strings.Split(firstLine, " ")

	if len(array) < 3 {
		log.Fatalf("Please set reference values in proper format")
	}

	value1, err1 := strconv.ParseFloat(array[1], 32)
	if err1 == nil {
		*refernces = append(*refernces, value1)
	} else {
		log.Fatalf("Reference value1 is not set properly")
	}
	value2, err2 := strconv.ParseFloat(array[2], 32)
	if err2 == nil {
		*refernces = append(*refernces, value2)
	} else {
		log.Fatalf("Reference value2 is not set properly")
	}
	// fmt.Printf("%+v\n", *refernces)
}

func addHumidity(line string, humidity *Humidity) string {
	// split them base on whitespace and grab index 1 and 2
	var array = strings.Split(line, " ")
	if len(array) < 2 {
		log.Fatalf("Please set Device Name in proper format for Humidity")
	}
	humidity.Name = array[1]
	humidity.Quality = "OK"
	humidity.Discarded = false
	return array[1]
}

func addThermometer(line string, thermometer *Thermometer) string {
	// split them base on whitespace and grab index 1 and 2
	var array = strings.Split(line, " ")
	if len(array) < 2 {
		log.Fatalf("Please set Device Name in proper format for Thermometer")
	}
	resetThermometerData(thermometer)
	thermometer.Name = array[1]
	return array[1]
}

func resetThermometerData(thermometer *Thermometer) {
	thermometer.Name = ""
	thermometer.Quality = ""
	thermometer.Sum = 0
	thermometer.N = 0
	thermometer.SumSquareRoot = 0
	thermometer.Mean = 0
	thermometer.SD = 0
	thermometer.CheckMeanWithin = true
}

func resetHumidityData(humidity *Humidity) {
	humidity.Name = ""
	humidity.Quality = ""
	humidity.Discarded = false
}

func checkThermometerData(line string, name string, refernce float64, thermometer *Thermometer) {
	if thermometer.Name != name {
		log.Fatalf("Device name is not matched, correct the format or the order")
	}
	var data = strings.Split(line, " ")
	if len(data) < 2 {
		log.Fatalf("Please set data in proper format for %s Thermometer", name)
	}

	value, err := strconv.ParseFloat(data[2], 32)
	if err != nil {
		log.Fatalf("data value for %s is not a number", name)
	}

	thermometer.Sum += value
	thermometer.N += 1
	thermometer.Mean = thermometer.Sum / (float64(thermometer.N))
	// sum((current - avg)^2)
	thermometer.SumSquareRoot += math.Pow(value-thermometer.Mean, 2)
	thermometer.SD = math.Sqrt(thermometer.SumSquareRoot / (float64(thermometer.N) - 1))

	thermometer.Quality = "precise"
	// if the mean of the readings is within 0.5 degrees of the refernce
	if thermometer.Mean-0.5 <= refernce && refernce <= thermometer.Mean+0.5 {
		// thermometer.CheckMeanWithin = true
		if thermometer.SD < 3 {
			thermometer.Quality = "ultra precise"
		} else if thermometer.SD < 5 {
			thermometer.Quality = "very precise"
		}
	}
}

func checkHumidityData(line string, name string, refernce float64, humidity *Humidity) {
	if humidity.Name != name {
		log.Fatalf("Device name is not matched, correct the format or the order")
	}
	var data = strings.Split(line, " ")
	if len(data) < 2 {
		log.Fatalf("Please set data in proper format for %s Humidity", name)
	}

	value, err := strconv.ParseFloat(data[2], 32)
	if err != nil {
		log.Fatalf("data value for %s is not a number", name)
	}

	if !humidity.Discarded {
		var upper = refernce * 1.01
		var lower = refernce * 0.99
		if !(lower <= value && value <= upper) {
			humidity.Discarded = true
			humidity.Quality = "discard"
		}
	}
}

func printLastDevice(lastDevice string, isThermometer bool,
	thermometer *Thermometer, humidity *Humidity) {

	// make sure there is a lastDevice exist
	if len(lastDevice) > 0 {
		if isThermometer {
			if thermometer.Name != lastDevice {
				log.Fatalf("Device name is not matched, correct the format or the order")
			}
			fmt.Printf("%s: %s\n", thermometer.Name, thermometer.Quality)
			resetThermometerData(thermometer)
		} else {
			if humidity.Name != lastDevice {
				log.Fatalf("Device name is not matched, correct the format or the order")
			}
			fmt.Printf("%s: %s\n", humidity.Name, humidity.Quality)
			resetHumidityData(humidity)
		}
	} else {
		fmt.Printf("Output\n")
	}
}

func main() {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open")
	}

	var refernces []float64
	thermometer := Thermometer{
		Name:          "",
		Quality:       "",
		Sum:           0,
		N:             0,
		SumSquareRoot: 0,
		Mean:          0,
		SD:            0,
	}
	humidity := Humidity{
		Name:      "",
		Quality:   "",
		Discarded: false,
	}
	var lastDevice = ""
	var isThermometer = true

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var line = scanner.Text()
		// detect the first line
		if len(refernces) == 0 {
			detectReferences(line, referenceKeyword, &refernces)
		} else if strings.Contains(line, deviceType1) { // detect devices' name for thermometer
			printLastDevice(lastDevice, isThermometer, &thermometer, &humidity)
			isThermometer = true
			lastDevice = addThermometer(line, &thermometer)
		} else if strings.Contains(line, deviceType2) { // detect devices' name for humidity
			printLastDevice(lastDevice, isThermometer, &thermometer, &humidity)
			isThermometer = false
			lastDevice = addHumidity(line, &humidity)
		} else if strings.Contains(line, lastDevice) { // Add datapoint
			if isThermometer {
				checkThermometerData(line, lastDevice, refernces[0], &thermometer)
			} else {
				checkHumidityData(line, lastDevice, refernces[1], &humidity)
			}
		}
	}
	// print the last device
	printLastDevice(lastDevice, isThermometer, &thermometer, &humidity)
	file.Close()
}
