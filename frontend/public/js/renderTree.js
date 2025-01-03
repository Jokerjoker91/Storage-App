import { getToken } from "./utils.js";

// Fetch and render file tree structure
export async function renderFileTree() {
  const fileTree = document.getElementById("fileTree");
  const folderSelect = document.getElementById("folderSelect");

  try {
    var token = getToken();
    const response = await fetch("/api/get-bucket-contents", {
      headers: {
        Authorization: `Bearer ${token}`, // Include JWT token
      },
    });

    if (!response.ok) {
      if (response.status === 401) {
        alert("Unauthorized: Please log in.");
      }
      throw new Error("Failed to fetch folder structure");
    }

    const folderStructure = await response.json();

    // 1. Render the file tree
    const structureHTML = `<ul>${generateFolderHTML(folderStructure)}</ul>`;
    fileTree.innerHTML = structureHTML;

    // 2. Populate the folder dropdown (reuse the response)
    populateFolderDropdown(folderStructure, folderSelect);
  } catch (error) {
    console.error("Error fetching or rendering folder structure:", error);
    fileTree.innerHTML = "<p>Error loading folder structure</p>";
  }
}

// Helper function to recursively generate folder HTML
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

// Helper function to populate folder dropdown
function populateFolderDropdown(folderStructure, folderSelect) {
  folderSelect.innerHTML = '<option value="Root">Root</option>'; // Default to Root

  function addFolders(folders, prefix = "") {
    folders.forEach((folder) => {
      const folderPath = prefix ? `${prefix}/${folder.name}` : folder.name;

      const option = document.createElement("option");
      option.value = folderPath;
      option.textContent = folderPath;
      folderSelect.appendChild(option);

      if (folder.subFolders && folder.subFolders.length > 0) {
        addFolders(folder.subFolders, folderPath);
      }
    });
  }

  addFolders(folderStructure.subFolders);
}
