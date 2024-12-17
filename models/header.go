package models

type ReqHeader struct {
	DeviceType    string `header:"P-DeviceType"`
	Authorization string `header:"Authorization"`
}
