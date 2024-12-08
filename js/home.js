// Simulate file upload and folder structure rendering
document.getElementById("uploadButton").addEventListener("click", () => {
  alert("Files uploaded successfully!");
  // Simulate updating the file structure
  renderFileTree();
});

// Render file structure (dummy data for now)
function renderFileTree() {
  const fileTree = document.getElementById("fileTree");
  const structure = `
      <ul>
        <li>Folder 1
          <ul>
            <li>File 1-1.jpg</li>
            <li>File 1-2.png</li>
          </ul>
        </li>
        <li>Folder 2
          <ul>
            <li>File 2-1.docx</li>
          </ul>
        </li>
      </ul>
    `;
  fileTree.innerHTML = structure;
}

// Drag-and-Drop Handling
const uploadZone = document.getElementById("uploadZone");
uploadZone.addEventListener("dragover", (e) => {
  e.preventDefault();
  uploadZone.style.background = "#e1e1e1";
});
uploadZone.addEventListener("dragleave", () => {
  uploadZone.style.background = "#f9f9f9";
});
uploadZone.addEventListener("drop", (e) => {
  e.preventDefault();
  uploadZone.style.background = "#f9f9f9";
  const files = Array.from(e.dataTransfer.files);
  console.log("Dropped files:", files);
  alert(`${files.length} file(s) dropped!`);
});

// Function to toggle the visibility of the dropdown menu
function toggleDropdown() {
  const dropdownMenu = document.getElementById("dropdownMenu");
  dropdownMenu.classList.toggle("show");
}

// Close the dropdown if clicked outside
document.addEventListener("click", (event) => {
  const dropdownMenu = document.getElementById("dropdownMenu");
  const avatarContainer = document.querySelector(".avatar-container");

  // Close dropdown if click happens outside the avatar container
  if (!avatarContainer.contains(event.target)) {
    dropdownMenu.classList.remove("show");
  }
});
