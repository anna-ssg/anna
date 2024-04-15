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
git clone https://github.com/acmpesuecc/anna /tmp/bench/anna
git clone https://github.com/anirudhRowjee/saaru /tmp/bench/saaru

# copy benchmark file
cp /tmp/bench/anna/site/content/posts/bench.md /tmp/bench/test.md 

# build saaru & anna
echo "build SSGs"
cd /tmp/bench/anna && go build && cd ..
cd saaru && cargo build --release && mv ./target/release/saaru . && cd ..

## setup hugo
hugo new site /tmp/bench/hugo; cd /tmp/bench/hugo
hugo new theme mytheme; echo "theme = 'mytheme'" >> hugo.toml; cd /tmp/bench

## setup 11ty
mkdir /tmp/bench/11ty -p

# clean content/* dirs
echo "Cleaning content directories"
rm -rf /tmp/bench/anna/site/content/posts/*
rm -rf /tmp/bench/saaru/docs/src/*
rm -rf /tmp/bench/hugo/content/*
rm -rf /tmp/bench/11ty/*

# create multiple copies of the test file
echo "Spawning $files different markdown files..."
for ((i = 0; i < files; i++)); do
    cp /tmp/bench/test.md "/tmp/bench/anna/site/content/posts/test_$i.md"
    cp /tmp/bench/test.md "/tmp/bench/saaru/docs/src/test_$i.md"
    cp /tmp/bench/test.md "/tmp/bench/hugo/content/test_$i.md"
    cp /tmp/bench/test.md "/tmp/bench/11ty/test_$i.md"
done

# run hyperfine
echo -e "\n"
echo "running benchmark: $files md files and $warm warmup runs"
echo -e "\n"
cd /tmp/bench/11ty && hyperfine -p 'sync' -w $warm "npx @11ty/eleventy" && cd ..
cd /tmp/bench/hugo && hyperfine -p 'sync' -w $warm "hugo" && cd .
cd /tmp/bench/saaru && hyperfine -p 'sync' -w $warm "./saaru --base-path ./docs" && cd ..
cd /tmp/bench/anna && hyperfine -p 'sync' -w $warm "./anna" && cd ..
echo -e "\n"
