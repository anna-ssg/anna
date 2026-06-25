#!/bin/bash
set -e

files=1000
warm=10

REPO_ROOT=$(pwd)
BENCH_DIR="$REPO_ROOT/tmp/bench"

rm -rf "$BENCH_DIR"

# build local anna
mkdir -p "$BENCH_DIR/anna"
cp -r site/ "$BENCH_DIR/anna/site"

go build -o "$BENCH_DIR/anna/anna"

# deps
if ! command -v hyperfine &>/dev/null; then
    echo "hyperfine is not installed. Please install hyperfine to continue."
    exit 1
fi

if ! command -v hugo &>/dev/null; then
    echo "hugo is not installed. Please install hugo to continue."
    exit 1
fi

# cloning candidates
echo "clone SSGs"

git clone --depth=1 https://github.com/anirudhRowjee/saaru "$BENCH_DIR/saaru"
git clone --depth=1 https://github.com/NavinShrinivas/sapling "$BENCH_DIR/sapling"
git clone --depth=1 https://github.com/anna-ssg/anna "$BENCH_DIR/anna-main"

# download latest anna release
echo "download latest anna release"

REPO="anna-ssg/anna"
ARCH="Linux_x86_64"

LATEST_TAG=$(
    curl -s https://api.github.com/repos/$REPO/releases/latest |
    sed -n 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/p'
)

if [[ -z "$LATEST_TAG" ]]; then
    echo "Failed to determine latest release tag"
    exit 1
fi

mkdir -p "$BENCH_DIR/anna-release"

TARBALL="anna_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$TARBALL"

curl -Ls "$URL" | tar -xz -C "$BENCH_DIR/anna-release"

if [[ ! -f "$BENCH_DIR/anna-release/anna" ]]; then
    echo "anna binary not found after extraction"
    exit 1
fi

chmod +x "$BENCH_DIR/anna-release/anna"

RELEASE_SHA=$(
    git ls-remote --tags https://github.com/anna-ssg/anna.git \
    "refs/tags/$LATEST_TAG^{}" | awk '{print $1}'
)

if [[ -z "$RELEASE_SHA" ]]; then
    RELEASE_SHA=$(
        git ls-remote --tags https://github.com/anna-ssg/anna.git \
        "refs/tags/$LATEST_TAG" | awk '{print $1}'
    )
fi

# show commit hashes
get_commit_info() {
    local repo_dir=$1
    local name=$2
    local url=$3

    INFO=$(cd "$repo_dir" && git log --oneline -1)
    HASH=$(echo "$INFO" | cut -d' ' -f1)
    MSG=$(echo "$INFO" | cut -d' ' -f2-)
    FULL_HASH=$(cd "$repo_dir" && git log --format=%H -1)

    echo "$name: [\`$HASH\`]($url/commit/$FULL_HASH) $MSG"
}

echo "commit hashes:"

get_commit_info "$REPO_ROOT" "anna" "https://github.com/anna-ssg/anna" > "$BENCH_DIR/commit_hashes.txt"
get_commit_info "$BENCH_DIR/saaru" "saaru" "https://github.com/anirudhRowjee/saaru" >> "$BENCH_DIR/commit_hashes.txt"
get_commit_info "$BENCH_DIR/sapling" "sapling" "https://github.com/NavinShrinivas/sapling" >> "$BENCH_DIR/commit_hashes.txt"
get_commit_info "$BENCH_DIR/anna-main" "anna-main" "https://github.com/anna-ssg/anna" >> "$BENCH_DIR/commit_hashes.txt"

echo "anna-release ($LATEST_TAG): [\`${RELEASE_SHA:0:7}\`](https://github.com/anna-ssg/anna/commit/$RELEASE_SHA)" >> "$BENCH_DIR/commit_hashes.txt"

cat "$BENCH_DIR/commit_hashes.txt"

echo "build SSGs"

cargo build --release --manifest-path "$BENCH_DIR/sapling/Cargo.toml"
cargo build --release --manifest-path "$BENCH_DIR/saaru/Cargo.toml"

(
    cd "$BENCH_DIR/anna-main"
    go build -o anna
)

# setup hugo
hugo new site "$BENCH_DIR/hugo"
cd "$BENCH_DIR/hugo"

