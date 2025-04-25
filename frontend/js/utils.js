/**
 * Utility functions for Advanced Blockchain Frontend
 */

// Format a timestamp to a human-readable date/time
function formatTimestamp(timestamp) {
    const date = new Date(timestamp);
    return date.toLocaleString();
}

// Format a timestamp to a relative time (e.g., "2 mins ago")
function formatRelativeTime(timestamp) {
    const now = new Date();
    const date = new Date(timestamp);
    const diffSeconds = Math.floor((now - date) / 1000);

    if (diffSeconds < 60) {
        return 'Just now';
    } else if (diffSeconds < 3600) {
        const mins = Math.floor(diffSeconds / 60);
        return `${mins} min${mins > 1 ? 's' : ''} ago`;
    } else if (diffSeconds < 86400) {
        const hours = Math.floor(diffSeconds / 3600);
        return `${hours} hour${hours > 1 ? 's' : ''} ago`;
    } else {
        const days = Math.floor(diffSeconds / 86400);
        return `${days} day${days > 1 ? 's' : ''} ago`;
    }
}

// Format a hash to a shortened version (e.g., "0x8a7d...7f6a")
function formatHash(hash, length = 8) {
    if (!hash || hash.length < length * 2) return hash;
    return `${hash.substring(0, length)}...${hash.substring(hash.length - length)}`;
}

