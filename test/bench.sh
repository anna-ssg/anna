#!/bin/bash
set -e
clear

# parameters
files=1000
warm=10

# cleanup
cleanup() {
    echo "Cleaning up"
    rm -rf /tmp/bench
}
trap cleanup EXIT

# echo "bash: cleanup"
cd /tmp || exit
rm -rf bench
mkdir bench
cd bench

# check if hyperfine is installed
if ! command -v hyperfine &>/dev/null; then
    echo "hyperfine is not installed. Please install hyperfine to continue."
    exit 1
fi

# cloning candidates
echo "Cloning SSGs"
git clone https://github.com/acmpesuecc/anna
git clone https://github.com/anirudhRowjee/saaru
sleep 1; clear

# benchmark file
cp anna/site/content/posts/bench.md test.md 

# checkout at v1 and commit
cd anna && git checkout v1.0.0-alpha; cd ..
cd saaru && git checkout c17930724fcaad67e1cfa3cc969667d043c1b826; cd ..
sleep 1; clear

# build SSGs
cd anna; go build; cd ..
cd saaru; cargo build --release; mv ./target/release/saaru .; cd ..
sleep 1; clear

# clean content dirs (no md files other than test.md)
echo "Cleaning content directories"
rm -rf anna/site/content/posts/*
rm -rf saaru/docs/src/*

# create multiple copies of the test file
echo "Spawning $files different markdown files..."
for ((i = 0; i < files; i++)); do
    cp test.md "anna/site/content/posts/test_$i.md"
    cp test.md "saaru/docs/src/test_$i.md"
done
sleep 1; clear

# run the benchmark
echo "Running benchmark with $files md files and $warm warmup runs"
cd anna; hyperfine -p 'sync' -w $warm "./anna"; cd ..
cd saaru; hyperfine -p 'sync' -w $warm "./saaru --base-path ./docs"; cd ..

echo "Benchmarking finished."
