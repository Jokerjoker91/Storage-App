// Simulate file upload and folder structure rendering
document.getElementById("uploadButton").addEventListener("click", () => {
  const files = document.getElementById("fileInput").files;
  const fileList = [];

  // Collect the files and their relative paths (folder structure)
  for (const file of files) {
    fileList.push({
      name: file.name,
      path: file.webkitRelativePath, // Relative path simulating folder structure
      file: file,
    });
  }

  // Send the file list to the backend for upload
  uploadFilesToBackend(fileList);

  // Simulate updating the file structure
  renderFileTree();
});

// Function to upload files to the backend
function uploadFilesToBackend(fileList) {
  fetch("http://localhost:8080/upload-folder", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      files: fileList.map((item) => ({
        name: item.name,
        path: item.path,
      })),
    }),
  })
    .then((response) => response.json()) // Ensure it's treated as JSON
    .then((data) => {
      if (data.success) {
        alert(data.message); // Display success message
      } else {
        alert("Error uploading files: " + data.message);
      }
    })
    .catch((error) => {
      console.error("Error uploading files:", error);
      alert("Error uploading files.");
    });
}

// Helper function to recursively generate HTML for the folder structure
function generateFolderHTML(folder) {
  let html = `<li>${folder.name}`;

  if (folder.files && folder.files.length > 0) {
    html += "<ul>";
    folder.files.forEach((file) => {
      html += `<li>${file}</li>`;
    });
    html += "</ul>";
  }

  if (folder.subFolders && folder.subFolders.length > 0) {
    html += "<ul>";
    folder.subFolders.forEach((subFolder) => {
      html += generateFolderHTML(subFolder);
    });
    html += "</ul>";
  }

  html += "</li>";
  return html;
}

// Render file structure dynamically from API
async function renderFileTree() {
  const fileTree = document.getElementById("fileTree");

  try {
    const response = await fetch("/api/get-bucket-contents");
    if (!response.ok) {
      throw new Error("Failed to fetch folder structure");
    }

    const folderStructure = await response.json();

    // Generate HTML for the folder structure
    const structureHTML = `<ul>${generateFolderHTML(folderStructure)}</ul>`;
    fileTree.innerHTML = structureHTML;
  } catch (error) {
    console.error("Error fetching or rendering folder structure:", error);
    fileTree.innerHTML = "<p>Error loading folder structure</p>";
  }
}

// Call the function to render the file tree
renderFileTree();

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

  // Convert dropped files into a structure with paths
  const fileList = files.map((file) => ({
    name: file.name,
    path: file.webkitRelativePath, // Folder structure
    file: file,
  }));

  // Send the file list to the backend for upload
  uploadFilesToBackend(fileList);
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
