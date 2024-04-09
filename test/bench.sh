#!/bin/bash
set -e
clear

# parameters
echo "bash: number of md files to create?"
read files
echo "bash: how many warm-up runs?"
read warm

# echo "bash: cleanup"
cd /tmp || exit
rm -rf bench
mkdir bench; cd bench

# check if hyperfine is installed
if ! command -v hyperfine &>/dev/null; then
    echo "bash!: hyperfine is not installed. Please install hyperfine to continue."
    exit 1
fi

# cloning candidates
echo "git: cloning ssgs"
git clone --depth=1 https://github.com/acmpesuecc/anna
git clone --depth=1 https://github.com/anirudhRowjee/saaru

# benchmark file
cp anna/site/content/posts/bench.md test.md

# build SSGs
cd anna; go build; cd ..
cd saaru; cargo build --release; mv ./target/release/saaru .; cd ..
## pnpx @11ty/eleventy

# fixes
## saaru fixes (jinja added in the anna bench.md)
## anna fixes (none)
## 11ty fixes (yet to bench)

# clean content dirs (no md files other than test.md)
rm -rf anna/site/content/posts/*
rm -rf saaru/docs/src/*

# create multiple copies of the test file
echo "bash: spawning $files different markdown files..."
for ((i = 0; i < files; i++)); do
    cp test.md "anna/site/content/posts/test_$i.md"
    cp test.md "saaru/docs/src/test_$i.md"
done

# run the benchmark
echo "hyperfine: running benchmark with hyperfine..."
cd anna; hyperfine -p 'sync' -w $warm "./anna"; cd ..
cd saaru; hyperfine -p 'sync' -w $warm "./saaru --base-path ./docs"; cd ..

# cleanup
echo "bash: cleaning up"
rm -rf bench
