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

# cloning candidates
echo "clone SSGs"
git clone https://github.com/acmpesuecc/anna /tmp/bench/anna
git clone https://github.com/anirudhRowjee/saaru /tmp/bench/saaru
sleep 1

# copy benchmark file
cp /tmp/bench/anna/site/content/posts/bench.md /tmp/bench/test.md 

# checkout specific versions
cd /tmp/bench/anna && git checkout v1.0.0-alpha && cd ..
cd saaru && git checkout c17930724fcaad67e1cfa3cc969667d043c1b826 && cd ..
sleep 1

# build SSGs
echo "build SSGs"
cd /tmp/bench/anna && go build && cd ..
cd saaru && cargo build --release && mv ./target/release/saaru . && cd ..
echo "finished building"

# clean content/* dirs
echo "Cleaning content directories"
rm -rf /tmp/bench/anna/site/content/posts/*
rm -rf /tmp/bench/saaru/docs/src/*

# create multiple copies of the test file
echo "Spawning $files different markdown files..."
for ((i = 0; i < files; i++)); do
    cp /tmp/bench/test.md "/tmp/bench/anna/site/content/posts/test_$i.md"
    cp /tmp/bench/test.md "/tmp/bench/saaru/docs/src/test_$i.md"
done
sleep 1

# cooldown
echo "cooldown for 30s"
sleep 30

# run hyperfine
echo "running benchmark: $files md files and $warm warmup runs"
cd /tmp/bench/anna && hyperfine -p 'sync' -w $warm "./anna" && cd ..
cd saaru && hyperfine -p 'sync' -w $warm "./saaru --base-path ./docs" && cd ..
