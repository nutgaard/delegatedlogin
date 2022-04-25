package copy_process

import (
	cp "github.com/otiai10/copy"
	"github.com/rs/zerolog/log"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

func CopyAndProcessFiles(src, dest string, processor func([]byte) []byte) {
	err := cp.Copy(src, dest, cp.Options{OnDirExists: func(src, dest string) cp.DirExistsAction {
		return cp.Replace
	}})
	if err != nil {
		log.Err(err).Msgf("Could not create tmp folder")
	}

	err = filepath.Walk("/tmp/www", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Err(err).Msgf("Could not walk tmp folder")
			return err
		}

		if info.IsDir() {
			return nil
		}

		return processFile(path, processor)
	})
	if err != nil {
		log.Err(err).Msgf("Could not walk tmp folder")
	}
}

func processFile(path string, processor func([]byte) []byte) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Err(err).Msgf("Could not read file %s", path)
		return err
	}

	content = processor(content)
	err = ioutil.WriteFile(path, content, 0444)
	if err != nil {
		log.Err(err).Msgf("Could not write content to file %s", path)
		return err
	}
	return nil
}
