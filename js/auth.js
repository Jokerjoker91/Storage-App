// Handle login form submission
document.getElementById("loginForm").addEventListener("submit", (e) => {
  e.preventDefault(); // Prevent the default form submission

  // Retrieve email and password values
  const email = document.getElementById("email").value;
  const password = document.getElementById("password").value;

  console.log("Email:", email);
  console.log("Password:", password);

  // Simulate successful login
  alert("Login successful! Redirecting to the home page...");
  window.location.href = "home.html"; // Redirect to the authenticated home page
});
