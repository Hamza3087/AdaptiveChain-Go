/**
 * Chart.js integration for Advanced Blockchain Frontend
 */

// Chart.js CDN
document.write('<script src="https://cdn.jsdelivr.net/npm/chart.js@3.9.1/dist/chart.min.js"></script>');

// Chart color schemes
const chartColors = {
    primary: 'rgba(58, 134, 255, 0.8)',
    secondary: 'rgba(131, 56, 236, 0.8)',
    success: 'rgba(56, 176, 0, 0.8)',
    warning: 'rgba(255, 190, 11, 0.8)',
    danger: 'rgba(255, 0, 110, 0.8)',
    primaryLight: 'rgba(58, 134, 255, 0.2)',
    secondaryLight: 'rgba(131, 56, 236, 0.2)',
    successLight: 'rgba(56, 176, 0, 0.2)',
    warningLight: 'rgba(255, 190, 11, 0.2)',
    dangerLight: 'rgba(255, 0, 110, 0.2)',
};

// Store chart instances to update them later
const chartInstances = {};

// Create a line chart
function createLineChart(canvasId, data, options = {}) {
    // Get the canvas element
    const canvas = document.getElementById(canvasId);
    if (!canvas) {
        console.error(`Canvas element with ID '${canvasId}' not found`);
        return null;
    }

    // Create chart configuration
    const config = {
        type: 'line',
        data: {
            labels: data.labels || [],
            datasets: data.datasets.map((dataset, index) => {
                const colorIndex = index % 5;
                const colors = [
                    chartColors.primary,
                    chartColors.secondary,
                    chartColors.success,
                    chartColors.warning,
                    chartColors.danger
                ];
                const backgroundColors = [
                    chartColors.primaryLight,
                    chartColors.secondaryLight,
                    chartColors.successLight,
                    chartColors.warningLight,
                    chartColors.dangerLight
                ];
                
                return {
                    label: dataset.label,
                    data: dataset.data,
                    borderColor: dataset.borderColor || colors[colorIndex],
                    backgroundColor: dataset.backgroundColor || backgroundColors[colorIndex],
                    borderWidth: dataset.borderWidth || 2,
                    tension: dataset.tension || 0.4,
                    fill: dataset.fill !== undefined ? dataset.fill : true,
                    pointRadius: dataset.pointRadius || 3,
                    pointHoverRadius: dataset.pointHoverRadius || 5
                };
            })
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'top',
                    labels: {
                        usePointStyle: true,
                        padding: 20
                    }
                },
                tooltip: {
                    mode: 'index',
                    intersect: false,
                    backgroundColor: 'rgba(0, 0, 0, 0.7)',
                    padding: 10,
                    cornerRadius: 4,
                    caretSize: 5
                },
                ...options.plugins
            },
            scales: {
                x: {
                    grid: {
                        display: false
                    },
                    ticks: {
                        padding: 10
                    }
                },
                y: {
                    beginAtZero: true,
                    grid: {
                        color: 'rgba(0, 0, 0, 0.05)'
                    },
                    ticks: {
                        padding: 10
                    }
                },
                ...options.scales
            },
            interaction: {
                mode: 'nearest',
                axis: 'x',
                intersect: false
            },
            ...options
        }
    };

    // Create the chart
    const chart = new Chart(canvas, config);
    
    // Store the chart instance for later updates
    chartInstances[canvasId] = chart;
    
    return chart;
}

