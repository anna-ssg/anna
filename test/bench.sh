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
echo "clone SSGs"

# Saaru & Sapling: fresh clones
git clone --depth=1 https://github.com/anirudhRowjee/saaru $BENCH_DIR/saaru
git clone --depth=1 https://github.com/NavinShrinivas/sapling $BENCH_DIR/sapling

# show commit hashes
get_commit_info() {
  local repo_dir=$1
  local name=$2
  local url=$3
  INFO=$(cd $repo_dir && git log --oneline -1)
  HASH=$(echo $INFO | cut -d' ' -f1)
  MSG=$(echo $INFO | cut -d' ' -f2-)
  FULL_HASH=$(cd $repo_dir && git log --format=%H -1)
  echo "$name: [\`$HASH\`]($url/commit/$FULL_HASH) $MSG"
}

echo "commit hashes:"
get_commit_info $REPO_ROOT "anna" "https://github.com/anna-ssg/anna" > $BENCH_DIR/commit_hashes.txt
get_commit_info $BENCH_DIR/saaru "saaru" "https://github.com/anirudhRowjee/saaru" >> $BENCH_DIR/commit_hashes.txt
get_commit_info $BENCH_DIR/sapling "sapling" "https://github.com/NavinShrinivas/sapling" >> $BENCH_DIR/commit_hashes.txt
cat $BENCH_DIR/commit_hashes.txt

echo "build SSGs"
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
