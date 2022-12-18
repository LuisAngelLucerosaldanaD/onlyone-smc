package ocr

type ResponseImgToBox struct {
	Code  int    `json:"code"`
	Data  []Box  `json:"data"`
	Error bool   `json:"error"`
	Msg   string `json:"msg"`
	Type  string `json:"type"`
}

type Box struct {
	BlockNum int    `json:"block_num"`
	Bottom   int    `json:"bottom"`
	Height   int    `json:"height"`
	Left     int    `json:"left"`
	Level    int    `json:"level"`
	LineNum  int    `json:"line_num"`
	PageNum  int    `json:"page_num"`
	Right    int    `json:"right"`
	Text     string `json:"text"`
	Top      int    `json:"top"`
	Width    int    `json:"width"`
}

type PersonInformation struct {
	Names          string `json:"names"`
	Surnames       string `json:"surnames"`
	IdentityNumber string `json:"identity_number"`
}
