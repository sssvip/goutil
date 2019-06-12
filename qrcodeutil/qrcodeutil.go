package qrcodeutil

import (
	"github.com/sssvip/qrterminal"
	"os"
)

var qrcodeConfig = qrterminal.Config{
	Level:          qrterminal.M,
	Writer:         os.Stdout,
	HalfBlocks:     true,
	BlackChar:      qrterminal.BLACK_BLACK,
	WhiteBlackChar: qrterminal.WHITE_BLACK,
	WhiteChar:      qrterminal.WHITE_WHITE,
	BlackWhiteChar: qrterminal.BLACK_WHITE,
	QuietZone:      1,
}

func QRCodeOutPut(text string) {
	qrterminal.GenerateWithConfig(text, qrcodeConfig)
}
