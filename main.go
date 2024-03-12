package main

import (
    "fmt"
	"bytes"
	"errors"
    "bufio"
	"io"
	"os"
	"strconv"
)

type PGMImage struct {
    Width int
    Height int
    MaxVal int
    Pixels [][]uint8
} 

func ReadPGM(filePath string) (*PGMImage, error) {
    if len(filePath) < 5 {
        return nil, errors.New("Error: Invalid file name")
    }
    fileType := filePath[len(filePath)-4:]
    if fileType != ".pgm" {
        return nil, errors.New("Error: File type must be .pgm")
    }

    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    magicNumber := make([]byte, 2) 
    _, err = file.Read(magicNumber)
    if err != nil {
        return nil, err
    }
    
    p5 := []byte{80, 53}
    p2 := []byte{80, 50}
    if !(bytes.Equal(p5, magicNumber)) {
        if (bytes.Equal(p2, magicNumber)) {
            return nil, errors.New("Error: This program only supports \"P5\" a.k.a. raw pbm images")
        }
        return nil, errors.New("Error: Invalid magic number in pbm header.")
    }

    reader := bufio.NewReader(file)

    width, err := readNextByte(reader, "Error: End of file reached before width could be completely read!")
    if err != nil {
        return nil, err
    }
    
    widthLeftover, err := readUntilWhitespace(reader, "Error: End of file reached before width could be completely read!")
    if err != nil {
        return nil, err
    }

    width = append(width, widthLeftover...)
    widthStr := string(width)
    widthInt, err := strconv.Atoi(widthStr)
    if err != nil {
        return nil, errors.New("Width is not a valid number in ASCII")
    }

    height, err := readNextByte(reader, "Error: End of file reached before height could be completely read!")
    if err != nil {
        return nil, err
    }

    heightLeftover, err := readUntilWhitespace(reader, "Error: End of file reached before height could be completely read!")
    if err != nil {
        return nil, err
    }

    height = append(height, heightLeftover...)
    heightStr := string(height)
    heightInt, err := strconv.Atoi(heightStr)
    if err != nil {
        return nil, errors.New("Height is not a valid number in ASCII")
    }
    
    maxVal, err := readNextByte(reader, "Error: End of file reached before maximum value could be completely read!")
    if err != nil {
        return nil, err
    }

    maxValLeftover, err := readUntilWhitespace(reader, "Error: End of file reached before maximum value could be completely read!")
    if err != nil {
        return nil, err
    }

    maxVal = append(maxVal, maxValLeftover...)
    maxValStr := string(maxVal)
    maxValInt, err := strconv.Atoi(maxValStr) 
    if err != nil {
        return nil, errors.New("Error: Maximum value is not a valid number in ASCII")
    }
    if maxValInt > 255 {
        return nil, errors.New("Error: This program only handles up to a max value of 255")
    }

    image := make([][]uint8, heightInt)
    for i := range image {
        image[i] = make([]uint8, widthInt)
    }

    firstPixelByte, err := readNextByte(reader, "Error: End of file reached before first pixel intensity could be completely read!")
    if err != nil {
        return nil, err
    }
    image[0][0] = uint8(firstPixelByte[0])
    
    for y := 0; y < heightInt; y++ {
        for x := 0; x < widthInt; x++ {
            if ! (y==0 && x==0) {
                nextPixelByte, err := reader.ReadByte()
                if err != nil {
                    if err == io.EOF {
                        return nil, fmt.Errorf("Error: End of file reached before all pixel intensities could be completely read! pixel: %d, %d", x, y)
                    }
                    return nil, err
                }
                image[y][x] = uint8(nextPixelByte)
            }
        }
    }

    finalImage := &PGMImage {
        Width: widthInt,
        Height: heightInt,
        MaxVal: maxValInt,
        Pixels: image,
    }

    return finalImage, nil
}

func readNextByte(reader *bufio.Reader, errMsg string) ([]byte, error) {
    value := make([]byte, 0)
    whitespaceBytes := []byte{9, 10, 11, 12, 13, 32}
    for {
        nextByte, err := reader.ReadByte()
        if err != nil {
            if err == io.EOF {
                return nil, errors.New(errMsg)
            }
            return nil, err
        }
        if !bytes.ContainsRune(whitespaceBytes, rune(nextByte)) {
            value = append(value, nextByte)
            break
        }
    }
    return value, nil
}

func readUntilWhitespace(reader *bufio.Reader, errMsg string) ([]byte, error) {
    values := make([]byte, 0)
    whitespaceBytes := []byte{9, 10, 11, 12, 13, 32}
    for {
        nextByte, err := reader.ReadByte()
        if err != nil {
            if err == io.EOF {
                return nil, errors.New(errMsg)
            }
            return nil, err
        }
        if bytes.ContainsRune(whitespaceBytes, rune(nextByte)) {
            break
        }
        values = append(values, nextByte)
    }
    return values, nil
}
