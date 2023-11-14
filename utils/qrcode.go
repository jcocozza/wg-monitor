package utils
import (
	"github.com/skip2/go-qrcode"
	"log/slog"
)

func GenerateQRCode(data []byte) {
	err := qrcode.WriteFile(string(data), qrcode.Medium, 256, "qrcode.png")
	if err != nil {
		slog.Debug("Failed to generate QR code:", err)
		return
	}
	slog.Debug("QR code generated successfully!")
}