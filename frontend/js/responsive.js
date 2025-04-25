/**
 * Responsive functionality for the Advanced Blockchain Frontend
 */

document.addEventListener('DOMContentLoaded', function() {
    // Create sidebar toggle button
    const sidebarToggle = document.createElement('button');
    sidebarToggle.className = 'sidebar-toggle';
    sidebarToggle.innerHTML = '☰ Menu';
    document.body.appendChild(sidebarToggle);
    
    // Get sidebar element
    const sidebar = document.querySelector('.sidebar');
    
    // Toggle sidebar when button is clicked
    sidebarToggle.addEventListener('click', function() {
        sidebar.classList.toggle('show');
        
        // Change button text based on sidebar state
        if (sidebar.classList.contains('show')) {
            sidebarToggle.innerHTML = '✕ Close';
        } else {
            sidebarToggle.innerHTML = '☰ Menu';
        }
    });
    
    // Close sidebar when clicking outside
    document.addEventListener('click', function(event) {
        if (!sidebar.contains(event.target) && event.target !== sidebarToggle) {
            sidebar.classList.remove('show');
            sidebarToggle.innerHTML = '☰ Menu';
        }
    });
    
    // Close sidebar when clicking on a link
    const sidebarLinks = document.querySelectorAll('.sidebar-menu a');
    sidebarLinks.forEach(link => {
        link.addEventListener('click', function() {
            if (window.innerWidth <= 992) {
                sidebar.classList.remove('show');
                sidebarToggle.innerHTML = '☰ Menu';
            }
        });
    });
    
    // Handle window resize
    window.addEventListener('resize', function() {
        if (window.innerWidth > 992) {
            sidebar.classList.remove('show');
            sidebarToggle.innerHTML = '☰ Menu';
        }
    });
});
