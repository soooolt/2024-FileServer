/*******************************************************************************
* ファイル名: fileops.go
* 内容: ファイルを扱うライブラリ
* 作成日: 2024年3月15日
* 更新日: 2024年3月15日
*******************************************************************************/
package fileops

import (
	"io/ioutil"
	"log"
	"strings"
)

/*--構造体の定義---------------------------------------------------------------*/
/* HTMLファイルに埋め込むファイル情報 */
type FileInfo struct {
	ID   int      /* ファイルID */
	TYPE string   /* ファイルタイプ */
	Name string   /* ファイル名 */
	Path string   /* ファイルパス */
	Tags []string /* タグ */
}

/*--定数の定義-----------------------------------------------------------------*/
const AnimatedDir = "contents/animated" /* アニメーションファイルのディレクトリ */
const ImageDir = "contents/image"       /* 画像ファイルのディレクトリ */
const MovieDir = "contents/movie"       /* 動画ファイルのディレクトリ */
const MusicDir = "contents/music"       /* 音楽ファイルのディレクトリ */

/*--関数の定義-----------------------------------------------------------------*/
/* ディレクトリ内のファイルタイプと名前を取得する */
func GetFileList(dir string) []FileInfo {
	var Extension_String [][]string = [][]string{
		{"mp4"},                /* ID:1 = movie*/
		{"gif", "webp"},        /* ID:2 = animated*/
		{"png", "jpeg", "jpg"}, /* ID:3 = image*/
		{"mp3"},                /* ID:4 = audio*/
	}

	/* ディレクトリ内のファイルを取得 */
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	/* ファイル情報を取得 */
	var fileinfoList []FileInfo
	for _, file := range files {
		var fileinfo FileInfo
		fileinfo.Name = file.Name()
		fileinfo.Path = dir + "/" + file.Name()

		/* ファイルタイプを取得 */
		for i, extension := range Extension_String {
			for _, ext := range extension {
				if strings.HasSuffix(fileinfo.Name, ext) {
					switch i {
					case 1:
						fileinfo.TYPE = "movie"
					case 2:
						fileinfo.TYPE = "animated"
					case 3:
						fileinfo.TYPE = "image"
					case 4:
						fileinfo.TYPE = "audio"
					}
				}
			}
		}

		/* ファイル情報を追加 */
		fileinfoList = append(fileinfoList, fileinfo)
	}

	return fileinfoList
}
