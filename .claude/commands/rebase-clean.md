---
description: mainにsoft resetして試行錯誤が見えないよう綺麗にコミットを整理し直す
allowed-tools: [Bash, Read, AskUserQuestion]
---

# Rebase Clean

ブランチのコミットをsoft resetし、レビュワーが読みやすいように論理単位ごとに綺麗なコミットに整理し直す。

## 前提条件

- 現在のブランチがmain以外であること
- mainからの差分コミットが存在すること

## 手順

### Step 1: merge-baseの特定と現状把握

**重要**: `main` ではなく `merge-base` を使う。mainが進んでいる場合、`git reset --soft main` はmainの新しい変更も巻き込んでしまう。

```bash
MERGE_BASE=$(git merge-base main HEAD)
echo "Merge base: $MERGE_BASE"
git log --oneline main..HEAD           # mainとの差分（mainが進んだ分も含む）
git log --oneline $MERGE_BASE..HEAD    # ブランチ固有のコミットのみ
git log --oneline $MERGE_BASE..main    # mainが進んだコミット（あれば表示）
```

全コミットの内容を確認し、以下を分析する:

- 各コミットの変更内容
- fixコミットやtry-and-errorの痕跡
- 論理的にまとめるべき単位

### Step 2: コミット整理計画

以下の原則でコミットをグループ化する:

1. **1コミット = 1つの論理的変更単位**（proto、ドメイン層、インフラ層、UI等）
2. **fixコミットは親コミットに吸収**（試行錯誤を隠す）
3. **レイヤー順にコミット**（依存関係の下流から上流へ）
4. **生成コードは元となるコミットに含める**

計画をユーザーに提示し、承認を得る。

### Step 3: Soft Reset（merge-baseへ）

```bash
git reset --soft $MERGE_BASE
git reset HEAD
```

**注意**: `main` ではなく `$MERGE_BASE` へリセットすること。これによりブランチ固有の変更のみが対象になる。

### Step 4: 計画に従ってコミット

ファイルを論理単位ごとに`git add`し、コミットメッセージはConventional Commits形式で作成。

コミットメッセージの方針:

- 既存コミットメッセージを参考にしつつ、整理後の内容に合わせて書き直す
- `Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>` を末尾に付与

### Step 5: 検証（整理後）

```bash
git log --oneline $MERGE_BASE..HEAD
git status  # nothing to commit, working tree clean
```

全変更が漏れなくコミットされていることを確認。

### Step 6: 最新mainへのrebase

mainが進んでいた場合のみ実行する（Step 1で `$MERGE_BASE..main` にコミットがあった場合）。

```bash
git fetch origin main
git rebase main
```

コンフリクトが発生した場合はユーザーに報告し、解決方針を相談する。

### Step 7: 最終検証・push

```bash
git log --oneline main..HEAD
git status
git push --force-with-lease
```
