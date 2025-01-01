package pkg

import (
	"os/user"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func GetExpandedFile(filename string) string {
	usr, err := user.Current()
	if err != nil {
		log.Err(err).Str("filename", filename).Msg("could not expand file")
		return filename
	}

	dir := usr.HomeDir

	if strings.HasPrefix(filename, "~/") {
		return filepath.Join(dir, filename[2:])
	}

	if strings.HasPrefix(filename, "$HOME/") {
		return filepath.Join(dir, filename[6:])
	}

	return filename
}
