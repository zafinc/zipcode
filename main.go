package main

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const (
	CsvUrl       = "https://www.post.japanpost.jp/zipcode/dl/utf/zip/utf_ken_all.zip"
	DownloadFile = "download.zip"
	CsvFile      = "utf_ken_all.csv"
)

type Record struct {
	Zipcode          string `json:"zipcode"`
	PrefectureCode   string `json:"prefectureCode"`
	PrefectureNumber int    `json:"prefectureNumber"`
	Prefecture       string `json:"prefecture"`
	City             string `json:"city"`
	Town             string `json:"town"`
}

var Prefectures = map[string]string{
	"北海道":  "01",
	"青森県":  "02",
	"岩手県":  "03",
	"宮城県":  "04",
	"秋田県":  "05",
	"山形県":  "06",
	"福島県":  "07",
	"茨城県":  "08",
	"栃木県":  "09",
	"群馬県":  "10",
	"埼玉県":  "11",
	"千葉県":  "12",
	"東京都":  "13",
	"神奈川県": "14",
	"新潟県":  "15",
	"富山県":  "16",
	"石川県":  "17",
	"福井県":  "18",
	"山梨県":  "19",
	"長野県":  "20",
	"岐阜県":  "21",
	"静岡県":  "22",
	"愛知県":  "23",
	"三重県":  "24",
	"滋賀県":  "25",
	"京都府":  "26",
	"大阪府":  "27",
	"兵庫県":  "28",
	"奈良県":  "29",
	"和歌山県": "30",
	"鳥取県":  "31",
	"島根県":  "32",
	"岡山県":  "33",
	"広島県":  "34",
	"山口県":  "35",
	"徳島県":  "36",
	"香川県":  "37",
	"愛媛県":  "38",
	"高知県":  "39",
	"福岡県":  "40",
	"佐賀県":  "41",
	"長崎県":  "42",
	"熊本県":  "43",
	"大分県":  "44",
	"宮崎県":  "45",
	"鹿児島県": "46",
	"沖縄県":  "47",
}

func download(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzip(src string, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	defer reader.Close()

	for _, file := range reader.File {
		input, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}

		defer input.Close()

		outputPath := filepath.Join(dest, file.Name)
		output, err := os.Create(outputPath)
		if err != nil {
			log.Fatal(err)
		}

		defer output.Close()

		_, err = io.Copy(output, input)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func main() {
	downloadPath := filepath.Join(os.TempDir(), DownloadFile)
	if err := download(downloadPath, CsvUrl); err != nil {
		log.Fatal(err)
	}

	if err := unzip(downloadPath, os.TempDir()); err != nil {
		log.Fatal(err)
	}

	csvPath := filepath.Join(os.TempDir(), CsvFile)
	fp, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := fp.Close(); err != nil {
			log.Fatal(err)
		}

		if err := os.Remove(downloadPath); err != nil {
			log.Fatal(err)
		}

		if err := os.Remove(csvPath); err != nil {
			log.Fatal(err)
		}
	}()

	reader := csv.NewReader(fp)

	grouped := make(map[string]map[string]Record)

	for {
		columns, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		prefectureCode := Prefectures[columns[6]]
		prefectureId, err := strconv.Atoi(prefectureCode)
		if err != nil {
			log.Fatal(err)
		}

		record := Record{
			Zipcode:          columns[2],
			PrefectureCode:   prefectureCode,
			PrefectureNumber: prefectureId,
			Prefecture:       columns[6],
			City:             columns[7],
			Town:             columns[8],
		}

		head := record.Zipcode[:3]

		//if head != "152" {
		//	continue
		//}

		if grouped[head] == nil {
			grouped[head] = map[string]Record{}
		}

		grouped[head][record.Zipcode] = record
	}

	// write to json file
	for head, records := range grouped {
		path := filepath.Join("./data", head+".json")
		file, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(records); err != nil {
			log.Fatal(err)
		}

		if err := file.Close(); err != nil {
			log.Println(err)
		}

		fmt.Println(path)
	}
}
