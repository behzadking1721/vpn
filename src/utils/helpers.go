package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

// GenerateID تولید یک ID منحصر به فرد
func GenerateID() string {
	// تولید یک آرایه 16 بایتی تصادفی
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		// بازگشت به ID مبتنی بر زمان در صورت شکست تولید تصادفی
		return fmt.Sprintf("id_%d", time.Now().UnixNano())
	}

	// قالب‌بندی به عنوان رشته هگزادسیمال
	return fmt.Sprintf("%x", bytes)
}

// FormatBytes قالب‌بندی تعداد بایت به فرمت قابل خواندن
func FormatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// CalculatePingTime محاسبه زمان پینگ به میلی‌ثانیه
func CalculatePingTime(start, end time.Time) int {
	return int(end.Sub(start).Milliseconds())
}

// Base64Encode کدگذاری Base64
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode کدگشایی Base64
func Base64Decode(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}
