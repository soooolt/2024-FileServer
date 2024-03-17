/*******************************************************************************
* ファイル名: main.go
* 内容: バックエンドhttpsサーバー
* 作成日: 2024年3月15日
* 更新日: 2024年3月15日
*******************************************************************************/

package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"maia.go/library/dbutills"
)

/*--定数の定義-----------------------------------------------------------------*/
const SearchHTML = "search.html" /* 検索画面のHTMLファイル */
const ViewHTML = "view.html"     /* 閲覧画面のHTMLファイル */
const StaticDir = "static"       /* 静的ファイルのディレクトリ */

/*--構造体の定義---------------------------------------------------------------*/

/* ユーザーのフィルタ状態 */
type FilterPanel struct {
	FP_movie    bool /* 動画ファイル表示許可 */
	FP_animated bool /* 動く画像ファイル表示許可 */
	FP_image    bool /* 画像ファイル表示許可 */
	FP_audio    bool /* 音声ファイル表示許可 */
	FP_manga    bool /* 漫画ファイル表示許可 */
	FP_3D       bool /* 3Dタグファイル表示許可 */
	FP_2D       bool /* 2Dタグファイル表示許可 */
	FP_Real     bool /* Realタグファイル表示許可 */
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
		search.POST("/filter", FilterPanelEndPoint)
		search.GET("/filter", SearchEndPoint)
	}

	/* ファイル閲覧画面のレスポンス */
	view := router.Group("/view")
	{
		view.GET("/", ViewFileEndPoint)
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

	/* ページ番号を取得（デフォルトは1） */
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	/* ページあたりの要素数 */
	perPage := 20

	/* ファイル情報をページに応じて取り出す */
	start := (page - 1) * perPage
	end := start + perPage
	if end > len(fileinfos) {
		end = len(fileinfos)
	}
	imagefileinfos := fileinfos[start:end]

	/* logにFileinfoを出力 */
	fmt.Println(imagefileinfos)

	c.HTML(200, SearchHTML, gin.H{
		"fileinfos": imagefileinfos,
		"prevPage":  page - 1,
		"nextPage":  page + 1,
		"hasPrev":   page > 1,
		"hasNext":   end < len(fileinfos),
		"total":     len(fileinfos), // ファイル情報の総数
	})
}

/* タグ検索 */
func SearctTagEndPoint(c *gin.Context) {
	/* クエリパラメータを取得 */
	AND_Tag := c.QueryArray("AND_Tag")
	OR_Tag := c.QueryArray("OR_Tag")
	NOT_Tag := c.QueryArray("NOT_Tag")

	/* データベースからファイル情報を取得 */
	fileinfos := dbutills.GetFileInfo(AND_Tag, OR_Tag, NOT_Tag)

	/* ページ番号を取得（デフォルトは1） */
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	/* ページあたりの要素数 */
	perPage := 20

	/* ファイル情報をページに応じて取り出す */
	start := (page - 1) * perPage
	end := start + perPage
	if end > len(fileinfos) {
		end = len(fileinfos)
	}
	imagefileinfos := fileinfos[start:end]

	/* logにFileinfoを出力 */
	fmt.Println(imagefileinfos)

	c.HTML(200, SearchHTML, gin.H{
		"fileinfos": imagefileinfos,
		"prevPage":  page - 1,
		"nextPage":  page + 1,
		"hasPrev":   page > 1,
		"hasNext":   end < len(fileinfos),
		"total":     len(fileinfos), // ファイル情報の総数
	})
}

/* フィルタ情報から取得情報をフィルタして返す */
func FileInfoFilter(fileinfo []dbutills.FileInfo, filterPanel FilterPanel) []dbutills.FileInfo {
	var result []dbutills.FileInfo
	for _, info := range fileinfo {
		/* フィルタパネルの状態に応じてフィルタ */
		if filterPanel.FP_movie && info.TYPE == "movie" {
			result = append(result, info)
		}
		if filterPanel.FP_animated && info.TYPE == "animated" {
			result = append(result, info)
		}
		if filterPanel.FP_image && info.TYPE == "image" {
			result = append(result, info)
		}
		if filterPanel.FP_audio && info.TYPE == "audio" {
			result = append(result, info)
		}
		if filterPanel.FP_manga && info.TYPE == "manga" {
			result = append(result, info)
		}
		/* タグを用いたフィルタ */
		/*--未実装--*/
	}
	return result
}

/* フィルタパネル */
func FilterPanelEndPoint(c *gin.Context) {
	/* フィルタパネルの状態を取得 */
	filterPanel := FilterPanel{
		FP_movie:    c.PostForm("FP_movie") == "true",
		FP_animated: c.PostForm("FP_animated") == "true",
		FP_image:    c.PostForm("FP_image") == "true",
		FP_audio:    c.PostForm("FP_audio") == "true",
		FP_manga:    c.PostForm("FP_manga") == "true",
		FP_3D:       c.PostForm("FP_3D") == "true",
		FP_2D:       c.PostForm("FP_2D") == "true",
		FP_Real:     c.PostForm("FP_Real") == "true",
	}
	/* logにFilterPanelの情報を出力 */
	fmt.Println(filterPanel)

	/* データベースからファイル情報を取得 */
	filebuf := dbutills.GetFileInfo(nil, nil, nil)

	/* フィルタパネルの状態に応じてフィルタ */
	fileinfos := FileInfoFilter(filebuf, filterPanel)

	/* ページ番号を取得（デフォルトは1） */
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	/* ページあたりの要素数 */
	perPage := 20

	/* ファイル情報をページに応じて取り出す */
	start := (page - 1) * perPage
	end := start + perPage
	if end > len(fileinfos) {
		end = len(fileinfos)
	}
	imagefileinfos := fileinfos[start:end]

	/* logにFileinfoを出力 */
	fmt.Println(imagefileinfos)

	c.HTML(200, SearchHTML, gin.H{
		"fileinfos": imagefileinfos,
		"prevPage":  page - 1,
		"nextPage":  page + 1,
		"hasPrev":   page > 1,
		"hasNext":   end < len(fileinfos),
		"total":     len(fileinfos), // ファイル情報の総数
	})
}

/*--ファイル閲覧画面-----------------------------------------------------------*/
/* ファイル閲覧 */
func ViewFileEndPoint(c *gin.Context) {
	/* リクエストからファイル名を取得 */
	filename := c.Query("filename")

	/* データベースからファイル情報を取得 */
	fileinfo := dbutills.GetFileInfoByName(filename)

	/* ファイル情報と一緒にHTMLをレンダリング */
	c.HTML(200, "view.html", gin.H{
		"TYPE": fileinfo.TYPE,
		"Path": fileinfo.Path,
		"Name": fileinfo.Name,
		"Tags": fileinfo.Tags,
	})
}
