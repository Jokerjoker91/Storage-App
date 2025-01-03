import {
  handleFileInput,
  uploadFilesToBackend,
  setupDragAndDrop,
} from "./upload.js";
import { renderFileTree } from "./renderTree.js";

document.addEventListener("DOMContentLoaded", () => {
  const fileInput = document.getElementById("fileInput");

  // Initialize rendering file tree and populating dropdown
  renderFileTree();

  // Initialize file input handler and drag-and-drop
  handleFileInput(fileInput, uploadFilesToBackend, renderFileTree);
  setupDragAndDrop(uploadFilesToBackend);
});

// Function to toggle the visibility of the dropdown menu
function toggleDropdown() {
  const dropdownMenu = document.getElementById("dropdownMenu");
  dropdownMenu.classList.toggle("show");
}

window.toggleDropdown = toggleDropdown;

// Close the dropdown if clicked outside
document.addEventListener("click", (event) => {
  const dropdownMenu = document.getElementById("dropdownMenu");
  const avatarContainer = document.querySelector(".avatar-container");

  // Close dropdown if click happens outside the avatar container
  if (!avatarContainer.contains(event.target)) {
    dropdownMenu.classList.remove("show");
  }
});
