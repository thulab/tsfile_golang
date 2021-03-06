package main

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"tsfile/common/constant"
	"tsfile/encoding/decoder"
	"tsfile/timeseries/read"

	"tsfile/timeseries/read/reader/impl/basic"
)

func TestRead(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Error:", err)
		}
	}()

	file := "D:/test.ts"
	f := new(read.TsFileSequenceReader)
	f.Open(file)
	defer f.Close()

	headerString := f.ReadHeadMagic()
	log.Println("Header string: " + headerString)

	tailerString := f.ReadTailMagic()
	log.Println("Tail string: " + tailerString)

	fileMetadata := f.ReadFileMetadata()
	log.Println("File version: " + strconv.Itoa(fileMetadata.GetCurrentVersion()))

	for f.HasNextRowGroup() {
		groupHeader := f.ReadRowGroupHeader()
		log.Println("row group: " + groupHeader.GetDevice() + ", chunk number: " + strconv.Itoa(int(groupHeader.GetNumberOfChunks())) + ", end posistion: " + strconv.FormatInt(f.Pos(), 10))
		for i := 0; i < int(groupHeader.GetNumberOfChunks()); i++ {
			chunkHeader := f.ReadChunkHeader()
			log.Println("  chunk: " + chunkHeader.GetSensor() + ", page number: " + strconv.Itoa(chunkHeader.GetNumberOfPages()) + ", end posistion: " + strconv.FormatInt(f.Pos(), 10))
			defaultTimeDecoder := decoder.CreateDecoder(constant.PLAIN, constant.INT64)
			valueDecoder := decoder.CreateDecoder(chunkHeader.GetEncodingType(), chunkHeader.GetDataType())
			for j := 0; j < chunkHeader.GetNumberOfPages(); j++ {
				pageHeader := f.ReadPageHeader(chunkHeader.GetDataType())
				log.Println("    page dps: " + strconv.Itoa(int(pageHeader.GetNumberOfValues())) + ", page data size: " + strconv.Itoa(int(pageHeader.GetCompressedSize())) + ", end posistion: " + strconv.FormatInt(f.Pos(), 10))

				pageData := f.ReadPage(pageHeader, chunkHeader.GetCompressionType())
				reader1 := &basic.PageDataReader{DataType: chunkHeader.GetDataType(), ValueDecoder: valueDecoder, TimeDecoder: defaultTimeDecoder}
				reader1.Read(pageData)
				for reader1.HasNext() {
					pair, _ := reader1.Next()
					log.Println("      (time,value): " + strconv.FormatInt(pair.Timestamp, 10) + ", " + fmt.Sprintf("%v", pair.Value))
				}
			}
		}
	}
}
