const hamburger = document.getElementById("hamburger");
const navLinks = document.getElementById("nav-links");

// Toggle navbar visibility on hamburger click
hamburger.addEventListener("click", () => {
  navLinks.classList.toggle("show");
});
