# Photo Storage App

## Overview

The **Photo Storage App** allows users to securely upload, organize, and access family photos. It uses AWS S3 Glacier for storage and provides a simple, user-friendly interface for managing photos. This project includes frontend pages with login functionality, a home page with file upload capabilities, and a dashboard for tracking login and file upload activity.

## Features

- **Authentication**: Users can log in with an email and password or through Google OAuth.
- **File Upload/Download**: Users can upload files to AWS S3 Glacier.
- **Folder Structure View**: Displays the user's folders and files that are uploaded to AWS S3 Glacier.
- **Dashboard**: A page that shows login data, photo upload/download data, and top folders/files sorted by memory size.

## File Structure

/photo-storage-app /css /main.css # Main styles for the app /dashboard.css # Styles specific to the dashboard page /js /auth.js # JavaScript for authentication logic /dashboard.js # JavaScript for dashboard logic /images # Images for the app /index.html # Homepage for unauthenticated users /login.html # Login page /home.html # Home page after login /dashboard.html # Dashboard for viewing app data /profile.html # Profile page (to be added later) /readme.md # Project documentation (this file)

## Frontend Structure

### 1. **index.html**

- **Unauthenticated homepage** with a navigation bar, a "Login" button, and a welcoming message.
- Includes a collapsible navbar that adjusts to mobile screen sizes with a hamburger menu.

### 2. **login.html**

- Login page where users can authenticate using an email/password or Google OAuth.
- Upon successful login, users are redirected to the home page.

### 3. **home.html**

- The main page after login where users can upload files and view their uploaded folder structure.
- Features:
  - **Avatar**: Displays the user’s initials (e.g., "AS") in a round avatar, which shows a dropdown with options for "My Profile" and "Logout".
  - **File Upload Section**: Users can drag-and-drop files to upload them to AWS S3 Glacier.
  - **Folder Structure Section**: Displays the folders and files uploaded to AWS.

### 4. **dashboard.html**

- Displays a dashboard with three sections:
  1.  **Login Data**: Shows a table of login information (user, IP, device, date).
  2.  **Photo Upload/Download Data**: Tracks user activity for file uploads/downloads.
  3.  **Top Folders and Files by Size**: Displays the largest folders and files by memory size.

## CSS

### 1. **main.css**

- Styles the general layout, navbar, and homepage.
- Includes responsive styling for mobile screens.

### 2. **dashboard.css**

- Styles the dashboard page with specific table layouts for login data, upload/download data, and top folders/files.

## JavaScript

### 1. **auth.js**

- Placeholder file for handling authentication logic, such as email/password validation and OAuth integration.

### 2. **dashboard.js**

- Placeholder file for handling the dynamic behavior of the dashboard page (e.g., fetching and displaying data for login activity, file uploads, etc.).

## AWS Integration

- **S3 Glacier**: Files uploaded via the home page are stored in AWS S3 Glacier, a cold storage solution ideal for long-term storage of infrequently accessed data.

## How to Run Locally

1. Clone this repository to your local machine:
   ```bash
   git clone https://github.com/Jokerjoker91/Storage-App.git
   cd photo-storage-app
   Open the index.html file in your browser to view the app.
   ```

The app is designed to run with the frontend alone for now. The backend logic (for authentication and AWS S3 Glacier interaction) will be implemented later in the project.

To Do
Implement the backend in Golang for authentication and communication with AWS S3 Glacier.
Implement the "Profile" page where users can edit their personal information.
Link the frontend to real authentication and file management via AWS.
License
