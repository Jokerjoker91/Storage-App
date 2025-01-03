import { getToken } from "./utils.js";

// Handle file selection and prepare for upload
export function handleFileInput(
  fileInput,
  uploadFilesToBackend,
  renderFileTree
) {
  document.getElementById("uploadButton").addEventListener("click", () => {
    const files = fileInput.files;

    if (files.length === 0) {
      alert("No files selected.");
      return;
    }

    const fileList = [];

    // Collect files and paths (handle both file and folder selection)
    for (const file of files) {
      fileList.push({
        name: file.name,
        path: file.webkitRelativePath || file.name, // Folder path or standalone file name
        file: file,
      });
    }

    // Send file list to backend for upload
    uploadFilesToBackend(fileList);

    // Optionally update the file tree after upload
    renderFileTree();
  });
}

// Upload files to the backend
export function uploadFilesToBackend(fileList) {
  const formData = new FormData();

  // Retrieve the selected folder from the dropdown
  const selectedFolder = document.getElementById("folderSelect").value;

  // Append files to the FormData object
  for (const fileData of fileList) {
    // Encode the filename
    const encodedFilename = encodeURIComponent(fileData.name);

    // Determine the upload path based on the selected folder
    const uploadPath =
      !selectedFolder || selectedFolder === "Root"
        ? encodedFilename
        : `${selectedFolder}/${encodedFilename}`;
    formData.append("files", fileData.file, uploadPath); // Include path as the name
  }

  // Include the selected folder as part of the payload
  //formData.append("folder", selectedFolder);

  var token = getToken();
  fetch("/api/upload-folder", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
    },
    body: formData,
  })
    .then((response) => {
      if (!response.ok) {
        if (response.status === 401) {
          alert("Unauthorized: Please log in.");
          window.location.href = "/login.html";
        }
        throw new Error("Server error occurred");
      }
      return response.json();
    })
    .then((data) => {
      console.log("Response JSON:", data);
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

// Drag-and-drop functionality
export function setupDragAndDrop(uploadFilesToBackend) {
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

    if (files.length === 0) {
      alert("No files dropped.");
      return;
    }

    const fileList = files.map((file) => ({
      name: file.name,
      path: file.webkitRelativePath || file.name,
      file: file,
    }));

    // Send the file list to the backend for upload
    uploadFilesToBackend(fileList);
  });
}
