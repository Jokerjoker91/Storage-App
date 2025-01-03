// Retrieve the JWT token from localStorage
export function getToken() {
  const token = localStorage.getItem("jwtToken");
  if (!token) {
    alert("Unauthorized: Please log in.");
    window.location.href = "./login.html"; // Redirect to login page
    return null;
  }
  return token;
}
