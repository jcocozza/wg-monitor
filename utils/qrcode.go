package utils
import (
	"github.com/skip2/go-qrcode"
	"log/slog"
)

func GenerateQRCode(path string, data []byte) {
	slog.Debug("Generating QR code...")
	err := qrcode.WriteFile(string(data), qrcode.Medium, 256, path)
	if err != nil {
		slog.Error("Failed to generate QR code:", err)
		return
	}
}

func QRCodeData(data []byte, size int) []byte {
	slog.Debug("Generating QR code data...")
	qrCode, err := qrcode.New(string(data), qrcode.Medium)

	if err != nil {
		slog.Error("Failed to generate QR code:", err)
		return []byte{}
	}

	pngData, err := qrCode.PNG(size)

	if err != nil {
		slog.Error("Failed to generate QR code PNG:", err)
		return []byte{}
	}

	return pngData
}