function lightScheme() {
  document.documentElement.style.setProperty("--background", "#D2E0FB");
  document.documentElement.style.setProperty("--primary-colour", "#000000");
  document.documentElement.style.setProperty("--font-colour", "#000000");
  document.documentElement.style.setProperty(
    "--heading-font-colour",
    "#000000",
  );
  document.documentElement.style.setProperty("--navbar-link-colour", "#3a5da0");
}

function darkScheme() {
  document.documentElement.style.setProperty("--background", "#000000");
  document.documentElement.style.setProperty("--primary-colour", "#ffffe6");
  document.documentElement.style.setProperty("--font-colour", "#ffffffba");
  document.documentElement.style.setProperty(
    "--heading-font-colour",
    "#ffffe6",
  );
  document.documentElement.style.setProperty("--navbar-link-colour", "#ffffe6");
}

// store user toggle preference in local storage. Site must remember what theme was last used
// and apply it when the user returns to the site.
function ThemeSwitch() {
  // on click of button with class toggle-theme, store theme in local and switch
  var theme = localStorage.getItem("theme");
  if (theme === "light") {
    localStorage.setItem("theme", "dark");
    darkScheme();
  } else {
    localStorage.setItem("theme", "light");
    lightScheme();
  }
}

document.addEventListener("DOMContentLoaded", () => {
  window.onload = function () {
    const toggle = document.getElementById("theme-toggle");
    toggle.onclick = function () {
      // toggle.style.transition = "transform .5s cubic-bezier(1,0,0,1)";
      // if (toggle.style.transform === "rotate(1080deg)") {
      //     toggle.style.transform = "rotate(0deg)";
      // } else {
      //     toggle.style.transform = "rotate(1080deg)";
      // }
      toggle.style.scale = "scale .5s cubic-bezier(1,0,0,1)";
      toggle.style.scale = "scale(.95)";
      toggle.style.scale = "scale(1)";
      ThemeSwitch();
    };
  };
  // check local storage for theme
  const theme = localStorage.getItem("theme");
  // if no theme exists
  if (!theme) {
    lightScheme();
  }
  if (theme === "light") {
    lightScheme();
  } else {
    darkScheme();
  }
});
