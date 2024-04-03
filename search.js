let go;

// Load the Wasm module
async function init() {
    go = new Go();
    const wasm = await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject);
    go.run(wasm.instance);
}

// Function to search the website using the Wasm module
function search(keyword) {
    const result = go.exports.searchWebsite("http://localhost:8000", keyword);
    // Display the result, you can customize this based on your UI requirements
    console.log("Search result:", result);
    changePageContent(result);
}

// Function to change the content of the page with the provided links
function changePageContent(links) {
    // Replace the existing page content with the links
    document.body.innerHTML = links.join("<br>");
}

// Initialize the Wasm module
init();
