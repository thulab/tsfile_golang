package tsFileWriter

/**
 * @Package Name: pageWriter
 * @Author: steven yao
 * @Email:  yhp.linux@gmail.com
 * @Create Date: 18-8-31 下午3:51
 * @Description:
 */

import (
	"bytes"
	"encoding/binary"
	"tsfile/common/constant"
	"tsfile/common/log"
	"tsfile/compress"
	"tsfile/file/header"
	"tsfile/file/metadata/statistics"
	"tsfile/timeseries/write/sensorDescriptor"
)

type PageWriter struct {
	compressor      *compress.Encompress
	desc            *sensorDescriptor.SensorDescriptor
	pageBuf         *bytes.Buffer
	totalValueCount int64
	maxTimestamp    int64
	minTimestamp    int64
}

func (p *PageWriter) WritePageHeaderAndDataIntoBuff(dataBuffer *bytes.Buffer, valueCount int, sts statistics.Statistics, maxTimestamp int64, minTimestamp int64) int {
	if p.desc.GetCompresstionType() == int16(constant.UNCOMPRESSED) {
		//this uncompressedSize should be calculate from timeBuf and valueBuf
		uncompressedSize := dataBuffer.Len()
		var compressedSize int = uncompressedSize
		pageHeader, pageHeaderErr := header.NewPageHeader(
			int32(uncompressedSize), int32(compressedSize),
			int32(valueCount), sts, maxTimestamp,
			minTimestamp, p.desc.GetTsDataType())
		if pageHeaderErr != nil {
			log.Error("init pageHeader error: ", pageHeaderErr)
		}
		pageBuf := p.pageBuf
		//pageHeader.PageHeaderToMemory(p.pageBuf, p.desc.GetTsDataType())
		binary.Write(pageBuf, binary.BigEndian, pageHeader.GetUncompressedSize())
		binary.Write(pageBuf, binary.BigEndian, pageHeader.GetCompressedSize())
		binary.Write(pageBuf, binary.BigEndian, pageHeader.GetNumberOfValues())
		binary.Write(pageBuf, binary.BigEndian, pageHeader.Max_timestamp())
		binary.Write(pageBuf, binary.BigEndian, pageHeader.Min_timestamp())
		statistics.Serialize(*(pageHeader.GetStatistics()), pageBuf, p.desc.GetTsDataType())

		pageBuf.Write(dataBuffer.Bytes())
		p.totalValueCount += int64(valueCount)
	} else {
		//this uncompressedSize should be calculate from timeBuf and valueBuf
		uncompressedSize := dataBuffer.Len()

		// write pageData to pageBuf
		//声明一个空的slice,容量为dataBuffer的长度
		dataSlice := make([]byte, dataBuffer.Len())
		//把buf的内容读入到timeSlice内,因为timeSlice容量为timeSize,所以只读了timeSize个过来
		dataBuffer.Read(dataSlice)

		var compressedSize int
		var enc []byte
		aSlice := make([]byte, 0)
		enc = p.compressor.GetEncompressor(p.desc.GetCompresstionType()).Encompress(aSlice, dataSlice)
		compressedSize = len(enc)

		pageHeader, pageHeaderErr := header.NewPageHeader(int32(uncompressedSize), int32(compressedSize), int32(valueCount), sts, maxTimestamp, minTimestamp, p.desc.GetTsDataType())
		if pageHeaderErr != nil {
			log.Error("init pageHeader error: ", pageHeaderErr)
		}
		// write pageheader to pageBuf
		//log.Info("start to flush a page header into buffer, buf pos: %d", p.pageBuf.Len())
		pageHeader.PageHeaderToMemory(p.pageBuf, p.desc.GetTsDataType())
		//log.Info("pageHeader: %v", pageHeader)
		//log.Info("finished to flush a page header into buffer, buf pos: %d", p.pageBuf.Len())

		//// write pageData to pageBuf
		////声明一个空的slice,容量为dataBuffer的长度
		//dataSlice := make([]byte, dataBuffer.Len())
		////把buf的内容读入到timeSlice内,因为timeSlice容量为timeSize,所以只读了timeSize个过来
		//dataBuffer.Read(dataSlice)
		//log.Info("start to flush a page data into buffer, buf pos: %d", p.pageBuf.Len())
		p.pageBuf.Write(enc)
		//log.Info("finished to flush a page data into buffer, buf pos: %d", p.pageBuf.Len())
		p.totalValueCount += int64(valueCount)
	}
	return 0
}

func (p *PageWriter) WriteAllPagesOfSeriesToTsFile(tsFileIoWriter *TsFileIoWriter, seriesStatistics statistics.Statistics, numOfPage int) int64 {
	if p.minTimestamp == -1 {
		log.Error("Write page error, minTime: %s, maxTime: %s")
	}
	// write trunk header to file
	chunkHeaderSize := tsFileIoWriter.StartFlushChunk(p.desc, p.desc.GetCompresstionType(), p.desc.GetTsDataType(), p.desc.GetTsEncoding(), seriesStatistics, p.maxTimestamp, p.minTimestamp, p.pageBuf.Len(), numOfPage)
	preSize := tsFileIoWriter.GetPos()
	// write all pages to file
	tsFileIoWriter.WriteBytesToFile(p.pageBuf)
	//// after write page, reset pageBuf
	//p.pageBuf.Reset()
	dataSize := tsFileIoWriter.GetPos() - preSize
	chunkSize := int64(chunkHeaderSize) + dataSize
	tsFileIoWriter.EndChunk(chunkSize, p.totalValueCount)
	return chunkSize
}

func (p *PageWriter) Reset() {
	p.minTimestamp = -1
	p.pageBuf.Reset()
	p.totalValueCount = 0
	return
}

func (p *PageWriter) EstimateMaxPageMemSize() int {
	pageSize := p.pageBuf.Len()
	pageHeaderSize := header.CalculatePageHeaderSize(p.desc.GetTsDataType())
	return pageSize + pageHeaderSize
}

func (p *PageWriter) GetCurrentDataSize() int {
	size := p.pageBuf.Len()
	return size
}

func NewPageWriter(sd *sensorDescriptor.SensorDescriptor) (*PageWriter, error) {
	return &PageWriter{
		desc:       sd,
		compressor: sd.GetCompressor(),
		pageBuf:    bytes.NewBuffer([]byte{}),
	}, nil
}
