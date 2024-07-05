document.addEventListener('DOMContentLoaded', function () {
  fetch('https://raw.githubusercontent.com/anna-ssg/themes/main/themes.json')
    .then(response => response.json())
    .then(data => {
      const selectElement = document.getElementById('themeURL');
      data.themes.forEach(theme => {
        const option = document.createElement('option');
        option.value = `${theme.name}`;
        option.textContent = theme.name;
        selectElement.appendChild(option);
      });
    })
    .catch(error => console.error('Error fetching themes:', error));
});

const slides = document.querySelectorAll(".content");
let currentSlide = 0;
let author, siteTitle, baseURL, themeURL;

function showSlide(index) {
  slides.forEach((slide, i) => {
    slide.style.display = i === index ? "block" : "none";
  });
  updateProgress();
}

function updateProgress() {
  const progressContainer = document.querySelector(".progress-container");
  const oldWidth = window.getComputedStyle(progressContainer).getPropertyValue("--new-width");
  const newWidth = ((currentSlide + 1) / slides.length) * 100 + "%";
  progressContainer.style.setProperty("--old-width", oldWidth);
  progressContainer.style.setProperty("--new-width", newWidth);
  progressContainer.style.animation = "progressAnimation 1s ease-in-out both";
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
  const nameRegex = /^[a-zA-Z0-9\s]+$/;
  const urlRegex = /^(https?:\/\/)([\da-z.-]+)\.([a-z.]{2,6})(\/[^\/\s]+)*$/;
  const authorButton = document.getElementById("authorButton");
  const siteTitleButton = document.getElementById("siteTitleButton");
  const baseURLButton = document.getElementById("baseURLButton");
  authorButton.disabled = !(author && author.match(nameRegex));
  siteTitleButton.disabled = !(siteTitle && siteTitle.match(nameRegex));
  baseURLButton.disabled = !(baseURL && baseURL.match(urlRegex));
  document.getElementById("author").classList.toggle("valid", author && author.match(nameRegex));
  document.getElementById("siteTitle").classList.toggle("valid", siteTitle && siteTitle.match(nameRegex));
  document.getElementById("baseURL").classList.toggle("valid", baseURL && baseURL.match(urlRegex));
}

function submitForm() {
  const checkboxes = document.querySelectorAll('.nav-checkboxes input[type="checkbox"]');
  const navbarOptions = Array.from(checkboxes)
    .filter(checkbox => checkbox.checked)
    .map(checkbox => {
      let navbarElementsMap = {};
      navbarElementsMap[checkbox.className] = checkbox.value;
      return navbarElementsMap;
    });

  const formData = JSON.stringify({
    author: author,
    siteTitle: siteTitle,
    baseURL: baseURL,
    themeURL: themeURL,
    navbar: navbarOptions,
  });

  const confettiSettings = {
    particles: {
      number: {
        value: 10,
      },
      size: {
        value: 2,
      },
      shape: {
        type: "circle",
      },
      move: {
        speed: 6,
      },
      color: {
        value: "#00FFFF",
      },
      opacity: {
        value: 0.8,
      },
    },
  };

  tsParticles.load("confetti", confettiSettings);
  showSlide(slides.length - 1);

  setTimeout(() => {
    window.location.href = "http://localhost:8000"; // Change URL as needed
  }, 2000); // 2 seconds delay before redirecting

  fetch("/submit", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: formData
  })
}

// Initialize the wizard by showing the first slide
showSlide(currentSlide);

// Add event listeners to call checkFormValidity when input fields change
document.getElementById("author").addEventListener("input", checkFormValidity);
document.getElementById("siteTitle").addEventListener("input", checkFormValidity);
document.getElementById("baseURL").addEventListener("input", checkFormValidity);
document.getElementById("themeURL").addEventListener("input", checkFormValidity);