// Create a bar chart
function createBarChart(canvasId, data, options = {}) {
    // Get the canvas element
    const canvas = document.getElementById(canvasId);
    if (!canvas) {
        console.error(`Canvas element with ID '${canvasId}' not found`);
        return null;
    }

    // Create chart configuration
    const config = {
        type: 'bar',
        data: {
            labels: data.labels || [],
            datasets: data.datasets.map((dataset, index) => {
                const colorIndex = index % 5;
                const colors = [
                    chartColors.primary,
                    chartColors.secondary,
                    chartColors.success,
                    chartColors.warning,
                    chartColors.danger
                ];
                
                return {
                    label: dataset.label,
                    data: dataset.data,
                    backgroundColor: dataset.backgroundColor || colors[colorIndex],
                    borderColor: dataset.borderColor || colors[colorIndex],
                    borderWidth: dataset.borderWidth || 1,
                    borderRadius: dataset.borderRadius || 4,
                    barPercentage: dataset.barPercentage || 0.6,
                    categoryPercentage: dataset.categoryPercentage || 0.8
                };
            })
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'top',
                    labels: {
                        usePointStyle: true,
                        padding: 20
                    }
                },
                tooltip: {
                    mode: 'index',
                    intersect: false,
                    backgroundColor: 'rgba(0, 0, 0, 0.7)',
                    padding: 10,
                    cornerRadius: 4,
                    caretSize: 5
                },
                ...options.plugins
            },
            scales: {
                x: {
                    grid: {
                        display: false
                    },
                    ticks: {
                        padding: 10
                    }
                },
                y: {
                    beginAtZero: true,
                    grid: {
                        color: 'rgba(0, 0, 0, 0.05)'
                    },
                    ticks: {
                        padding: 10
                    }
                },
                ...options.scales
            },
            ...options
        }
    };

    // Create the chart
    const chart = new Chart(canvas, config);
    
    // Store the chart instance for later updates
    chartInstances[canvasId] = chart;
    
    return chart;
}

// Create a doughnut chart
function createDoughnutChart(canvasId, data, options = {}) {
    // Get the canvas element
    const canvas = document.getElementById(canvasId);
    if (!canvas) {
        console.error(`Canvas element with ID '${canvasId}' not found`);
        return null;
    }

    // Default colors
    const defaultColors = [
        chartColors.primary,
        chartColors.secondary,
        chartColors.success,
        chartColors.warning,
        chartColors.danger
    ];

    // Create chart configuration
    const config = {
        type: 'doughnut',
        data: {
            labels: data.labels || [],
            datasets: [{
                data: data.data || [],
                backgroundColor: data.backgroundColor || defaultColors,
                borderColor: data.borderColor || 'white',
                borderWidth: data.borderWidth || 2,
                hoverOffset: data.hoverOffset || 4
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'right',
                    labels: {
                        usePointStyle: true,
                        padding: 20
                    }
                },
                tooltip: {
                    backgroundColor: 'rgba(0, 0, 0, 0.7)',
                    padding: 10,
                    cornerRadius: 4,
                    caretSize: 5
                },
                ...options.plugins
            },
            cutout: options.cutout || '70%',
            ...options
        }
    };

    // Create the chart
    const chart = new Chart(canvas, config);
    
    // Store the chart instance for later updates
    chartInstances[canvasId] = chart;
    
    return chart;
}

// Update an existing chart
function updateChart(canvasId, newData, options = {}) {
    const chart = chartInstances[canvasId];
    if (!chart) {
        console.error(`Chart with ID '${canvasId}' not found`);
        return;
    }

    // Update labels if provided
    if (newData.labels) {
        chart.data.labels = newData.labels;
    }

    // Update datasets if provided
    if (newData.datasets) {
        newData.datasets.forEach((newDataset, datasetIndex) => {
            if (datasetIndex < chart.data.datasets.length) {
                // Update existing dataset
                Object.keys(newDataset).forEach(key => {
                    chart.data.datasets[datasetIndex][key] = newDataset[key];
                });
            } else {
                // Add new dataset
                chart.data.datasets.push(newDataset);
            }
        });
    } else if (newData.data) {
        // For doughnut charts
        chart.data.datasets[0].data = newData.data;
    }

    // Update options if provided
    if (options && Object.keys(options).length > 0) {
        chart.options = { ...chart.options, ...options };
    }

    // Update the chart
    chart.update();
}

// Export chart functions
const ChartUtils = {
    createLineChart,
    createBarChart,
    createDoughnutChart,
    updateChart,
    chartColors
};
