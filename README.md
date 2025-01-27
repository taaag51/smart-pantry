# Smart-pantry

Smart-pantry is a web application designed to help users manage their pantry inventory efficiently. Built with a Go backend using the Echo framework and a React frontend, Smart-pantry offers a seamless experience for tracking, adding, and organizing pantry items.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Add Items:** Easily add new items to your pantry with detailed information.
- **View Inventory:** View a comprehensive list of all pantry items.
- **Edit & Delete:** Modify or remove items as needed.
- **Search & Filter:** Quickly find items using search and filter functionalities.
- **Responsive Design:** Accessible on both desktop and mobile devices.

## Getting Started

### Prerequisites

Ensure you have the following installed on your system:

- **Go**: 1.21
- **Node.js**: 18.x
- **Docker** (optional, for containerization)

### Installation

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/yourusername/smart-pantry.git
   cd smart-pantry
   ```

2. **Backend Setup:**

   ```bash
   cd backend
   go mod download
   ```

3. **Frontend Setup:**

   ```bash
   cd ../frontend
   npm install
   ```

4. **Environment Variables:**

   - Duplicate the `.env.example` file and rename it to `.env`.
   - Fill in the necessary environment variables as per your setup.

5. **Run the Application:**

   ```bash
   cd ../
   docker-compose up --build
   ```

6. **Access the Application:**
   Open your browser and navigate to `http://localhost:8080`.

## Project Structure

```
smart-pantry/
├── backend/
│   ├── cmd/
│   │   └── server.go
│   ├── configs/
│   │   └── config.go
│   ├── internal/
│   │   ├── controllers/
│   │   │   └── pantry_controller.go
│   │   ├── models/
│   │   │   └── pantry_model.go
│   │   ├── routes/
│   │   │   └── api.go
│   │   └── services/
│   │       └── pantry_service.go
│   ├── migrations/
│   ├── scripts/
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── package.json
│   ├── src/
│   │   ├── components/
│   │   │   ├── PantryForm.tsx
│   │   │   └── PantryItemList.tsx
│   │   ├── pages/
│   │   │   └── PantryPage.tsx
│   │   ├── services/
│   │   │   └── api.ts
│   │   ├── styles/
│   │   ├── utils/
│   │   ├── App.tsx
│   │   └── index.tsx
├── docker-compose.yml
├── .env.example
├── .gitignore
└── README.md
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
