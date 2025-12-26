#!/bin/bash
# parameters
files=1000
warm=10

BASE_DIR=$(pwd)

# cleanup
cleanup() {
    echo "cleaning up"
    rm -rf $BASE_DIR/tmp/bench
}
trap cleanup EXIT

# check if hyperfine is installed
if ! command -v hyperfine &>/dev/null; then
    echo "hyperfine is not installed. Please install hyperfine to continue."
    exit 1
fi

# check if hugo is installed
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
  if [ ! -d "$dir" ]; then
    git clone --depth=1 "$repo" "$dir"
  else
    cd "$dir" && git pull --depth=1 --ff-only
    if [ $? -ne 0 ]; then
      cd "$(dirname "$dir")"
      rm -rf "$dir"
      git clone --depth=1 "$repo" "$dir"
    fi
  fi
}

clone_or_pull https://github.com/anna-ssg/anna $BASE_DIR/tmp/bench/anna
clone_or_pull https://github.com/anirudhRowjee/saaru $BASE_DIR/tmp/bench/saaru
clone_or_pull https://github.com/NavinShrinivas/sapling $BASE_DIR/tmp/bench/sapling

# copy benchmark file
cp $BASE_DIR/tmp/bench/anna/site/content/posts/bench.md $BASE_DIR/tmp/bench/test.md

echo ""
echo "build SSGs"
echo ""
cd $BASE_DIR/tmp/bench/anna && go build && cd ../..
cd $BASE_DIR/tmp/bench/anna && GOEXPERIMENT=greenteagc go build -o anna_greentea && cd ../..

# build rust based SSGs (edit this block if they are already installed)
cd $BASE_DIR/tmp/bench/sapling && cargo build --release && mv target/release/sapling .
cd $BASE_DIR/tmp/bench/saaru && cargo build --release && mv target/release/saaru .

## setup hugo
hugo new site $BASE_DIR/tmp/bench/hugo; cd $BASE_DIR/tmp/bench/hugo
hugo new theme mytheme; echo "theme = 'mytheme'" >> hugo.toml; cd ../..

## setup 11ty

# clean content/* dirs
echo ""
echo "Cleaning content directories"
echo ""
rm -rf $BASE_DIR/tmp/bench/anna/site/content/posts/*
rm -rf $BASE_DIR/tmp/bench/saaru/docs/src/*
rm -rf $BASE_DIR/tmp/bench/sapling/benchmark/content/blog/*
rm -rf $BASE_DIR/tmp/bench/hugo/content/*

# create multiple copies of the test file
echo ""
echo "Spawning $files different markdown files..."
for ((i = 0; i < files; i++)); do
    cp $BASE_DIR/tmp/bench/test.md "$BASE_DIR/tmp/bench/anna/site/content/posts/test_$i.md"
    cp $BASE_DIR/tmp/bench/test.md "$BASE_DIR/tmp/bench/saaru/docs/src/test_$i.md"
    cp $BASE_DIR/tmp/bench/test.md "$BASE_DIR/tmp/bench/sapling/benchmark/content/blogs/test_$i.md"
    cp $BASE_DIR/tmp/bench/test.md "$BASE_DIR/tmp/bench/hugo/content/test_$i.md"
done
echo ""

# begin benchmark
echo "Begin Benchmark"

# run hyperfine
echo "running benchmark: $files md files and $warm warmup runs"
echo ""
hyperfine -p 'sync' -w $warm \
  "cd $BASE_DIR/tmp/bench/hugo && hugo" \
  "cd $BASE_DIR/tmp/bench/anna && ./anna" \
  "cd $BASE_DIR/tmp/bench/anna && ./anna_greentea" \
  "cd $BASE_DIR/tmp/bench/saaru && ./saaru --base-path ./docs" \
  "cd $BASE_DIR/tmp/bench/sapling/benchmark && ./../sapling run"
echo ""
echo "End Benchmark"