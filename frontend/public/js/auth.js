document.getElementById("loginForm").addEventListener("submit", async (e) => {
  e.preventDefault(); // Prevent default form submission

  // Retrieve email and password values
  const email = document.getElementById("email").value;
  const password = document.getElementById("password").value;

  try {
    // Send login details to the backend
    const response = await fetch("http://localhost:8080/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
    });

    const data = await response.json();
    if (data.success) {
      alert(data.message);
      window.location.href = "home.html"; // Redirect to authenticated home page
    } else {
      alert(data.message);
    }
  } catch (error) {
    console.error("Error:", error);
    alert("An error occurred. Please try again.");
  }
});
