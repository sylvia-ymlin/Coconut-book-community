package database

import (
	"sync"

	"github.com/sylvia-ymlin/Coconut-book-community/config"
	"github.com/sylvia-ymlin/Coconut-book-community/internal/pkg/storage"
)

var videoSaver *storage.LocalDouyinVedioSaver
var videoInitOnce sync.Once

func GetVideoSaver() storage.VideoStorageService[storage.SimpleObject] {
	videoInitOnce.Do(func() {
		videoSaver = storage.InitLocalOSS(GetMysqlDB(), config.GetVedioConfig().BasePath)
	})
	return videoSaver
}
