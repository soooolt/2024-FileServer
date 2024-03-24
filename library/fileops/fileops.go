/*******************************************************************************
* ファイル名: fileops.go
* 内容: ファイルを扱うライブラリ
* 作成日: 2024年3月15日
* 更新日: 2024年3月15日
*******************************************************************************/
package fileops

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
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

/* ファイルの種類 */
var Extension_String = [][]string{
	{"mp4"},                /* ID:1 = movie*/
	{"gif", "webp", "avi"}, /* ID:2 = animated*/
	{"png", "jpeg", "jpg"}, /* ID:3 = image*/
	{"mp3"},                /* ID:4 = audio*/
}

/*--定数の定義-----------------------------------------------------------------*/
const AnimatedDir = "contents/animated" /* アニメーションファイルのディレクトリ */
const ImageDir = "contents/image"       /* 画像ファイルのディレクトリ */
const MovieDir = "contents/movie/file"  /* 動画ファイルのディレクトリ */
const MusicDir = "contents/music"       /* 音楽ファイルのディレクトリ */

const FrameDir = "contents/movie/frame" /* フレームファイルのディレクトリ */

/*--関数の定義-----------------------------------------------------------------*/
/* DBファイルを削除 */
func DeleteDB() {
	/* DBファイルを削除 */
	err := os.Remove("contents/FileInfo.db")
	if err != nil {
		log.Fatal(err)
	}
}

/* ディレクトリ内のファイルタイプと名前を取得する */
func GetFileList(dir string) []FileInfo {
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

/* ファイルタイプを判定 */
func GetFileType(filename string) string {
	/* ファイルタイプを取得 */
	for i, extension := range Extension_String {
		for _, ext := range extension {
			/* Extension_Stringの要素内の拡張子があった場合にタイプを返す */
			if strings.HasSuffix(filename, ext) {
				switch i {
				case 0:
					return "movie"
				case 1:
					return "animated"
				case 2:
					return "image"
				case 3:
					return "audio"
				}
			}
		}
	}
	return ""
}

/* ファイル保存場所指定 */
func FileNametoDir(filename string) (Path string, Type string) {
	/* ファイルタイプを判定 */
	filetype := GetFileType(filename)

	var filepath string
	/* ファイルタイプごとに保存ディレクトリを作成して返す */
	switch filetype {
	case "movie":
		filepath = MovieDir + "/" + filename
	case "animated":
		filepath = AnimatedDir + "/" + filename
	case "image":
		filepath = ImageDir + "/" + filename
	case "audio":
		filepath = MusicDir + "/" + filename
	}

	return filepath, filetype
}

/* アップロードされたファイルを保存 */
func SaveUploadFile(file multipart.File, handler *multipart.FileHeader) {
	/* ファイル名を取得 */
	filename := handler.Filename

	/* 保存先を判定してパスを受け取る */
	filepath, filetype := FileNametoDir(filename)

	/* ファイルを保存 */
	out, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	/* movieファイルのみサムネイル作成 */
	if filetype == "movie" {
		SaveMovieFrame(filepath, "1")
	}
}

func SaveMovieFrame(videoName string, framenum string) string {
	/* 画像の保存先パスを作成 */
	imagePath := FrameDir + "/" + videoName + ".jpg"

	/* 動画ファイルのパスを作成 */
	videoPath := MovieDir + "/" + videoName

	/* ffmpegでフレームを取得 */
	cmd := exec.Command("ffmpeg", "-y", "-i", videoPath, "-vframes", framenum, imagePath)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to extract first frame from video: %s", videoName)
		log.Printf("Failed with error: %v", err)
		log.Printf("ffmpeg output:\n %v", out.String())
		log.Printf("ffmpeg error:\n %v", stderr.String())
	}

	return imagePath
}
