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
)

/*--定数の定義-----------------------------------------------------------------*/
const FilePath = "contents/FileInfo.db" /* データベースのファイルパス */

/*--構造体の定義---------------------------------------------------------------*/

/*--関数の定義-----------------------------------------------------------------*/
/* データベースの初期化 */
func InitDB() {
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
