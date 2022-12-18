package ocr

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"onlyone_smc/internal/logger"
	"regexp"
	"strings"
)

var (
	DNIRgx    = "^[\\d\\-\\d]{8}?$"
	CedulaRgx = "^[\\d\\-\\d]{7,16}?$"
)

func ConsumeOcrToBox(image []byte, filename, url string) ([]Box, error) {

	imgReader := bytes.NewReader(image)
	body := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(body)

	writer, err := multipartWriter.CreateFormFile("img-encoding", filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(writer, imgReader); err != nil {
		return nil, err
	}

	multipartWriter.Close()

	request, err := http.NewRequest(http.MethodPost, url, body)
	request.Header.Add("Content-Type", multipartWriter.FormDataContentType())

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		logger.Error.Printf("no se  puedo enviar la petición: %v  -- log: ", err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error.Printf("no se pudo ejecutar defer body close: %v  -- log: ", err)
		}
	}(resp.Body)

	rsBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Printf("no se  puedo obtener respuesta: %v  -- log: ", err)
		return nil, err
	}

	resOcr := ResponseImgToBox{}

	err = json.Unmarshal(rsBody, &resOcr)
	if err != nil {
		logger.Error.Printf("no se pudo parsear la respuesta del servio ocr: %v  -- log: ", err)
		return nil, err
	}

	if resOcr.Error {
		logger.Error.Printf("error al consumir el servicio ocr: " + resOcr.Msg)
		return nil, fmt.Errorf(resOcr.Msg)
	}

	return resOcr.Data, nil
}

func GetDniInformation(url, filename string, imageBytes []byte) (*PersonInformation, error) {
	res := PersonInformation{}

	resBoxes, err := ConsumeOcrToBox(imageBytes, filename, url)
	if err != nil {
		logger.Error.Println("error al obtener la información del dni: " + err.Error())
		return nil, err
	}
	surname := findValueInBlockByNumAndLine("Primer Apellido", resBoxes)
	secondSurname := findValueInBlockByNumAndLine("Segundo Apellido", resBoxes)
	names := findValueInBlockByNumAndLine("Pre Nombres", resBoxes)
	dni := findValueInBlocksByRegex(resBoxes, "[0-9]{8}", DNIRgx, "", "")

	res.Surnames = strings.TrimSpace(surname + " " + secondSurname)
	res.Names = strings.TrimSpace(names)
	res.IdentityNumber = dni

	return &res, nil
}

func GetCedulaInformation(url, filename string, image []byte) (*PersonInformation, error) {
	res := PersonInformation{}

	resBoxes, err := ConsumeOcrToBox(image, filename, url)
	if err != nil {
		logger.Error.Println("error al obtener la información de la cedula: " + err.Error())
		return nil, err
	}

	dni := findValueInBlocksByRegex(resBoxes, "[0-9]{7,16}", CedulaRgx, ".", "")

	res.IdentityNumber = dni

	return &res, nil
}

func findBlockByInjector(injector string, data []Box) *Box {
	for _, item := range data {
		if item.Text == injector || strings.Contains(item.Text, injector) {
			return &item
		}
	}
	return nil
}

func findValueInBlockByNumAndLine(injector string, data []Box) string {
	boxFound := findBlockByInjector(injector, data)
	if boxFound == nil {
		return ""
	}
	for _, item := range data {
		if boxFound.BlockNum == item.BlockNum && item.LineNum > boxFound.LineNum {
			return item.Text
		}
	}
	return ""
}

func findValueInBlocksByRegex(boxes []Box, regexFind, regexValidator, replacer, replacerValue string) string {
	for _, box := range boxes {
		valueRxg := findValueByRegex(box.Text, regexFind, regexValidator, replacer, replacerValue)
		if valueRxg != "" {
			return valueRxg
		}
	}
	return ""
}

func findValueByRegex(value, regexFind, regexValidator, replacer, replacerValue string) string {
	newValue := value
	if replacer != "" {
		newValue = strings.ReplaceAll(value, replacer, replacerValue)
	}
	rgx := regexp.MustCompile(regexFind)
	matches := rgx.FindAllString(newValue, -1)
	if matches == nil {
		return ""
	}
	for _, match := range matches {
		regex := regexp.MustCompile(regexValidator)
		if regex.MatchString(match) {
			return match
		}
	}
	return ""
}
