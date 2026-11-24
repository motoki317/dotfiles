#!/usr/bin/env -S uv run --script
# /// script
# requires-python = ">=3.12"
# dependencies = ["sudachipy==0.7.0", "sudachidict-core==20261015"]
# ///
"""List words that your recent posts use but native engineers' writing rarely uses.

Some of these words are calques: a word-by-word translation of an English term (深さ for depth).
Statistics cannot tell a calque from a rare but natural word, so a person judges each candidate
and adds the calques to words/flag.txt or words/flag-patterns.txt.

--sqlite and --sql: the query returns one row per post, with its Slack mrkdwn text in the first column.
--native: "lemma<TAB>count" lines of native text, made with the Sudachi dictionary pinned above.
  A new dictionary version changes lemmas, so recount the native text when you change the pin.
--reviewed: one lemma per line. The script skips these lemmas and the words in --words.
Each output line: posts that use the lemma, the lemma, its native count, and the shortest sentence that uses it.
"""
import argparse
import collections
import re
import sqlite3

from sudachipy import Dictionary

CONTENT = {"名詞", "動詞", "形容詞", "形状詞"}
LATIN = re.compile(r"^[A-Za-z0-9_.\-/#@:+]+$")
KATAKANA = re.compile(r"^[ァ-ヶー・]+$")
JAPANESE = re.compile(r"[ぁ-んァ-ヶ一-龥々]")
# A sentence about wording quotes the word or names the topic: it discusses a calque instead of using one.
ABOUT_WORDS = re.compile(r"「|直訳|変な日本語|不思議な日本語|tech-writing|lint")


def sentences(mrkdwn):
    text = re.sub(r"```.*?```", "\n", mrkdwn, flags=re.S)
    text = re.sub(r"<(?:https?://|mailto:|#|@|!)[^>|]*(?:\|([^>]*))?>", lambda m: m.group(1) or " ", text)
    text = re.sub(r"^&gt;.*$", "", text, flags=re.M)  # a quoted line is someone else's text
    text = text.replace("&amp;", "&").replace("&lt;", "<").replace("&gt;", ">")
    text = re.sub(r"`[^`\n]*`", " CODE ", text)
    text = re.sub(r"https?://\S+", " ", text)
    text = re.sub(r":[a-z0-9_+\-]+:", " ", text)
    text = re.sub(r"[*_~]", "", text)
    for line in re.split(r"(?<=[。！？!?])|\n", text):
        s = re.sub(r"^\s*(?:[-*•◦▪︎・]|\d+\.|#+|>)\s*", "", line).strip()
        if len(JAPANESE.findall(s)) >= 6:
            yield s


def lemmas(tok, sentence):
    """Lemmas of content words, without proper nouns, numbers, nouns like こと, Latin words, and katakana loanwords."""
    for m in tok.tokenize(sentence):
        pos = m.part_of_speech()
        content = (pos[0] in CONTENT and pos[1] != "数詞" and pos[2] != "助数詞可能") or pos[:2] == ("接尾辞", "名詞的")
        if not content or pos[1] in ("固有名詞", "非自立可能") or LATIN.match(m.surface()):
            continue
        lemma = m.normalized_form()
        if not KATAKANA.match(lemma):
            yield lemma


def known(paths):
    """The first column of each list. A flag-patterns.txt name ends in the words to use in （）, which the split drops."""
    out = set()
    for path in paths:
        with open(path) as f:
            for line in f:
                if line.strip() and not line.startswith("#"):
                    out.add(line.split("\t")[0].split("（")[0].strip())
    return out


def main():
    p = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    p.add_argument("--sqlite", required=True)
    p.add_argument("--sql", required=True)
    p.add_argument("--native", required=True)
    p.add_argument("--words", required=True, help="the tech-writing words/ folder")
    p.add_argument("--reviewed", required=True)
    p.add_argument("--top", type=int, default=40)
    p.add_argument("--min-posts", type=int, default=2)
    p.add_argument("--rare-per-million", type=float, default=5.0)
    a = p.parse_args()

    native = collections.Counter()
    with open(a.native) as f:
        for line in f:
            lemma, n = line.rstrip("\n").split("\t")
            native[lemma] = int(n)
    native_total = sum(native.values())
    skip = known([f"{a.words}/{name}" for name in ("flag.txt", "flag-patterns.txt", "allow.txt")] + [a.reviewed])

    tok = Dictionary(dict="core").tokenizer(mode="C")
    posts = [row[0] for row in sqlite3.connect(f"file:{a.sqlite}?mode=ro", uri=True).execute(a.sql)]
    posts_using = collections.defaultdict(set)  # lemma -> indexes of the posts that use it
    example = {}
    for i, post in enumerate(posts):
        for s in sentences(post):
            if ABOUT_WORDS.search(s):
                continue
            for lemma in lemmas(tok, s):
                posts_using[lemma].add(i)
                if lemma not in example or len(s) < len(example[lemma]):
                    example[lemma] = s

    rows = []
    for lemma, ps in posts_using.items():
        n = native[lemma]
        if len(ps) >= a.min_posts and lemma not in skip and n * 1e6 / native_total <= a.rare_per_million:
            rows.append((-len(ps), n, lemma))
    rows.sort()
    print(f"# {len(posts)} posts, {len(rows)} candidates")
    for minus_posts, n, lemma in rows[:a.top]:
        print(f"{-minus_posts}\t{lemma}\tnative={n} ({n * 1e6 / native_total:.1f}/M)\t{example[lemma][:100]}")


if __name__ == "__main__":
    main()
