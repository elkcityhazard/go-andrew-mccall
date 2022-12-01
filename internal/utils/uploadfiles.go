package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Tools struct {
	MaxFileSize      int
	AllowedFileTypes []string
}

type UploadedFile struct {
	NewFileName      string
	OriginalFileName string
	FileSize         int64
}

// the random string source

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"

//	RandomString returns a string of random characters of length n using random string source

func (t *Tools) RandomString(n int) string {
	// define two variables
	//	s is a slice of rune of n length
	//  r uses the random string source to seed a slice of runes
	s, r := make([]rune, n), []rune(randomStringSource)

	// range over s which is a fixed length slice
	// use rand crypto package Prime method
	// generate a random number with x,y values where
	// x gets converted from a big int to a Uint64, and y is normalized to uint64 from the length of the random string source slice
	//	finally update each index of the slice of runes with the index of r based on the modulus of x%y

	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}

	// return the random string

	return string(s)
}

func (t *Tools) CreateDirectoryIfNone(r *http.Request, dirname string) (string, error) {
	const mode = 0655
	if _, err := os.Stat(path.Join("./static/uploads/", dirname)); os.IsNotExist(err) {
		err := os.MkdirAll(path.Join("./static/uploads/", dirname), 0655)

		if err != nil {
			return "", err
		}

	}
	return path.Join("./static/uploads/", dirname), nil
}

func (t *Tools) UploadSingleFile(r *http.Request, uploadDir string, rename ...bool) (*UploadedFile, error) {

	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}

	files, err := t.UploadFiles(r, uploadDir, renameFile)
	if err != nil {
		return nil, err
	}

	if len(files) > 0 {
		return files[0], nil
	} else {
		return nil, err
	}

}

func (t *Tools) UploadFiles(r *http.Request, uploadDir string, rename ...bool) ([]*UploadedFile, error) {
	// handle whether or not we are going to rename the file
	// variatic parameter rename

	renameFile := true

	if len(rename) > 0 {
		renameFile = rename[0]
	}

	// uploadedFiles is what is going to be returned
	var uploadedFiles []*UploadedFile

	// set a sensible default if none is set

	if t.MaxFileSize == 0 {
		t.MaxFileSize = 1024 * 1024 * 1024
	}

	_, err := t.CreateDirectoryIfNone(r, uploadDir)

	if err != nil {
		return nil, err
	}

	// Parse the multipart form

	err = r.ParseMultipartForm(int64(t.MaxFileSize))

	fmt.Println(r.MultipartForm.File)

	if err != nil {
		return nil, errors.New("the uploaded file is too big")
	}

	// loop through the multipart form files
	// then loop through the headers for each file
	// anytime you are deferring something in a loop,
	// you need to inline a function

	for _, fheaders := range r.MultipartForm.File {
		for _, hdr := range fheaders {

			uploadedFiles, err := func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {

				//	have a place to store uploaded File
				var uploadedFile UploadedFile

				// open the file

				infile, err := hdr.Open()
				//	check for an error
				if err != nil {
					return nil, err
				}

				// defer closing the file

				defer infile.Close()

				// create a buffer of 512 bytes (to get mime/type)

				buff := make([]byte, 512)

				// read into the buffer

				_, err = infile.Read(buff)

				// check for error again

				if err != nil {
					return nil, err
				}

				// TODO: check to see if the file type is permitted

				// assume we don't want the file type to be uploaded

				allowed := false

				// detect the filetype

				fileType := http.DetectContentType(buff)

				if err != nil {
					return nil, err
				}

				// define some hardcoded allowed file types

				//t.AllowedFileTypes

				// check to see if allowed filetypes slice is populated

				if len(t.AllowedFileTypes) > 0 {
					//	range through the file types and check the equality of the fileType to the allowedTypes
					for _, x := range t.AllowedFileTypes {
						if strings.EqualFold(fileType, x) {
							//	if matches, set allowed true to let the server know the filetype is permitted
							allowed = true
						}
					}
				} else {
					// if file types aren't defined, allow all files
					allowed = true
				}

				if !allowed {
					return nil, errors.New("uploaded filetype is not permitted")
				}

				// handle whether or not we are renaming the file

				if renameFile {
					uploadedFile.NewFileName = fmt.Sprintf("%s%s", t.RandomString(25), filepath.Ext(hdr.Filename))
				} else {
					uploadedFile.NewFileName = hdr.Filename
				}
				uploadedFile.OriginalFileName = hdr.Filename

				// initiate the file that is going to be written to disk

				var outfile *os.File

				//	defer closing the outfile until the function exits
				defer outfile.Close()

				//	reset the file to the beginning of the file

				_, err = infile.Seek(0, 0)

				//	handle joining the uploadDir and the filepath directory
				if outfile, err = os.Create(filepath.Join("./static/uploads/"+uploadDir, uploadedFile.NewFileName)); err != nil {
					return nil, err
				} else {
					// get the filesize
					fileSize, err := io.Copy(outfile, infile)

					if err != nil {
						return nil, err
					}
					uploadedFile.FileSize = fileSize
				}
				// append teh current uploadedFile to the uploadedFiles slice
				uploadedFiles = append(uploadedFiles, &uploadedFile)

				return uploadedFiles, nil
			}(uploadedFiles)
			if err != nil {
				return uploadedFiles, err
			}

		}
	}
	return uploadedFiles, nil
}
