// Load the header and sidebar
fetch("templates/header.txt")
  .then((response) => response.text())
  .then((data) => {
    const headerPlaceholder = document.getElementById("header-placeholder");
    headerPlaceholder.outerHTML = data;
  })
  .catch((error) => {
    console.error("Error loading header:", error);
  });

fetch("templates/sidebar.txt")
  .then((response) => response.text())
  .then((data) => {
    const sidebarPlaceholder = document.getElementById("sidebar-placeholder");
    sidebarPlaceholder.outerHTML = data;
  })
  .catch((error) => {
    console.error("Error loading sidebar:", error);
  });

// Add event listeners after the content is loaded
document.addEventListener("DOMContentLoaded", () => {
  // Add any necessary event listeners or additional functionality here
  console.log("Header and sidebar loaded successfully");
});
