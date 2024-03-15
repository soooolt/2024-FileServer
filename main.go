/*******************************************************************************
* ファイル名: main.go
* 内容: バックエンドhttpsサーバー
* 作成日: 2024年3月15日
* 更新日: 2024年3月15日
*******************************************************************************/

package main

import (
	"github.com/gin-gonic/gin"
	"maia.go/library/dbutills"
)

/*--定数の定義-----------------------------------------------------------------*/
const SearchHTML = "page/search.html" /* 検索画面のHTMLファイル */
const ViewHTML = "page/view.html"     /* 閲覧画面のHTMLファイル */
const StaticDir = "static"            /* 静的ファイルのディレクトリ */

/*--構造体の定義---------------------------------------------------------------*/
/* HTMLファイルに埋め込むファイル情報 */
type FileInfo struct {
	ID   int      /* ファイルID */
	TYPE int      /* ファイルタイプ */
	Name string   /* ファイル名 */
	Tag  []string /* タグ */
}

/*--関数の定義-----------------------------------------------------------------*/
func main() {
	dbutills.InitDB()
}

// func main() {
// 	/*--メイン---------------------------------------------*/
// 	router := gin.Default()

// 	/* homeレスポンス */
// 	// router.GET("/", )

// 	/* ファイル検索画面のレスポンス */
// 	search := router.Group("/search")
// 	{
// 		search.GET("/Tag", SearctTagEndPoint)
// 	}

// 	/* ファイル閲覧画面のレスポンス */
// 	view := router.Group("/view")
// 	{
// 		view.GET("/File", ViewFileEndPoint)
// 	}

// 	/* 静的ファイルのレスポンス */
// 	router.Static("/static", StaticDir)

// 	/* サーバ立ち上げ */
// 	go router.Run(":8080")

// 	/*--デバッグ用-----------------------------------------*/
// 	/* URLをコンソールに表示 */
// 	fmt.Println("http://localhost:8080")
// }

/*--ファイル検索画面-----------------------------------------------------------*/
/* 検索なし */
func SearchEndPoint(c *gin.Context) {

}

/* タグ検索 */
func SearctTagEndPoint(c *gin.Context) {

}

/*--ファイル閲覧画面-----------------------------------------------------------*/
/* ファイル閲覧 */
func ViewFileEndPoint(c *gin.Context) {

}
