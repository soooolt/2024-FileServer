/*******************************************************************************
* ファイル名: script.js
* 内容: 全ページ共通のスクリプト
* 作成日: 2024年3月15日
* 更新日: 2024年3月15日
*******************************************************************************/

/*--グローバル変数定義----------------------------------------------------------*/
/* フィルタパネル選択状態保存用 */
let filterPanel = {
    'FP_movie': false,
    'FP_animated': false,
    'FP_image': false,
    'FP_audio': false,
    'FP_manga': false,
    'FP_3D': false,
    'FP_2D': false,
    'FP_Real': false
};

/* タグの選択状態保存用 */
let tags = {
    'AND': [],
    'OR': [],
    'NOT': []
};

/* ボタンとフィルタパネルのマッピング */
const buttonToFilterMap = {
    'Movie': 'FP_movie',
    'Gif': 'FP_animated',
    'Image': 'FP_image',
    'Audio': 'FP_audio',
    'Manga': 'FP_manga',
    '3D': 'FP_3D',
    '2D': 'FP_2D',
    'Real': 'FP_Real'
};

/*--関数定義-------------------------------------------------------------------*/
/* /searchエンドポイントへfilterPanelとtagsをパラメータとしてPOST送信する */
function postFilter() {
    let query = '';
    for (const tagType in tags) {
        for (const tag of tags[tagType]) {
            if (tag) {
                query += `&${tagType}=${tag}`;
            }
        }
    }

    for (const filterName in filterPanel) {
        if (filterPanel[filterName]) {
            query += `&${filterName}=true`;
        }
    }

    // POSTリクエストを送信
    fetch('/search', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: query
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.text();
    })
    .then(html => {
        document.querySelector('#content').innerHTML = html;
    })
    .catch(error => console.error('Error:', error));
}

/* ローカルストレージからfilterPanelとtagsを復元する */
function restoreFilter() {
    /* タグの選択状態 */
    tags['AND'] = JSON.parse(localStorage.getItem('AND')) || [];
    tags['OR'] = JSON.parse(localStorage.getItem('OR')) || [];
    tags['NOT'] = JSON.parse(localStorage.getItem('NOT')) || [];

    /* 各検索フォームの入力値を更新 */
    document.querySelector('input[name="and_tag"]').value = tags['AND'].join(',');
    document.querySelector('input[name="or_tag"]').value = tags['OR'].join(',');
    document.querySelector('input[name="not_tag"]').value = tags['NOT'].join(',');

    /* ローカルストレージからフィルタパネルの状態を復元 */
    for (const filterName in filterPanel) {
        filterPanel[filterName] = localStorage.getItem(filterName) === 'true';
    }

    /* フィルタパネルの状態を更新 */
    document.querySelector('.filter-panel').querySelectorAll('.toggle').forEach((button) => {
        const buttonName = button.textContent.trim();
        const filterName = buttonToFilterMap[buttonName];
        updateButtonStyle(button, filterPanel[filterName]);
    });
}

/* ボタンのスタイルを変更 */
function updateButtonStyle(button, state) {
    if (state) {
        button.classList.remove('btn-outline-secondary');
        button.classList.add('btn-secondary');
        button.setAttribute('aria-pressed', 'true');
    } else {
        button.classList.remove('btn-secondary');
        button.classList.add('btn-outline-secondary');
        button.setAttribute('aria-pressed', 'false');
    }
}

/*--イベントリスナー------------------------------------------------------------*/
/* ページ読み込み後 */
document.addEventListener('DOMContentLoaded', (event) => {
    /* ローカルストレージから状態を復元 */
    restoreFilter();

    /* フィルタパネル*/
    document.querySelector('.filter-panel').querySelectorAll('.toggle').forEach((button) => {
        /* ボタン名とフィルタ名の取得 */
        const buttonName = button.textContent.trim();
        const filterName = buttonToFilterMap[buttonName];

        /*--クリック時のイベント-- */
        button.addEventListener('click', () => {
            /* 変数更新 */
            filterPanel[filterName] = !filterPanel[filterName];

            /* ローカルストレージに状態を保存 */
            localStorage.setItem(filterName, filterPanel[filterName]);

            /* ボタンのスタイルを更新 */
            updateButtonStyle(button, filterPanel[filterName]);
        });
    });

    /* 検索ボタン */
    document.querySelector('.search-trigger').addEventListener('click', (event) => {
        /* フォームの送信をキャンセル */
        event.preventDefault();

        /* 検索フォームの内容を取得 */
        tags['AND'] = document.querySelector('input[name="and_tag"]').value.split(',');
        tags['OR'] = document.querySelector('input[name="or_tag"]').value.split(',');
        tags['NOT'] = document.querySelector('input[name="not_tag"]').value.split(',');

        /* 検索フォームの内容をローカルストレージに保存 */
        localStorage.setItem('AND', JSON.stringify(tags['AND']));
        localStorage.setItem('OR', JSON.stringify(tags['OR']));
        localStorage.setItem('NOT', JSON.stringify(tags['NOT']));

        /* フィルタを送信 */
        postFilter();
    });
});