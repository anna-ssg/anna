#!/bin/bash
files=1000
warm=10
BASE_DIR=$(pwd)

# Determine repo root
if [ -d "$BASE_DIR/site" ]; then
  REPO_ROOT=$BASE_DIR
else
  REPO_ROOT=$(cd $BASE_DIR/.. && pwd)
fi

# Create bench directory
mkdir -p $BASE_DIR/tmp/bench

# deps
if ! command -v hyperfine &>/dev/null; then
    echo "hyperfine is not installed. Please install hyperfine to continue."
    exit 1
fi
if ! command -v hugo &>/dev/null; then
    echo "hugo is not installed. Please install hugo to continue."
fi

# cloning candidates
echo ""
echo "clone SSGs"
echo ""
clone_or_pull() {
  local repo=$1
  local dir=$2
  if [ "$repo" = "https://github.com/anna-ssg/anna" ]; then
    # For Anna, use current source code
    rm -rf "$dir"
    mkdir -p "$dir"
    cp -r "$REPO_ROOT/site" "$dir/"
    cd "$REPO_ROOT" && go build -o "$dir/anna"
    cd "$REPO_ROOT" && GOEXPERIMENT=greenteagc go build -o "$dir/anna_greentea"
  else
    # For others, always fresh clone
    rm -rf "$dir"
    git clone --depth=1 "$repo" "$dir"
  fi
}

clone_or_pull https://github.com/anna-ssg/anna $BASE_DIR/tmp/bench/anna
clone_or_pull https://github.com/anirudhRowjee/saaru $BASE_DIR/tmp/bench/saaru
clone_or_pull https://github.com/NavinShrinivas/sapling $BASE_DIR/tmp/bench/sapling

# show commit hashes
echo ""
echo "Commit hashes:"
echo "anna: $(cd $REPO_ROOT && git rev-parse HEAD)"
echo "saaru: $(cd $BASE_DIR/tmp/bench/saaru && git rev-parse HEAD)"
echo "sapling: $(cd $BASE_DIR/tmp/bench/sapling && git rev-parse HEAD)"
echo ""

cp $REPO_ROOT/site/content/posts/bench.md $BASE_DIR/tmp/bench/test.md

echo ""
echo "build SSGs"
echo ""
cd $BASE_DIR/tmp/bench/sapling && cargo build --release && mv target/release/sapling .
cd $BASE_DIR/tmp/bench/saaru && cargo build --release && mv target/release/saaru .

## setup hugo
rm -rf $BASE_DIR/tmp/bench/hugo
hugo new site $BASE_DIR/tmp/bench/hugo
cd $BASE_DIR/tmp/bench/hugo
hugo new theme mytheme
echo "theme = 'mytheme'" >> hugo.toml
cd ../..

# clean content/* dirs
rm -rf $BASE_DIR/tmp/bench/anna/site/content/posts/*
rm -rf $BASE_DIR/tmp/bench/saaru/docs/src/*
rm -rf $BASE_DIR/tmp/bench/sapling/benchmark/content/blog/*
rm -rf $BASE_DIR/tmp/bench/hugo/content/*

echo ""
echo "Spawning $files different markdown files..."
for ((i = 0; i < files; i++)); do
    cp $BASE_DIR/tmp/bench/test.md "$BASE_DIR/tmp/bench/anna/site/content/posts/test_$i.md"
    cp $BASE_DIR/tmp/bench/test.md "$BASE_DIR/tmp/bench/saaru/docs/src/test_$i.md"
    cp $BASE_DIR/tmp/bench/test.md "$BASE_DIR/tmp/bench/sapling/benchmark/content/blogs/test_$i.md"
    cp $BASE_DIR/tmp/bench/test.md "$BASE_DIR/tmp/bench/hugo/content/test_$i.md"
done
echo ""

echo "Begin Benchmark"
echo "running benchmark: $files md files and $warm warmup runs \n"
hyperfine -p 'sync' -w $warm \
  "cd $BASE_DIR/tmp/bench/saaru && ./saaru --base-path ./docs" \
  "cd $BASE_DIR/tmp/bench/sapling/benchmark && ./../sapling run" \
  "cd $BASE_DIR/tmp/bench/anna && ./anna_greentea" \
  "cd $BASE_DIR/tmp/bench/anna && ./anna" \
  "cd $BASE_DIR/tmp/bench/hugo && hugo"
echo "End Benchmark"
