package main

import (
	"fmt"
	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/go-audio/wav"
	"log"
	"os"
)

func ProcessSamples(model whisper.Model, samples []float32) (finalText string, err error) {
	// TODO: Fix err handling

	// Process samples
	context, err := model.NewContext()
	if err != nil {
		return finalText, err
	}

	if err := context.SetLanguage("auto"); err != nil {
		log.Printf("failed to set language to auto-detect")
		return finalText, err
	}

	if err := context.Process(samples, nil); err != nil {
		return finalText, err
	}

	// Print out the results
	fmt.Printf("Recognized: ")
	for {
		segment, err := context.NextSegment()
		if err != nil {
			break
		}
		fmt.Printf(" %s ", segment.Text)
		finalText += fmt.Sprintf(" %v", segment.Text)
	}

	return finalText, err

}

func GetSamplesFromFilePath(path string) (samples []float32, err error) {
	fmt.Printf("Loading %q\n", path)
	fh, err := os.Open(path)
	if err != nil {
		log.Print(err)
	}
	defer fh.Close()
	dec := wav.NewDecoder(fh)
	if buf, err := dec.FullPCMBuffer(); err != nil {
		log.Print(err)
	} else if dec.SampleRate != whisper.SampleRate {
		log.Printf("unsupported sample rate: %d", dec.SampleRate)
	} else if dec.NumChans != 1 {
		log.Printf("unsupported number of channels: %d", dec.NumChans)
	} else {
		samples = buf.AsFloat32Buffer().Data
	}
	fmt.Printf("Loaded %q, no. of samples = %d\n", path, len(samples))
	return samples, err
}

func GetModel() whisper.Model {
	var modelPath = "./whisper.cpp/bindings/go/models/ggml-small.en.bin"

	customPath, ok := os.LookupEnv("MODELPATH")
	if ok {
		log.Printf("MODELPATH set to %v\n", customPath)
		modelPath = customPath // If we found the env variable, set it, otherwise we will leave the default
	}

	// Load the model
	model, err := whisper.New(modelPath)
	if err != nil {
		panic(err)
	}
	return model
}
