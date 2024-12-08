const hamburger = document.getElementById("hamburger");
const navLinks = document.getElementById("nav-links");

// Toggle navbar visibility on hamburger click
hamburger.addEventListener("click", () => {
  navLinks.classList.toggle("show");
});

document.querySelector(".oauth-google").addEventListener("click", () => {
  alert("Google OAuth login will be implemented soon!");
});
