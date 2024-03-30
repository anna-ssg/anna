const slides = document.querySelectorAll('.content');
let currentSlide = 0;

// Declare global variables
let author, siteTitle, baseURL, themeURL;

function showSlide(index) {
    slides.forEach((slide, i) => {
        slide.style.display = i === index ? 'block' : 'none';
    });
    updateProgress();
}

function updateProgress() {
    const progressContainer = document.querySelector(".progress-container");
    const oldWidth = window.getComputedStyle(progressContainer).getPropertyValue('--new-width');
    const newWidth = ((currentSlide + 1) / slides.length) * 100 + "%";
    progressContainer.style.setProperty('--old-width', oldWidth);
    progressContainer.style.setProperty('--new-width', newWidth);
    progressContainer.style.animation = 'progressAnimation 1s ease-in-out both';
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
    var urlRegex = /^(https?:\/\/)([\da-z.-]+)\.([a-z.]{2,6})(\/[^\/\s]+)*$/;
                        
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

    tsParticles.load({
        id: "tsparticles",
        options: {
            "fullScreen": {
                "zIndex": 1
            },
            "emitters": [
                {
                    "position": {
                        "x": 0,
                        "y": 30
                    },
                    "rate": {
                        "quantity": 5,
                        "delay": 0.15
                    },
                    "particles": {
                        "move": {
                            "direction": "top-right",
                            "outModes": {
                                "top": "none",
                                "left": "none",
                                "default": "destroy"
                            }
                        }
                    }
                },
                {
                    "position": {
                        "x": 100,
                        "y": 30
                    },
                    "rate": {
                        "quantity": 5,
                        "delay": 0.15
                    },
                    "particles": {
                        "move": {
                            "direction": "top-left",
                            "outModes": {
                                "top": "none",
                                "right": "none",
                                "default": "destroy"
                            }
                        }
                    }
                }
            ],
            "particles": {
                "color": {
                    "value": [
                        "#ffffff",
                        "#FF0000"
                    ]
                },
                "move": {
                    "decay": 0.05,
                    "direction": "top",
                    "enable": true,
                    "gravity": {
                        "enable": true
                    },
                    "outModes": {
                        "top": "none",
                        "default": "destroy"
                    },
                    "speed": {
                        "min": 10,
                        "max": 50
                    }
                },
                "number": {
                    "value": 0
                },
                "opacity": {
                    "value": 1
                },
                "rotate": {
                    "value": {
                        "min": 0,
                        "max": 360
                    },
                    "direction": "random",
                    "animation": {
                        "enable": true,
                        "speed": 30
                    }
                },
                "tilt": {
                    "direction": "random",
                    "enable": true,
                    "value": {
                        "min": 0,
                        "max": 360
                    },
                    "animation": {
                        "enable": true,
                        "speed": 30
                    }
                },
                "size": {
                    "value": {
                        "min": 0,
                        "max": 2
                    },
                    "animation": {
                        "enable": true,
                        "startValue": "min",
                        "count": 1,
                        "speed": 16,
                        "sync": true
                    }
                },
                "roll": {
                    "darken": {
                        "enable": true,
                        "value": 25
                    },
                    "enable": true,
                    "speed": {
                        "min": 5,
                        "max": 15
                    }
                },
                "wobble": {
                    "distance": 30,
                    "enable": true,
                    "speed": {
                        "min": -7,
                        "max": 7
                    }
                },
                "shape": {
                    "type": "emoji",
                    "options": {
                        "emoji": {
                            "particles": {
                                "size": {
                                    "value": 8
                                }
                            },
                            "value": [
                                "🍚"
                            ]
                        }
                    }
                }
            }
        }
    });


    tsParticles.load("confetti", confettiSettings);
    currentSlide = currentSlide + 1;
    updateProgress()
    setTimeout(() => {
        window.location.href = 'http://localhost:8000';
    }, 5000);
}

showSlide(0);

// Add event listeners to call checkFormValidity when input fields change
document.getElementById("author").addEventListener("input", checkFormValidity);
document.getElementById("siteTitle").addEventListener("input", checkFormValidity);
document.getElementById("baseURL").addEventListener("input", checkFormValidity);
document.getElementById("themeURL").addEventListener("input", checkFormValidity);

// Confetti animation
const confettiSettings = {
    particles: {
        number: {
            value: 10
        },
        size: {
            value: 2
        },
        shape: {
            type: "circle"
        },
        move: {
            speed: 6
        },
        color: {
            value: "#00FFFF"
        },
        opacity: {
            value: 0.8
        }
    }
};
