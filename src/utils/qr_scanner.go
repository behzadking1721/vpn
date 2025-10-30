package utils

// QRScannerInterface رابط برای اسکنر کد QR
type QRScannerInterface interface {
	// ScanQRCode اسکن کد QR و بازگشت محتوای آن
	ScanQRCode() (string, error)
	
	// IsAvailable بررسی در دسترس بودن اسکنر
	IsAvailable() bool
}

// MockQRScanner اسکنر QR ساختگی برای تست
type MockQRScanner struct{}

// NewMockQRScanner ایجاد یک اسکنر QR ساختگی
func NewMockQRScanner() *MockQRScanner {
	return &MockQRScanner{}
}

// ScanQRCode اسکن کد QR ساختگی
func (mqs *MockQRScanner) ScanQRCode() (string, error) {
	// در پیاده‌سازی واقعی، اینجا دوربین را فعال کرده و کد QR را اسکن می‌کنیم
	// برای تست، یک نمونه لینک VMess برمی‌گردانیم
	return "vmess://eyJ2IjogIjIiLCAicHMiOiAi5rWL6K+V5YWl5a2X56ym5LiyIiwgImFkZCI6ICJleGFtcGxlLmNvbSIsICJwb3J0IjogIjQ0MyIsICJpZCI6ICJ0ZXN0LXV1aWQiLCAiYWlkIjogIjAiLCAic2N5IjogImF1dG8iLCAibmV0IjogInRjcCIsICJ0eXBlIjogIm5vbmUiLCAiaG9zdCI6ICIiLCAicGF0aCI6ICIiLCAidGxzIjogInRscyJ9", nil
}

// IsAvailable بررسی در دسترس بودن اسکنر
func (mqs *MockQRScanner) IsAvailable() bool {
	// در پیاده‌سازی واقعی، این تابع بررسی می‌کند که آیا دوربین در دسترس است یا نه
	return true
}