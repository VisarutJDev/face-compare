package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"

	"gocv.io/x/gocv"
)

// FaceDetectFromWebCam is a function to detect human face from web cam
func FaceDetectFromWebCam() {
	if len(os.Args) < 3 {
		fmt.Println("How to run:\n\tFaceDetectFromWebCam [camera ID] [classifier XML file]")
		return
	}

	// parse args
	deviceID, _ := strconv.Atoi(os.Args[1])
	xmlFile := os.Args[2]

	// open webcam
	webcam, err := gocv.VideoCaptureDevice(int(deviceID))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	// open display window
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}

	fmt.Printf("start reading camera device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))

		// draw a rectangle around each face on the original image,
		// along with text identifying as "Human"
		for _, r := range rects {
			gocv.Rectangle(&img, r, blue, 3)

			size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
		}

		// show the image in the window, and wait 1 millisecond
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}

// FaceDetectFromImage is a funtion to detect human face from 2 picture
func FaceDetectFromImage() {
	if len(os.Args) < 4 {
		fmt.Println("Please sent these argument to us.")
		fmt.Println("How to run:\n\t FaceDetectFromImage [path/to/file/image1.jpg] [path/to/file/image2.jpg] [classifier XML file]")
		return
	}
	// parse args
	imageFile1 := os.Args[1]
	imageFile2 := os.Args[2]
	xmlFile := os.Args[3]

	// prepare image matrix
	img1 := gocv.IMRead(imageFile1, gocv.IMReadAnyColor)
	img2 := gocv.IMRead(imageFile2, gocv.IMReadAnyColor)
	if img1.Empty() {
		log.Panic("Can not read Image file : ", imageFile1)
		return
	}
	if img2.Empty() {
		log.Panic("Can not read Image file : ", imageFile2)
		return
	}
	defer img1.Close()
	defer img2.Close()

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}

	// color for the rect when faces detected
	// blue := color.RGBA{0, 0, 255, 0}

	// detect faces
	rects1 := classifier.DetectMultiScale(img1)
	fmt.Printf("found %d faces\n", len(rects1))
	rects2 := classifier.DetectMultiScale(img2)
	fmt.Printf("found %d faces\n", len(rects2))

	// open display window
	window1 := gocv.NewWindow("Face Detect 1")
	defer window1.Close()

	// open display window
	window2 := gocv.NewWindow("Face Detect 2")
	defer window2.Close()

	// draw a rectangle around each face on the original image,
	// along with text identifying as "Human"
	var faceCrop1 gocv.Mat
	var faceCrop2 gocv.Mat
	images := make([]gocv.Mat, 2)
	for _, r := range rects1 {
		// gocv.Rectangle(&img1, r, blue, 3)
		faceCrop1 = img1.Region(r)
		// size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
		// pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
		// gocv.PutText(&img1, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
	}

	for _, r2 := range rects2 {
		// gocv.Rectangle(&img2, r2, blue, 3)
		faceCrop2 = img2.Region(r2)
		// size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
		// pt := image.Pt(r2.Min.X+(r2.Min.X/2)-(size.X/2), r2.Min.Y-2)
		// gocv.PutText(&img2, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
	}

	images[0] = faceCrop1
	images[1] = faceCrop2
	FaceCompare(images)
	// for {
	// 	// show the image in the window, and wait 1 millisecond
	// 	window1.IMShow(faceCrop1)
	// 	window2.IMShow(faceCrop2)
	// 	if window1.WaitKey(1) >= 0 || window2.WaitKey(1) >= 0 {
	// 		break
	// 	}

	// }

}
