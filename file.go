package gotgbot

import (
	"encoding/json"
	"errors"
	"io"
)

type InputFile struct {
	Name  string
	Data  io.Reader
	Value string
}

func (f InputFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Value)
}

var (
	ErrAttachmentKeyAlreadyExists = errors.New("key already exists")
	ErrEmptyInputFile             = errors.New("empty InputFile")
	ErrMultipleFileReferences     = errors.New("InputFile can have either data, or a value - not both")
)

func (f InputFile) Attach(key string, data map[string]InputFile) (string, error) {
	if f.Data == nil && f.Value == "" {
		return "", ErrEmptyInputFile
	}
	if f.Data != nil && f.Value != "" {
		return "", ErrMultipleFileReferences
	}

	if f.Value != "" {
		return f.Value, nil
	}

	if _, ok := data[key]; ok {
		return "", ErrAttachmentKeyAlreadyExists
	}
	data[key] = f
	return "attach://" + key, nil
}

// InputFileURL is used to send a file on the internet via a publicly accessible HTTP URL.
func InputFileURL(url string) InputFile {
	return InputFile{Value: url}
}

// InputFileID is used to send a file that is already present on telegram's servers, using its telegram file_id.
func InputFileID(fileID string) InputFile {
	return InputFile{Value: fileID}
}

// InputFileReader is used to send a file by a reader interface; such as a filehandle from os.Open().
//
// For example:
//
//	f, err := os.Open("some_file.go")
//	if err != nil {
//		return fmt.Errorf("failed to open file: %w", err)
//	}
//
//	m, err := b.SendDocument(<chat_id>, gotgbot.InputFileReader("source.go", f), <opts>)
func InputFileReader(name string, r io.Reader) InputFile {
	return InputFile{Name: name, Data: r}
}
