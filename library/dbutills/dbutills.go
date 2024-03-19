/*******************************************************************************
* ファイル名: dbutill.go
* 内容: データベースを扱うライブラリ
* 作成日: 2024年3月15日
* 更新日: 2024年3月15日
*******************************************************************************/
package dbutills

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"maia.go/library/fileops"
)

/*--定数の定義-----------------------------------------------------------------*/
const FilePath = "contents/FileInfo.db" /* データベースのファイルパス */

/*--ファイルディレクトリまでのパス--*/
const AnimatedDir = "contents/animated" /* アニメーションファイルのディレクトリ */
const ImageDir = "contents/image"       /* 画像ファイルのディレクトリ */
const MovieDir = "contents/movie"       /* 動画ファイルのディレクトリ */
const MusicDir = "contents/music"       /* 音楽ファイルのディレクトリ */

/*--構造体の定義---------------------------------------------------------------*/
/* HTMLファイルに埋め込むファイル情報 */
type FileInfo struct {
	ID   int      /* ファイルID */
	TYPE string   /* ファイルタイプ */
	Name string   /* ファイル名 */
	Path string   /* ファイルパス */
	Tags []string /* タグ */
}

/*--関数の定義-----------------------------------------------------------------*/
/* Mediaテーブル行をすべて削除する */
func DeleteMedia() {
	/*--初期化---------------------------------------------*/
	db, err := sql.Open("sqlite3", FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/*--データ削除-----------------------------------------*/
	_, err = db.Exec("DELETE FROM Media")
	if err != nil {
		log.Fatal(err)
	}
}

/* データベースの初期化 */
func InitDB_Table() {
	/*--初期化---------------------------------------------*/
	db, err := sql.Open("sqlite3", FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queries := []string{
		`CREATE TABLE MediaType (
            ID INTEGER PRIMARY KEY,
            TYPE TEXT NOT NULL
        );`,
		`CREATE TABLE Media (
            ID INTEGER PRIMARY KEY,
            NAME TEXT NOT NULL,
            MediaType_ID INTEGER,
            FOREIGN KEY(MediaType_ID) REFERENCES MediaType(ID)
        );`,
		`CREATE TABLE Tag (
            ID INTEGER PRIMARY KEY,
            TAG TEXT NOT NULL,
			TAG_JAPANESE TEXT
        );`,
		`CREATE TABLE MediaTag (
            Media_ID INTEGER,
            Tag_ID INTEGER,
            PRIMARY KEY(Media_ID, Tag_ID),
            FOREIGN KEY(Media_ID) REFERENCES Media(ID),
            FOREIGN KEY(Tag_ID) REFERENCES Tag(ID)
        );`,
	}

	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

/* ファイルパス作成関数 */
func GetFilePath(mediaType string, filename string) string {
	var filePath string
	switch mediaType {
	case "image":
		filePath = "../" + ImageDir + "/" + filename
	case "movie":
		filePath = "../" + MovieDir + "/" + filename
	case "animated":
		filePath = "../" + AnimatedDir + "/" + filename
	case "audio":
		filePath = "../" + MusicDir + "/" + filename
	}
	return filePath
}

/* タグ取得関数 */
func GetTags(fileID int) []string {
	/*--初期化---------------------------------------------*/
	db, err := sql.Open("sqlite3", FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/*--データ取得-----------------------------------------*/
	/* クエリの実行 */
	rows, err := db.Query("SELECT Tag.TAG FROM Tag, MediaTag WHERE MediaTag.Media_ID = ? AND MediaTag.Tag_ID = Tag.ID", fileID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	/* タグの取得 */
	var tags []string
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			log.Fatal(err)
		}
		tags = append(tags, tag)
	}

	return tags
}

/* ディレクトリ内の情報をDBに入力する */
func InitDB_Data() {
	/* Imageファイルを取得 */
	FileInfos := fileops.GetFileList(ImageDir)

	/* メディアテーブルに入力 */
	for _, fileInfo := range FileInfos {
		InsertFileInfo(fileInfo.Name, "image", fileInfo.Tags)
	}

	/* Movieファイルを取得 */
	FileInfos = fileops.GetFileList(MovieDir)

	/* メディアテーブルに入力 */
	for _, fileInfo := range FileInfos {
		InsertFileInfo(fileInfo.Name, "movie", fileInfo.Tags)
	}

	/* Animatedファイルを取得 */
	FileInfos = fileops.GetFileList(AnimatedDir)

	/* メディアテーブルに入力 */
	for _, fileInfo := range FileInfos {
		InsertFileInfo(fileInfo.Name, "animated", fileInfo.Tags)
	}
}

func GetFileInfo(AND_Tag []string, OR_Tag []string, NOT_Tag []string) []FileInfo {
	/*--初期化---------------------------------------------*/
	db, err := sql.Open("sqlite3", FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/*--データ取得-----------------------------------------*/
	/* クエリの作成 */
	query := "SELECT Media.ID, MediaType.TYPE, Media.NAME FROM Media, MediaType WHERE Media.MediaType_ID = MediaType.ID"
	var args []interface{}

	if len(AND_Tag) > 0 || len(OR_Tag) > 0 || len(NOT_Tag) > 0 {
		query += " AND Media.ID IN (SELECT Media_ID FROM MediaTag, Tag WHERE MediaTag.Tag_ID = Tag.ID"
	}

	if len(AND_Tag) > 0 {
		for _, tag := range AND_Tag {
			query += " AND Tag.TAG = ?"
			args = append(args, tag)
		}
	}

	if len(OR_Tag) > 0 {
		query += " AND ("
		for i, tag := range OR_Tag {
			if i > 0 {
				query += " OR"
			}
			query += " Tag.TAG = ?"
			args = append(args, tag)
		}
		query += ")"
	}

	if len(NOT_Tag) > 0 {
		for _, tag := range NOT_Tag {
			query += " AND Tag.TAG != ?"
			args = append(args, tag)
		}
	}

	if len(AND_Tag) > 0 || len(OR_Tag) > 0 || len(NOT_Tag) > 0 {
		query += ")"
	}

	/* クエリの実行 */
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	/* FileInfo型にデータを格納 */
	var fileInfos []FileInfo
	for rows.Next() {
		var fileInfo FileInfo
		/* ID, ファイルタイプ, ファイル名の格納 */
		err = rows.Scan(&fileInfo.ID, &fileInfo.TYPE, &fileInfo.Name)
		if err != nil {
			log.Fatal(err)
		}

		/* ファイルパスの作成 */
		fileInfo.Path = GetFilePath(fileInfo.TYPE, fileInfo.Name)

		/* タグの取得 */
		fileInfo.Tags = GetTags(fileInfo.ID)

		/* FileInfoを追加 */
		fileInfos = append(fileInfos, fileInfo)
	}

	// /* FileInfo型をlogに表示 */
	// for _, fileInfo := range fileInfos {
	// 	log.Println(fileInfo)
	// }

	return fileInfos
}

/* 新しいファイルの情報を入力 */
func InsertFileInfo(name string, mediaType string, tags []string) {
	/*--初期化---------------------------------------------*/
	db, err := sql.Open("sqlite3", FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/*--データ挿入-----------------------------------------*/
	/* MediaTypeのIDを取得 */
	var mediaTypeID int
	err = db.QueryRow("SELECT ID FROM MediaType WHERE TYPE = ?", mediaType).Scan(&mediaTypeID)
	if err != nil {
		log.Fatal(err)
	}

	/* Mediaの挿入 */
	_, err = db.Exec("INSERT INTO Media(NAME, MediaType_ID) VALUES(?, ?)", name, mediaTypeID)
	if err != nil {
		log.Fatal(err)
	}

	/* MediaのIDを取得 */
	var mediaID int
	err = db.QueryRow("SELECT ID FROM Media WHERE NAME = ?", name).Scan(&mediaID)
	if err != nil {
		log.Fatal(err)
	}

	/* TagのIDを取得 */
	var tagIDs []int
	for _, tag := range tags {
		var tagID int
		err = db.QueryRow("SELECT ID FROM Tag WHERE TAG = ?", tag).Scan(&tagID)
		if err != nil {
			log.Fatal(err)
		}
		tagIDs = append(tagIDs, tagID)
	}

	/* MediaTagの挿入 */
	for _, tagID := range tagIDs {
		_, err = db.Exec("INSERT INTO MediaTag(Media_ID, Tag_ID) VALUES(?, ?)", mediaID, tagID)
		if err != nil {
			log.Fatal(err)
		}
	}
}

/* タグとメディアの関連付けを行う */
func InsertMediaTag(mediaName string, tags []string) {
	/*--初期化---------------------------------------------*/
	db, err := sql.Open("sqlite3", FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/*--データ挿入-----------------------------------------*/
	/* MediaのIDを取得 */
	var mediaID int
	err = db.QueryRow("SELECT ID FROM Media WHERE NAME = ?", mediaName).Scan(&mediaID)
	if err != nil {
		log.Fatal(err)
	}

	/* TagのIDを取得 */
	var tagIDs []int
	for _, tag := range tags {
		var tagID int
		err = db.QueryRow("SELECT ID FROM Tag WHERE TAG = ?", tag).Scan(&tagID)
		if err != nil {
			log.Fatal(err)
		}
		tagIDs = append(tagIDs, tagID)
	}

	/* MediaTagの挿入 */
	for _, tagID := range tagIDs {
		_, err = db.Exec("INSERT INTO MediaTag(Media_ID, Tag_ID) VALUES(?, ?)", mediaID, tagID)
		if err != nil {
			log.Fatal(err)
		}
	}
}

/* 名前から1ファイルの情報を取得し、FileInfo型に変換して返す */
func GetFileInfoByName(filename string) FileInfo {

	/* DBを開く */
	db, err := sql.Open("sqlite3", FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/* クエリの実行 */
	rows, err := db.Query("SELECT Media.ID, MediaType.TYPE, Media.NAME FROM Media, MediaType WHERE Media.MediaType_ID = MediaType.ID AND Media.NAME = ?", filename)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	/* FileInfo型にデータを格納 */
	var fileinfo FileInfo
	for rows.Next() {
		/* ID, ファイルタイプ, ファイル名の格納 */
		err = rows.Scan(&fileinfo.ID, &fileinfo.TYPE, &fileinfo.Name)
		if err != nil {
			log.Fatal(err)
		}

		/* ファイルパスの作成 */
		fileinfo.Path = GetFilePath(fileinfo.TYPE, fileinfo.Name)

		/* タグの取得 */
		fileinfo.Tags = GetTags(fileinfo.ID)
	}

	/* FileInfo型をlogに表示 */
	log.Println(fileinfo)

	return fileinfo
}
