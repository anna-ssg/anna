// search.js

// Function to handle form submission
async function handleSubmit(event) {
    event.preventDefault();

    // Get the search query from the input field
    const searchInput = document.getElementById('search-input');
    const query = searchInput.value.trim();

    // Call the Go function with the search query
    const results = await searchFiles(query);

    // Write the results to a text file
    writeResultsToFile(query, results);
}

// Function to call the Go function (searchFiles) and retrieve results
async function searchFiles(query) {
    return new Promise((resolve, reject) => {
        // Call the Go function
        const result = searchFilesGo(query);
        // Resolve the promise with the result
        resolve(result);
    });
}

// Function to write the results to a text file
function writeResultsToFile(query, results) {
    // Format the results as a string
    const resultsString = results.join('\n');

    // Create a Blob with the results string
    const blob = new Blob([resultsString], { type: 'text/plain' });

    // Create a link element to download the Blob as a text file
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `search_results_${query}.txt`;
    link.click();
}

// Add event listener for form submission
const searchForm = document.getElementById('search-form');
searchForm.addEventListener('submit', handleSubmit);