// Format a number with commas (e.g., 1,234,567)
function formatNumber(number) {
    return number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

// Format a file size (e.g., "32.5 KB")
function formatFileSize(bytes) {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
    return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`;
}

// Format an amount with currency symbol (e.g., "1.25 ETH")
function formatAmount(amount, currency = 'ETH') {
    return `${parseFloat(amount).toFixed(2)} ${currency}`;
}

// Format an uptime duration (e.g., "3d 12h 45m")
function formatUptime(seconds) {
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);

    return `${days}d ${hours}h ${minutes}m`;
}

// Show an error message to the user
function showError(message, duration = 5000) {
    // Create error element if it doesn't exist
    let errorElement = document.getElementById('error-notification');
    if (!errorElement) {
        errorElement = document.createElement('div');
        errorElement.id = 'error-notification';
        errorElement.style.position = 'fixed';
        errorElement.style.top = '20px';
        errorElement.style.right = '20px';
        errorElement.style.backgroundColor = 'var(--danger-color)';
        errorElement.style.color = 'white';
        errorElement.style.padding = '15px 20px';
        errorElement.style.borderRadius = 'var(--border-radius)';
        errorElement.style.boxShadow = 'var(--box-shadow)';
        errorElement.style.zIndex = '1000';
        errorElement.style.maxWidth = '300px';
        errorElement.style.wordBreak = 'break-word';
        document.body.appendChild(errorElement);
    }

    // Set message and show
    errorElement.textContent = message;
    errorElement.style.display = 'block';

    // Hide after duration
    setTimeout(() => {
        errorElement.style.display = 'none';
    }, duration);
}

// Show a success message to the user
function showSuccess(message, duration = 3000) {
    // Create success element if it doesn't exist
    let successElement = document.getElementById('success-notification');
    if (!successElement) {
        successElement = document.createElement('div');
        successElement.id = 'success-notification';
        successElement.style.position = 'fixed';
        successElement.style.top = '20px';
        successElement.style.right = '20px';
        successElement.style.backgroundColor = 'var(--success-color)';
        successElement.style.color = 'white';
        successElement.style.padding = '15px 20px';
        successElement.style.borderRadius = 'var(--border-radius)';
        successElement.style.boxShadow = 'var(--box-shadow)';
        successElement.style.zIndex = '1000';
        successElement.style.maxWidth = '300px';
        successElement.style.wordBreak = 'break-word';
        document.body.appendChild(successElement);
    }

    // Set message and show
    successElement.textContent = message;
    successElement.style.display = 'block';

    // Hide after duration
    setTimeout(() => {
        successElement.style.display = 'none';
    }, duration);
}

// Create pagination elements
function createPagination(currentPage, totalPages, onPageChange) {
    const paginationElement = document.createElement('div');
    paginationElement.className = 'pagination';

    // Previous button
    const prevButton = document.createElement('div');
    prevButton.className = 'pagination-item';
    prevButton.textContent = '«';
    prevButton.onclick = () => {
        if (currentPage > 1) {
            onPageChange(currentPage - 1);
        }
    };
    paginationElement.appendChild(prevButton);

    // Page numbers
    const startPage = Math.max(1, currentPage - 2);
    const endPage = Math.min(totalPages, startPage + 4);

    for (let i = startPage; i <= endPage; i++) {
        const pageItem = document.createElement('div');
        pageItem.className = `pagination-item ${i === currentPage ? 'active' : ''}`;
        pageItem.textContent = i;
        pageItem.onclick = () => {
            if (i !== currentPage) {
                onPageChange(i);
            }
        };
        paginationElement.appendChild(pageItem);
    }

    // Next button
    const nextButton = document.createElement('div');
    nextButton.className = 'pagination-item';
    nextButton.textContent = '»';
    nextButton.onclick = () => {
        if (currentPage < totalPages) {
            onPageChange(currentPage + 1);
        }
    };
    paginationElement.appendChild(nextButton);

    return paginationElement;
}

// Create a loading spinner
function createLoadingSpinner() {
    const spinner = document.createElement('div');
    spinner.className = 'loading-spinner';
    spinner.innerHTML = `
        <style>
            .loading-spinner {
                display: flex;
                justify-content: center;
                align-items: center;
                height: 100px;
            }
            .spinner {
                width: 40px;
                height: 40px;
                border: 4px solid rgba(0, 0, 0, 0.1);
                border-radius: 50%;
                border-top-color: var(--primary-color);
                animation: spin 1s ease-in-out infinite;
            }
            @keyframes spin {
                to { transform: rotate(360deg); }
            }
        </style>
        <div class="spinner"></div>
    `;
    return spinner;
}

// Show loading state
function showLoading(container) {
    // Check if container exists
    if (!container) {
        console.error('Cannot show loading state: container is null or undefined');
        return;
    }

    // Clear container
    container.innerHTML = '';

    // Add spinner
    container.appendChild(createLoadingSpinner());
}

// Create a chart using Chart.js
function createChart(container, data, options = {}) {
    // Clear the container
    container.innerHTML = '';

    // Create a canvas element for the chart
    const canvas = document.createElement('canvas');
    const canvasId = 'chart-' + Math.random().toString(36).substring(2, 9);
    canvas.id = canvasId;
    canvas.style.width = '100%';
    canvas.style.height = '300px';

    // Add the canvas to the container
    container.appendChild(canvas);

    // Determine chart type based on data and options
    let chartType = options.chartType || 'line';

    // Create the chart
    if (typeof ChartUtils !== 'undefined') {
        if (chartType === 'line') {
            ChartUtils.createLineChart(canvasId, data, options);
        } else if (chartType === 'bar') {
            ChartUtils.createBarChart(canvasId, data, options);
        } else if (chartType === 'doughnut') {
            ChartUtils.createDoughnutChart(canvasId, data, options);
        } else {
            ChartUtils.createLineChart(canvasId, data, options);
        }
    } else {
        // Fallback if ChartUtils is not available
        const fallbackContainer = document.createElement('div');
        fallbackContainer.style.width = '100%';
        fallbackContainer.style.height = '300px';
        fallbackContainer.style.backgroundColor = '#f8f9fa';
        fallbackContainer.style.borderRadius = 'var(--border-radius)';
        fallbackContainer.style.display = 'flex';
        fallbackContainer.style.justifyContent = 'center';
        fallbackContainer.style.alignItems = 'center';
        fallbackContainer.style.color = '#888';
        fallbackContainer.textContent = options.placeholder || '[Chart Placeholder]';

        container.innerHTML = '';
        container.appendChild(fallbackContainer);
    }
}

// Export all utility functions
const Utils = {
    formatTimestamp,
    formatRelativeTime,
    formatHash,
    formatNumber,
    formatFileSize,
    formatAmount,
    formatUptime,
    showError,
    showSuccess,
    createPagination,
    createLoadingSpinner,
    showLoading,
    createChart
};
