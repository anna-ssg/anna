#!/bin/bash
# parameters
files=1000
warm=10

# cleanup
cleanup() {
    echo "cleaning up"
    rm -rf /tmp/bench
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
echo "clone SSGs"
git clone --depth=1 https://github.com/acmpesuecc/anna /tmp/bench/anna 
git clone --depth=1 https://github.com/anirudhRowjee/saaru /tmp/bench/saaru
git clone --depth=1 https://github.com/NavinShrinivas/sapling /tmp/bench/sapling

# copy benchmark file
cp /tmp/bench/anna/site/content/posts/bench.md /tmp/bench/test.md 

echo "build SSGs"
cd /tmp/bench/anna && go build && cd /tmp/bench

# build rust based SSGs (comment this block if you already have these installed)
# sapling
cd /tmp/bench/sapling && cargo build --release
sudo rm -rf /usr/local/src/sapling
sudo mkdir -p /usr/local/src/sapling
sudo cp -r /tmp/bench/sapling/* /usr/local/src/sapling/
cd /usr/local/src/sapling
sudo mkdir -p /usr/local/bin/
sudo rm -rf /usr/local/bin/sapling
sudo ln -sf /usr/local/src/sapling/target/release/sapling /usr/local/bin/sapling 
sudo chmod +x /usr/local/bin/sapling

# saaru
cd /tmp/bench/saaru && cargo build --release
sudo rm -rf /usr/local/src/saaru
sudo mkdir -p /usr/local/src/saaru
sudo cp -r /tmp/bench/saaru/* /usr/local/src/saaru/
cd /usr/local/src/saaru
sudo mkdir -p /usr/local/bin/
sudo rm -rf /usr/local/bin/saaru
sudo ln -sf /usr/local/src/saaru/target/release/saaru /usr/local/bin/saaru
sudo chmod +x /usr/local/bin/saaru

## setup hugo
hugo new site /tmp/bench/hugo; cd /tmp/bench/hugo
hugo new theme mytheme; echo "theme = 'mytheme'" >> hugo.toml; cd /tmp/bench

## setup 11ty
mkdir /tmp/bench/11ty -p

# clean content/* dirs
echo "Cleaning content directories"
rm -rf /tmp/bench/anna/site/content/posts/*
rm -rf /tmp/bench/saaru/docs/src/*
rm -rf /tmp/bench/sapling/benchmark/content/blog/*
rm -rf /tmp/bench/hugo/content/*
rm -rf /tmp/bench/11ty/*

# create multiple copies of the test file
echo "Spawning $files different markdown files..."
for ((i = 0; i < files; i++)); do
    cp /tmp/bench/test.md "/tmp/bench/anna/site/content/posts/test_$i.md"
    cp /tmp/bench/test.md "/tmp/bench/saaru/docs/src/test_$i.md"
    cp /tmp/bench/test.md "/tmp/bench/sapling/benchmark/content/blogs/test_$i.md"
    cp /tmp/bench/test.md "/tmp/bench/hugo/content/test_$i.md"
    cp /tmp/bench/test.md "/tmp/bench/11ty/test_$i.md"
done

# run hyperfine
echo -e "\n"
echo "running benchmark: $files md files and $warm warmup runs"
echo -e "\n"
hyperfine -p 'sync' -w $warm \
  "cd /tmp/bench/11ty && npx @11ty/eleventy" \
  "cd /tmp/bench/hugo && hugo" \
  "cd /tmp/bench/anna && ./anna"
  "cd /tmp/bench/saaru && saaru --base-path ./docs" \
  "cd /tmp/bench/sapling/benchmark && sapling run" \
echo -e "\n"

