document.addEventListener("DOMContentLoaded", () => {
  document.getElementById("loginForm").addEventListener("submit", async (e) => {
    e.preventDefault(); // Prevent default form submission

    // Retrieve email and password values
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    try {
      // Send login details to the backend
      const response = await fetch("/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();
      if (data.success) {
        // Store the JWT token in localStorage
        localStorage.setItem("jwtToken", data.token);

        // Redirect to the home page
        window.location.href = "home.html"; // Redirect to authenticated home page
      } else {
        alert(data.message);
      }
    } catch (error) {
      console.error("Login Error:", error);
      alert("An error occurred. Please try again.");
    }
  });
});