hugo new theme mytheme
echo "theme = 'mytheme'" >> hugo.toml

cd "$REPO_ROOT"

# clean content dirs
rm -rf "$BENCH_DIR/anna/site/content"
mkdir -p "$BENCH_DIR/anna/site/content/posts"

rm -rf "$BENCH_DIR/anna-main/site/content"
mkdir -p "$BENCH_DIR/anna-main/site/content/posts"

# anna-release needs a full site tree, not just content/
cp -r site/ "$BENCH_DIR/anna-release/site"

rm -rf "$BENCH_DIR/anna-release/site/content"
mkdir -p "$BENCH_DIR/anna-release/site/content/posts"

rm -rf "$BENCH_DIR/saaru/docs/src/"*
rm -rf "$BENCH_DIR/sapling/benchmark/content/blogs/"*
rm -rf "$BENCH_DIR/hugo/content/"*

# populate with test files
cp "$REPO_ROOT/site/content/posts/bench.md" "$BENCH_DIR/test.md"

for ((i=0; i<files; i++)); do
    cp "$BENCH_DIR/test.md" "$BENCH_DIR/anna/site/content/posts/test_$i.md"
    cp "$BENCH_DIR/test.md" "$BENCH_DIR/anna-main/site/content/posts/test_$i.md"
    cp "$BENCH_DIR/test.md" "$BENCH_DIR/anna-release/site/content/posts/test_$i.md"
    cp "$BENCH_DIR/test.md" "$BENCH_DIR/saaru/docs/src/test_$i.md"
    cp "$BENCH_DIR/test.md" "$BENCH_DIR/sapling/benchmark/content/blogs/test_$i.md"
    cp "$BENCH_DIR/test.md" "$BENCH_DIR/hugo/content/test_$i.md"
done

echo "running benchmark: $files md files and $warm warmup runs" | tee "$BENCH_DIR/bench_results.txt"

hyperfine -p 'sync' -w "$warm" \
  -n saaru        "cd $BENCH_DIR/saaru && ./target/release/saaru --base-path ./docs" \
  -n sapling      "cd $BENCH_DIR/sapling/benchmark && ./../target/release/sapling run" \
  -n anna-release "cd $BENCH_DIR/anna-release && ./anna" \
  -n anna-main    "cd $BENCH_DIR/anna-main && ./anna" \
  -n anna-local   "cd $BENCH_DIR/anna && ./anna" \
  -n hugo         "cd $BENCH_DIR/hugo && hugo" \
  | tee -a "$BENCH_DIR/bench_results.txt"

# generate markdown report
{
    echo "## Benchmark Results"
    echo

    echo "| Rank | Benchmark | Mean | Std Dev | User | System |"
    echo "|------|-----------|------|---------|------|--------|"

    awk '
    /^Benchmark [0-9]+:/ {
        bench=$3
    }

    /Time \(mean ± σ\):/ {
        mean=$4
        std=$7

        match($0, /User: ([0-9.]+ ms)/)
        user=substr($0, RSTART+6, RLENGTH-6)

        match($0, /System: ([0-9.]+ ms)/)
        system=substr($0, RSTART+8, RLENGTH-8)

        print bench "|" mean "|" std "|" user "|" system
    }
    ' "$BENCH_DIR/bench_results.txt" |
    sort -t'|' -k2,2n |
    awk -F'|' '
    {
        printf("| %d | %s | %s ms | %s ms | %s | %s |\n",
            NR, $1, $2, $3, $4, $5)
    }
    '

    echo
    echo "### Relative Performance"
    echo

    awk '
    /^Summary$/ { p=1; next }
    p && NF { print }
    ' "$BENCH_DIR/bench_results.txt"

    echo
    echo "## Commit Hashes"
    echo

    cat "$BENCH_DIR/commit_hashes.txt"

} > "$BENCH_DIR/bench_report.md"


echo
echo "Results written to:"
echo "  $BENCH_DIR/bench_results.txt"
echo "Commit hashes written to:"
echo "  $BENCH_DIR/commit_hashes.txt"
echo "Markdown report written to:"
echo "  $BENCH_DIR/bench_report.md"