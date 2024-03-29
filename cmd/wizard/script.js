const slides = document.querySelectorAll('.content');
let currentSlide = 0;

// Declare global variables
let author, siteTitle, baseURL, themeURL;

function showSlide(index) {
    slides.forEach((slide, i) => {
        slide.style.display = i === index ? 'block' : 'none';
    });
}

function nextSlide() {
    currentSlide = Math.min(currentSlide + 1, slides.length - 1);
    showSlide(currentSlide);
}

function prevSlide() {
    currentSlide = Math.max(currentSlide - 1, 0);
    showSlide(currentSlide);
}

function checkFormValidity() {
    author = document.getElementById("author").value;
    siteTitle = document.getElementById("siteTitle").value;
    baseURL = document.getElementById("baseURL").value;
    themeURL = document.getElementById("themeURL").value;

    var nameRegex = /^[a-zA-Z0-9\s]+$/;
    var urlRegex = /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*$/;
    
    var authorButton = document.getElementById("authorButton");
    var siteTitleButton = document.getElementById("siteTitleButton");
    var baseURLButton = document.getElementById("baseURLButton");

    authorButton.disabled = !(author && author.match(nameRegex));
    siteTitleButton.disabled = !(siteTitle && siteTitle.match(nameRegex));
    baseURLButton.disabled = !(baseURL && baseURL.match(urlRegex));

    // Add or remove 'valid' class based on validity
    document.getElementById("author").classList.toggle("valid", author && author.match(nameRegex));
    document.getElementById("siteTitle").classList.toggle("valid", siteTitle && siteTitle.match(nameRegex));
    document.getElementById("baseURL").classList.toggle("valid", baseURL && baseURL.match(urlRegex));
}

function submitForm() {
    var formData = JSON.stringify({
        "author": author,
        "siteTitle": siteTitle,
        "baseURL": baseURL,
        "themeURL": themeURL,
        "navbar": "index"
    });

    showSlide(slides.length - 1);
    fetch('/submit', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: formData
    });

    setTimeout(() => {
        window.location.href = 'http://localhost:8000';
    }, 3000);
}

showSlide(0);

// Add event listeners to call checkFormValidity when input fields change
document.getElementById("author").addEventListener("input", checkFormValidity);
document.getElementById("siteTitle").addEventListener("input", checkFormValidity);
document.getElementById("baseURL").addEventListener("input", checkFormValidity);
document.getElementById("themeURL").addEventListener("input", checkFormValidity);
