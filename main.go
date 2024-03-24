/*******************************************************************************
* ファイル名: main.go
* 内容: バックエンドhttpsサーバー
* 作成日: 2024年3月15日
* 更新日: 2024年3月15日
*******************************************************************************/

package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"maia.go/library/dbutills"
	"maia.go/library/fileops"
)

/*--定数の定義-----------------------------------------------------------------*/
const IndexHTML = "page.html"    /* 初期画面のHTMLファイル */
const SearchHTML = "search.html" /* 検索画面のHTMLファイル */
const ViewHTML = "view.html"     /* 閲覧画面のHTMLファイル */
const StaticDir = "static"       /* 静的ファイルのディレクトリ */

/*--構造体の定義---------------------------------------------------------------*/

/* タグ検索フォームの状態 */
type TagSearchForm struct {
	AND_Tag []string /* AND検索用タグ */
	OR_Tag  []string /* OR検索用タグ */
	NOT_Tag []string /* NOT検索用タグ */
}

/* ユーザーのフィルタパネルの状態 */
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

/* グローバル変数定義 */
var fileinfos []dbutills.FileInfo /* ファイル情報 */

/*--関数の定義-----------------------------------------------------------------*/
// dbutills package
/* データベースからHTMLファイルに埋め込むデータを送信する関数 */
// func GetFileInfo(AND_Tag[] string, OR_Tag[] string, NOT_Tag[] string) FileInfo[]

/*--関数の定義-----------------------------------------------------------------*/
/* 単体テスト用 main */
func main_() {
	/* テスト用データベース構築 */
	fileops.DeleteDB()
	dbutills.InitDB_Table()
	dbutills.InitDB_MediaType()
	dbutills.InitDB_Data()
}

/* 結合テスト用 main */
func main() {
	/* テスト用データベース構築 */
	// fileops.DeleteDB()
	// dbutills.InitDB_Table()
	// dbutills.InitDB_MediaType()
	// dbutills.InitDB_Data()

	/*--サーバーの設定-------------------------------------*/
	router := gin.Default()

	/* ginにテンプレートを渡す */
	router.LoadHTMLGlob("page/*.html")

	/* "/"の時にindex.htmlを返す */
	router.GET("/", IndexEndPoint)

	/* ファイル検索画面のレスポンス */
	search := router.Group("/search")
	{
		search.POST("/", SearchEndPoint)

		/* ページネーション */
		search.GET("/", paginationEndPoint)
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

/*--レスポンス以外の処理関数----------------------------------------------------*/
/* ページ用データの抽出 */
func ExtractionPage(fileinfos []dbutills.FileInfo, page int) ([]dbutills.FileInfo, int) {
	/* ページあたりの要素数 */
	perPage := 20

	/* ファイル情報をページに応じて取り出す */
	start := (page - 1) * perPage
	end := start + perPage
	if end > len(fileinfos) {
		end = len(fileinfos)
	}
	imagefileinfos := fileinfos[start:end]

	// /* logにFileinfoを出力 */
	// fmt.Println(imagefileinfos)

	return imagefileinfos, end
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

/* FileInfoをシャッフルする */
func ShuffleFileInfo(fileinfo []dbutills.FileInfo) []dbutills.FileInfo {
	/* 要素をシャッフルする---*/
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(fileinfo), func(i, j int) {
		fileinfo[i], fileinfo[j] = fileinfo[j], fileinfo[i]
	})

	/* シャッフルしたfileinfoを返す */
	return fileinfo
}

/* 日付、人気、ランダムでソート */
func SortFileInfo(fileinfo []dbutills.FileInfo, sort string) []dbutills.FileInfo {
	/* 日付でソート */
	if sort == "Date" {
		/*--未実装--*/
	}
	/* 人気でソート */
	if sort == "Popular" {
		/*--未実装--*/
	}
	/* ランダムでソート */
	if sort == "Random" {
		fileinfo = ShuffleFileInfo(fileinfo)
	}
	return fileinfo
}

/*--ファイル検索画面-----------------------------------------------------------*/
/* 初期画面のレスポンス */
func IndexEndPoint(c *gin.Context) {
	/* データベースからファイル情報を取得 */
	fileinfos = dbutills.GetFileInfo(nil, nil, nil)

	/* ページ番号を取得（デフォルトは1） */
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	/* ページデータ取得 */
	imagefileinfos, end := ExtractionPage(fileinfos, page)

	/* HTMLを返す */
	c.HTML(200, IndexHTML, gin.H{
		"fileinfos": imagefileinfos,
		"prevPage":  page - 1,
		"nextPage":  page + 1,
		"hasPrev":   page > 1,
		"hasNext":   end < len(fileinfos),
		"total":     len(fileinfos), // ファイル情報の総数
	})
}

/* ページネーション */
func paginationEndPoint(c *gin.Context) {
	/* ページ番号を取得（デフォルトは1） */
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	/* ページデータ取得 */
	imagefileinfos, end := ExtractionPage(fileinfos, page)

	/* ページデータをHTMLにレンダリング */
	c.HTML(200, IndexHTML, gin.H{
		"fileinfos": imagefileinfos,
		"prevPage":  page - 1,
		"nextPage":  page + 1,
		"hasPrev":   page > 1,
		"hasNext":   end < len(fileinfos),
		"total":     len(fileinfos), // ファイル情報の総数
	})
}

/* ファイル検索画面のレスポンス */
func SearchEndPoint(c *gin.Context) {
	/* POSTパラメータからTagを取得 */
	AND_Tag := c.PostFormArray("AND")
	OR_Tag := c.PostFormArray("OR")
	NOT_Tag := c.PostFormArray("NOT")

	/* データベースからファイル情報を取得 */
	fileinfos = dbutills.GetFileInfo(AND_Tag, OR_Tag, NOT_Tag)

	/* POSTパラメータからフィルタパネルの状態を取得 */
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

	/* フィルタパネルの状態に応じてフィルタ */
	filter_fileinfos := FileInfoFilter(fileinfos, filterPanel)

	/* filterPanelをfilter_fileinfosで上書き */
	fileinfos = filter_fileinfos

	/* sort指定を取得 */
	sort := c.PostForm("sort")

	/* sort指定に応じてソート */
	fileinfos = SortFileInfo(fileinfos, sort)

	/* ページ番号を取得（デフォルトは1） */
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	/* ページデータ取得 */
	imagefileinfos, end := ExtractionPage(filter_fileinfos, page)

	/* パラメータから取得したものをlogに表示 */
	fmt.Println("AND_Tag:", AND_Tag)
	fmt.Println("OR_Tag:", OR_Tag)
	fmt.Println("NOT_Tag:", NOT_Tag)
	fmt.Println("filterPanel:", filterPanel)
	fmt.Println("page:", page)
	fmt.Println("sort:", sort)

	/* 表示されるファイルのタイプと名前を一行表示 */
	for _, info := range imagefileinfos {
		fmt.Println(info.TYPE, info.Name)
	}

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
