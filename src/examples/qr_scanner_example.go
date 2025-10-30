package main

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	
	"github.com/makiuchi-d/gozxing/qrcode/encoder"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func main() {
	fmt.Println("مثال اسکن کد QR")
	fmt.Println("==============")
	
	// ایجاد یک کد QR نمونه برای تست
	qrText := "vmess://eyJ2IjogIjIiLCAicHMiOiAi5rWL6K+V5YWl5a2X56ym5LiyIiwgImFkZCI6ICJleGFtcGxlLmNvbSIsICJwb3J0IjogIjQ0MyIsICJpZCI6ICJ0ZXN0LXV1aWQiLCAiYWlkIjogIjAiLCAic2N5IjogImF1dG8iLCAibmV0IjogInRjcCIsICJ0eXBlIjogIm5vbmUiLCAiaG9zdCI6ICIiLCAicGF0aCI6ICIiLCAidGxzIjogInRscyJ9"
	qrImage := createSampleQRCode(qrText)
	
	// ذخیره تصویر QR برای مرجع
	file, err := os.Create("sample_qr.png")
	if err != nil {
		fmt.Printf("خطا در ایجاد فایل تصویر: %v\n", err)
		return
	}
	defer file.Close()
	
	png.Encode(file, qrImage)
	fmt.Println("تصویر QR نمونه در فایل sample_qr.png ذخیره شد")
	
	// ایجاد اسکنر QR
	scanner := utils.NewRealQRScanner()
	
	// اسکن کد QR
	result, err := scanner.ScanQRCode(qrImage)
	if err != nil {
		fmt.Printf("خطا در اسکن کد QR: %v\n", err)
		return
	}
	
	fmt.Printf("محتوای کد QR: %s\n", result)
	
	// پردازش محتوای QR
	fmt.Println("در حال پردازش محتوای کد QR...")
	if len(result) > 50 {
		fmt.Printf("محتوای کد QR (قسمت اول): %.50s...\n", result)
	} else {
		fmt.Printf("محتوای کد QR: %s\n", result)
	}
}

// createSampleQRCode ایجاد یک تصویر QR نمونه
func createSampleQRCode(content string) image.Image {
	// ایجاد کد QR
	qrCode, err := encoder.Encode(content, qrcode.ErrorCorrectionLevel_L, nil)
	if err != nil {
		fmt.Printf("خطا در ایجاد کد QR: %v\n", err)
		return nil
	}
	
	// تبدیل به تصویر
	width := qrCode.GetWidth()
	height := qrCode.GetHeight()
	
	// ایجاد تصویر جدید
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// رنگ‌بندی پیکسل‌ها
	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}
	
	// پر کردن تصویر
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if qrCode.Get(x, y) {
				img.Set(x, y, black)
			} else {
				img.Set(x, y, white)
			}
		}
	}
	
	// اضافه کردن حاشیه سفید
	border := 4
	newWidth := width + 2*border
	newHeight := height + 2*border
	borderedImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.Draw(borderedImage, borderedImage.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)
	draw.Draw(borderedImage, image.Rect(border, border, border+width, border+height), img, image.Point{}, draw.Src)
	
	return borderedImage
}