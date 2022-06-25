package aws_ia

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/textract"
	"onlyone_smc/internal/logger"
	"regexp"
	"strings"
)

const (
	AWS_BUCKET = "bjungle-files"
	AWS_REGION = "us-east-1"
)

var (
	blockMap  = make(map[string]*textract.Block)
	keyMap    = make(map[string]*textract.Block)
	valueMap  = make(map[string]*textract.Block)
	DNIRgx    = regexp.MustCompile("^[\\d\\-\\d]{7,10}?$")
	CedulaRgx = regexp.MustCompile("^((\\d{8})|(\\d{10})|(\\d{11})|(\\d{6}-\\d{5}))?$")
	user      User
)

func CompareFaces(face1, face2 []byte) (bool, error) {

	sess, err := sessionAws()
	if err != nil {
		return false, err
	}

	svc := rekognition.New(sess)
	input := &rekognition.CompareFacesInput{
		SimilarityThreshold: aws.Float64(90.000000),
		SourceImage: &rekognition.Image{
			Bytes: face1,
		},
		TargetImage: &rekognition.Image{
			Bytes: face2,
		},
	}

	result, err := svc.CompareFaces(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rekognition.ErrCodeInvalidParameterException:
				logger.Error.Println("InvalidParameterException:", aerr.Error())
			case rekognition.ErrCodeInvalidS3ObjectException:
				logger.Error.Println("ErrCodeInvalidS3ObjectException:", aerr.Error())
			case rekognition.ErrCodeImageTooLargeException:
				logger.Error.Println("ErrCodeImageTooLargeException:", aerr.Error())
			case rekognition.ErrCodeAccessDeniedException:
				logger.Error.Println("ErrCodeAccessDeniedException:", aerr.Error())
			case rekognition.ErrCodeInternalServerError:
				logger.Error.Println("ErrCodeInternalServerError:", aerr.Error())
			case rekognition.ErrCodeThrottlingException:
				logger.Error.Println("ErrCodeThrottlingException:", aerr.Error())
			case rekognition.ErrCodeProvisionedThroughputExceededException:
				logger.Error.Println("ErrCodeProvisionedThroughputExceededException:", aerr.Error())
			case rekognition.ErrCodeInvalidImageFormatException:
				logger.Error.Println("ErrCodeInvalidImageFormatException:", aerr.Error())
			default:
				logger.Error.Println("Other error:", aerr.Error())
			}
		} else {
			logger.Error.Println("Other error:", aerr.Error())
		}
		return false, err
	}
	var similarity float64

	if result.FaceMatches == nil || len(result.FaceMatches) == 0 {
		return false, fmt.Errorf("faces are not similar enough: %f", similarity)
	}

	similarity = *result.FaceMatches[0].Similarity

	if similarity <= 90 {
		return false, fmt.Errorf("faces are not similar enough: %f", similarity)
	}

	return true, nil
}

func GetUserFields(img []byte) (*User, error) {

	sess, err := sessionAws()
	if err != nil {
		return nil, err
	}

	svc := textract.New(sess)
	input := &textract.AnalyzeDocumentInput{
		Document: &textract.Document{
			Bytes: img,
		},
		FeatureTypes: []*string{
			aws.String("FORMS"),
		},
	}
	result, err := svc.AnalyzeDocument(input)
	if err != nil {
		return nil, err
	}

	getKeyValueMap(result.Blocks)
	for _, key := range keyMap {
		valueBlock := findValueBlock(key)
		keyValue := getText(key)
		value := getText(valueBlock)
		if keyValue != "" {
			setUserValues(keyValue, value)
		}
	}

	if user.IdentityNumber == "" {
		for _, block := range result.Blocks {
			if *block.BlockType == "LINE" {
				if DNIRgx.MatchString(*block.Text) {
					dni := *block.Text
					if len(dni) == 10 {
						user.IdentityNumber = dni[:len(dni)-2]
						break
					}

					if len(dni) == 9 {
						user.IdentityNumber = dni[:len(dni)-1]
						break
					}

					user.IdentityNumber = dni
					break
				}
				if CedulaRgx.MatchString(*block.Text) {
					user.IdentityNumber = strings.ReplaceAll(*block.Text, ".", "")
					break
				}
			}
		}
	}

	return &user, nil
}

func setUserValues(key, value string) {
	key = strings.ReplaceAll(key, ":", "")

	if key == "Pre Nombres" || key == "NOMBRES" || key == "Nombres" {
		user.Names = value
		return
	}

	if key == "APELLIDOS" || key == "Apellidos" {
		surname := strings.Split(value, " ")
		user.Surname = surname[0]
		if len(surname) >= 2 {
			user.SecondSurname = surname[1]
		}
		return
	}

	if key == "Primer Apellido" {
		user.Surname = value
		return
	}

	if key == "Segundo Apellido" {
		user.SecondSurname = value
		return
	}

	if key == "DNI" {
		if len(value) == 10 {
			user.IdentityNumber = value[:len(value)-2]
		} else if len(value) == 9 {
			user.IdentityNumber = value[:len(value)-1]
		} else {
			user.IdentityNumber = value
		}
		return
	}

	if key == "NUMERO" || key == "NÚMERO" || key == "Número" || key == "número" {
		user.IdentityNumber = strings.ReplaceAll(value, ".", "")
		return
	}
}

func getKeyValueMap(blocks []*textract.Block) {
	for _, block := range blocks {
		blockId := *block.Id
		blockMap[blockId] = block
		if *block.BlockType == "KEY_VALUE_SET" {
			for _, entityType := range block.EntityTypes {
				if *entityType == "KEY" {
					keyMap[blockId] = block
				} else {
					valueMap[blockId] = block
				}
			}
		}
	}
}

func findValueBlock(keyBlock *textract.Block) *textract.Block {
	for _, relationship := range keyBlock.Relationships {
		if *relationship.Type == "VALUE" {
			for _, id := range relationship.Ids {
				for _, block := range valueMap {
					if *id == *block.Id {
						return valueMap[*id]
					}
				}
			}
		}
	}

	return nil
}

func getText(keyBlock *textract.Block) string {
	var text string
	if keyBlock != nil && keyBlock.Relationships != nil {
		for _, line := range keyBlock.Relationships {
			if *line.Type == "CHILD" {
				for _, id := range line.Ids {
					word := blockMap[*id]
					if *word.BlockType == "WORD" {
						text += *word.Text + " "
					}
					if *word.BlockType == "SELECTION_ELEMENT" && *word.SelectionStatus == "SELECTED" {
						text += "X"
					}
				}
			}
		}
	}

	return strings.Trim(text, " ")
}

func sessionAws() (client.ConfigProvider, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWS_REGION),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		fmt.Println("error iniciando sesión con aws: #{err}")
		return sess, err
	}
	return sess, nil
}
