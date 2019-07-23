package main

import (
	"fmt"

	"gocv.io/x/gocv"
	"image"
)

// CreateImg create a solid image based on params
func CreateImgByBGR(sizex int, sizey int, b float64, g float64, r float64) gocv.Mat {
	img := gocv.NewMatWithSizeFromScalar(gocv.NewScalar(b, g, r, 255), sizex, sizey, gocv.MatTypeCV8UC3)
	return img
}

// ResizeVideo resize a video, sizeX/Y refer to width/height
func ResizeVideo(srcPath string, dstPath string, sizeX, sizeY int) (err error) {
	capt, err := gocv.VideoCaptureFile(srcPath)
	defer capt.Close()
	if err != nil {
		return
	}

	fps := capt.Get(gocv.VideoCaptureFPS)
	writer, err := gocv.VideoWriterFile(dstPath, "avc1", fps, sizeX, sizeY, true)
	defer writer.Close()
	if err != nil {
		return
	}

	frame := gocv.NewMat()
	dstFrame := gocv.NewMat()
	defer frame.Close()
	defer dstFrame.Close()
	for {
		if ok := capt.Read(&frame); ok {
			gocv.Resize(frame, &dstFrame, image.Pt(sizeX, sizeY), 0, 0, gocv.InterpolationLinear)
			//win.IMShow(dstFrame)
			//gocv.WaitKey(100)
			err = writer.Write(dstFrame)
			if err != nil {
				return
			}
		} else {
			break
		}
	}

	return nil
}


// Convert replace green screen with a custom color
func Convert(srcPath string, dstPath string, r, g, b float64) (err error) {
	lb := gocv.NewScalar(68, 84, 153, 255)
	ub := gocv.NewScalar(80, 255, 255, 255)

	hsv := gocv.NewMat()
	defer hsv.Close()
	mask := gocv.NewMat()
	defer mask.Close()
	mask_inv := gocv.NewMat()
	defer mask.Close()
	frame := gocv.NewMat()
	defer frame.Close()
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
	defer kernel.Close()
	result := gocv.NewMat()
	defer result.Close()


	capt, err := gocv.VideoCaptureFile(srcPath)
	defer capt.Close()
	if err != nil {
		return
	}

	fps := capt.Get(gocv.VideoCaptureFPS)
	width := int(capt.Get(gocv.VideoCaptureFrameWidth))
	height := int(capt.Get(gocv.VideoCaptureFrameHeight))
	writer, err := gocv.VideoWriterFile(dstPath, "avc1", fps, width, height, true)
	defer writer.Close()

	if err != nil {
		return
	}

	ColorBg := CreateImgByBGR(height, width, b, g, r)
	for {
		person := gocv.NewMat()
		BgPerson := gocv.NewMat()
		if ok := capt.Read(&frame); ok {
			gocv.CvtColor(frame, &hsv, gocv.ColorBGRToHSV)
			gocv.InRangeWithScalar(hsv, lb, ub, &mask)

			gocv.BitwiseNot(mask, &mask_inv)

			gocv.Erode(mask_inv, &mask, kernel)
			gocv.BitwiseNot(mask, &mask_inv)
			gocv.BitwiseAndWithMask(frame, frame, &person, mask)

			gocv.BitwiseAndWithMask(ColorBg, ColorBg, &BgPerson, mask_inv)

			gocv.Add(BgPerson, person, &result)

			err = writer.Write(result)
			if err != nil {
				fmt.Printf("err occur when write frame: %s", err)
			}
		} else {
			break
		}
	}
	return nil
}



func main() {
	err := ResizeVideo("video.mp4", "resize.mp4", 960, 720)
	if err != nil {
		panic(err)
	}
}

