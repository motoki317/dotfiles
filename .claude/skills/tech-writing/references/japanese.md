# 日本語

## 語を選ぶ

単語を選ぶときは、日本のソフトウェア技術者が仕様書やPRに書く語にする。英語の術語を一語ずつ訳さない。

| 一語ずつ訳した語 | 技術者が書く語 |
|---|---|
| この直し | この修正、このパッチ |
| 前景で動かす | フォアグラウンドで実行する |
| 偽のSlackサーバー | Slackのモックサーバー |

## 表記

- 和文と英単語、数字、コードスパンのあいだに半角空白を入れない（「このREADMEは」）。
- 道具や概念に意志を持たせない。物が主語の「〜てくれる」を書かない（「このキャッシュが負荷を減らしてくれる」→「このキャッシュで負荷が減る」）。

## 点検

日本語の文章を投稿する前と保存する前に、毎回次の2つのコマンドを実行する。Slack、GitHub、文書のどれでも同じにする。

```sh
python3 -I $HOME/.claude/skills/tech-writing/scripts/yomiyasu_lint.py --doc-type <種類> --words $HOME/.claude/skills/tech-writing/words - < <ファイル>
perl -CSD -Mutf8 -ne 'if(/^\s*```/){$c=!$c;next} next if $c||/^    /; print "$.: $_" if /[ぁ-んァ-ヶ一-龥々ー、。」）] [A-Za-z0-9`]|[A-Za-z0-9`%] [ぁ-んァ-ヶ一-龥々ー（「]/' <ファイル>
```

- `<種類>`は、チャットなら`mail`、PRとissueなら`pr`、そのほかの文書なら`article`にする。
- 作業フォルダの外のファイルを引数に渡すとリンターは読まないので、標準入力で渡す。
- perlは、和文と英数字のあいだに半角空白がある行を出す。
- `excess_list`と`excess_bold`の指摘は使わない。

他人の文章を書き直したときは、最後に`mktemp -d`で一時フォルダを作り、元の文と書き直した文をそこにコピーする。そのフォルダで次のコマンドを実行し、「意味が変わる恐れ」の候補を一つずつ確かめる。SKILL.md「3. Keep the meaning of the source」に反する文は元に戻す。

```sh
python3 -I $HOME/.claude/skills/tech-writing/scripts/yomiyasu_diff.py --words $HOME/.claude/skills/tech-writing/words <元の文> <書き直した文>
```
