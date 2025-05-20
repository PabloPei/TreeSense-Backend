package utils

import "time"

var argentinaLoc, _ = time.LoadLocation("America/Argentina/Buenos_Aires")

// TODO: hacer que convierta a cualquier horario de donde se haya cargado
func ConvertUTCToArgentina(t time.Time) time.Time {
	return t.In(argentinaLoc)
}
