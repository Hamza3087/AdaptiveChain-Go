# Advanced Blockchain Frontend

This is a comprehensive frontend interface for the Advanced Blockchain project. It provides a user-friendly way to interact with the sophisticated blockchain system, fully integrated with the backend API.

## Features

- **Dashboard**: Overview of blockchain statistics, recent transactions, and network health
- **Blocks Explorer**: View and explore blockchain blocks with detailed information
- **Transactions Explorer**: Track and analyze blockchain transactions
- **Node Management**: Monitor and manage blockchain nodes
- **Adaptive Merkle Forest Visualization**: Visualize and interact with the Adaptive Merkle Forest

## Pages

1. **Dashboard** (`index.html`): Main dashboard with key metrics and visualizations
2. **Blocks** (`blocks.html`): Explore blockchain blocks and their details
3. **Transactions** (`transactions.html`): View and analyze blockchain transactions
4. **Nodes** (`nodes.html`): Monitor and manage blockchain nodes
5. **Adaptive Merkle Forest** (`amf.html`): Visualize and interact with the AMF

## How to Use

1. Make sure your blockchain node is running
2. Open the `index.html` file in your web browser to access the dashboard
3. Navigate between pages using the sidebar menu
4. Interact with the visualizations and tables to explore blockchain data
5. Use the search functionality to find specific blocks, transactions, or nodes
6. View detailed information by clicking on items in the tables or visualizations

## API Integration

The frontend is fully integrated with the blockchain backend API. The integration is handled through the following files:

- `js/api.js`: Contains all API endpoints and functions for interacting with the backend
- `js/utils.js`: Contains utility functions for formatting data and handling UI elements

The API integration provides the following functionality:

- Real-time data fetching from the blockchain node
- Transaction creation and submission
- Node management (add, restart, remove)
- Block and transaction exploration
- AMF visualization and interaction

## Configuration

You can configure the frontend by modifying the following:

1. **API Endpoint**: Update the `API_BASE_URL` in `js/api.js` to point to your blockchain node's API endpoint
2. **Refresh Intervals**: Adjust the refresh intervals in the JavaScript code to control how often data is updated
3. **UI Customization**: Modify the CSS styles to match your branding

## Error Handling

The frontend includes comprehensive error handling:

- Connection errors when the blockchain node is not running
- API errors when requests fail
- User input validation
- Friendly error messages displayed to the user

## Technologies Used

- **HTML5**: For structure and content
- **CSS3**: For styling (with modern features like CSS Grid and Flexbox)
- **JavaScript**: For dynamic behavior and API integration (vanilla, no frameworks)
- **Responsive Design**: Works on all device sizes
- **Asynchronous API Calls**: Using modern async/await patterns

## Browser Compatibility

The frontend is compatible with all modern browsers, including:

- Google Chrome
- Mozilla Firefox
- Microsoft Edge
- Safari

## Extending the Frontend

You can extend the frontend by:

1. Adding new pages for additional features
2. Extending the API integration in `js/api.js`
3. Adding new visualization components
4. Implementing additional user interactions

## License

This frontend is part of the Advanced Blockchain project and is subject to the same license terms.
