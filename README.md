# push2polygon

## 概要
競プロの問題を作問する際、github と rime.py を用いるのが一般的です。しかし、プライベートコンテストを作るには codeforces の polygon を使うしかなく、polygon で要求されるコードの形と rime.py で要求されるコードの形が異なるため、rime.py を用いることができません。
push2polygon はその問題を解決するためのものです。push2polygon を用いれば、github と rime.py を用いて作問し、polygon にアップロードすることができます。

## 中身
polygon の API の Wrapper(https://pkg.go.dev/github.com/variety-jones/polygon) を用いています。テストケース生成コードは、書き換える必要があるので、構文解析をし、polygon に合う形に変換してからアップロードしています。

# 使う人に向けて
## polygon 上でするべきこと

- create problem
- checker 作らないときは選択
- statement 記述
- commit

## 使い方

- template をコピーして、rime が使えるようにフォルダを構成する
- WA, TLE など、想定される判定をフォルダ名の suffix に書く
- polygon.txt に問題の id, スコアを記述
- gen_function.txt に、テストに使う関数名と、その関数で生成するテストケース数を記述
- polygon.exe を実行

## 注意点

- 2 回目以降の実行でエラーが出た場合、polygon 上で今までの変更を discard すること


## 後回しの改善点

- statement 日本語を push できるようにする。
日本語を push しようとすると署名エラーが起こる。

- 部分点をつけれるようにする