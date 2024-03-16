/*******************************************************************************
* ファイル名: main.go
* 内容: バックエンドhttpsサーバー
* 作成日: 2024年3月15日
* 更新日: 2024年3月15日
*******************************************************************************/

package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"maia.go/library/dbutills"
)

/*--定数の定義-----------------------------------------------------------------*/
const SearchHTML = "search.html" /* 検索画面のHTMLファイル */
const ViewHTML = "view.html"     /* 閲覧画面のHTMLファイル */
const StaticDir = "static"       /* 静的ファイルのディレクトリ */

/*--構造体の定義---------------------------------------------------------------*/
/* HTMLファイルに埋め込むファイル情報 */
type FileInfo struct {
	ID   int      /* ファイルID */
	TYPE int      /* ファイルタイプ */
	Name string   /* ファイル名 */
	Path string   /* ファイルパス */
	Tags []string /* タグ */
}

/*--関数の定義-----------------------------------------------------------------*/
// dbutills package
/* データベースからHTMLファイルに埋め込むデータを送信する関数 */
// func GetFileInfo(AND_Tag[] string, OR_Tag[] string, NOT_Tag[] string) FileInfo[]

/*--関数の定義-----------------------------------------------------------------*/
/* 単体テスト用 main */
func main__() {
	dbutills.DeleteMedia()
	dbutills.InitDB_Data()
	dbutills.GetFileInfo(nil, nil, nil)
}

/* 結合テスト用 main */
func main() {
	/* テスト用データベース構築 */
	dbutills.DeleteMedia()
	dbutills.InitDB_Data()
	dbutills.GetFileInfo(nil, nil, nil)

	/*--サーバーの設定-------------------------------------*/
	router := gin.Default()

	/* ginにテンプレートを渡す */
	router.LoadHTMLGlob("page/*.html")

	/* "/"の時にindex.htmlを返す */
	router.GET("/", SearchEndPoint)

	/* ファイル検索画面のレスポンス */
	search := router.Group("/search")
	{
		search.GET("/Tag", SearctTagEndPoint)
	}

	/* ファイル閲覧画面のレスポンス */
	view := router.Group("/view")
	{
		view.GET("/File", ViewFileEndPoint)
	}

	/* 静的ファイルのレスポンス */
	router.Static("/static", StaticDir)
	router.Static("/contents", "./contents")

	/*--デバッグ用-----------------------------------------*/
	/* アクセス用のURLをコンソールに表示 */
	fmt.Println("http://localhost:8080/")

	router.Run()
}

/*--ファイル検索画面-----------------------------------------------------------*/
/* 検索なし */
func SearchEndPoint(c *gin.Context) {
	/* データベースから全ファイル情報を取得 */
	fileinfos := dbutills.GetFileInfo(nil, nil, nil)

	/* 画像ファイルのみを取り出す */
	var imagefileinfos []dbutills.FileInfo
	for _, fileinfo := range fileinfos {
		if fileinfo.TYPE == 2 {
			imagefileinfos = append(imagefileinfos, fileinfo)
		}
	}

	/* 画像ファイルの配列の中から20要素だけにする */
	if len(imagefileinfos) > 20 {
		imagefileinfos = imagefileinfos[:20]
	}

	/* HTMLファイルにデータを埋め込む */
	c.HTML(200, SearchHTML, gin.H{
		"fileinfos": imagefileinfos,
	})
}

/* タグ検索 */
func SearctTagEndPoint(c *gin.Context) {

}

/*--ファイル閲覧画面-----------------------------------------------------------*/
/* ファイル閲覧 */
func ViewFileEndPoint(c *gin.Context) {

}
