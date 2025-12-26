#!/bin/bash
files=1000
warm=10
REPO_ROOT=$(pwd)
BENCH_DIR=$REPO_ROOT/tmp/bench
rm -rf $BENCH_DIR

# build anna
mkdir -p $BENCH_DIR/anna
cp -r site/ $BENCH_DIR/anna/site
go build -o $BENCH_DIR/anna/anna
GOEXPERIMENT=greenteagc go build -o $BENCH_DIR/anna/anna_greentea

# deps
if ! command -v hyperfine &>/dev/null; then
    echo "hyperfine is not installed. Please install hyperfine to continue." && exit 1
fi
if ! command -v hugo &>/dev/null; then
    echo "hugo is not installed. Please install hugo to continue."
fi

# cloning candidates
echo "\n clone SSGs"

# Saaru & Sapling: fresh clones
git clone --depth=1 https://github.com/anirudhRowjee/saaru $BENCH_DIR/saaru
git clone --depth=1 https://github.com/NavinShrinivas/sapling $BENCH_DIR/sapling

# show commit hashes
echo "\n Commit hashes:"
ANNA_HASH=$(cd $REPO_ROOT && git log --format=%h -1)
echo "anna: [\`$ANNA_HASH\`](https://github.com/anna-ssg/anna/commit/$(cd $REPO_ROOT && git log --format=%H -1))" > $BENCH_DIR/commit_hashes.txt

SAARU_HASH=$(cd $BENCH_DIR/saaru && git log --format=%h -1)
echo "saaru: [\`$SAARU_HASH\`](https://github.com/anirudhRowjee/saaru/commit/$(cd $BENCH_DIR/saaru && git log --format=%H -1))" >> $BENCH_DIR/commit_hashes.txt

SAPLING_HASH=$(cd $BENCH_DIR/sapling && git log --format=%h -1)
echo "sapling: [\`$SAPLING_HASH\`](https://github.com/NavinShrinivas/sapling/commit/$(cd $BENCH_DIR/sapling && git log --format=%H -1))" >> $BENCH_DIR/commit_hashes.txt

cat $BENCH_DIR/commit_hashes.txt

echo "\n build SSGs"
cargo build --release --manifest-path $BENCH_DIR/sapling/Cargo.toml
cargo build --release --manifest-path $BENCH_DIR/saaru/Cargo.toml

## setup hugo
hugo new site $BENCH_DIR/hugo; cd $BENCH_DIR/hugo; hugo new theme mytheme
echo "theme = 'mytheme'" >> hugo.toml

# clean content/* dirs
rm -rf $BENCH_DIR/anna/site/content && mkdir -p $BENCH_DIR/anna/site/content/posts
rm -rf $BENCH_DIR/saaru/docs/src/*
rm -rf $BENCH_DIR/sapling/benchmark/content/blogs/*
rm -rf $BENCH_DIR/hugo/content/*

# populate with test files
cp $REPO_ROOT/site/content/posts/bench.md $BENCH_DIR/test.md
for ((i = 0; i < files; i++)); do
    cp $BENCH_DIR/test.md "$BENCH_DIR/anna/site/content/posts/test_$i.md"
    cp $BENCH_DIR/test.md "$BENCH_DIR/saaru/docs/src/test_$i.md"
    cp $BENCH_DIR/test.md "$BENCH_DIR/sapling/benchmark/content/blogs/test_$i.md"
    cp $BENCH_DIR/test.md "$BENCH_DIR/hugo/content/test_$i.md"
done

echo "running benchmark: $files md files and $warm warmup runs" > $BENCH_DIR/bench_results.txt
hyperfine -p 'sync' -w $warm \
  "cd $BENCH_DIR/saaru && ./target/release/saaru --base-path ./docs" \
  "cd $BENCH_DIR/sapling/benchmark && ./../target/release/sapling run" \
  "cd $BENCH_DIR/anna && ./anna_greentea" \
  "cd $BENCH_DIR/anna && ./anna" \
  "cd $BENCH_DIR/hugo && hugo" >> $BENCH_DIR/bench_results.txt
