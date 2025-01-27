# Project Structure

The Smart-pantry project is organized as follows:

```
smart-pantry/
├── backend/                   # Go backend
│   ├── cmd/
│   │   └── server.go          # Entry point for the backend server
│   ├── configs/
│   │   └── config.go          # Configuration settings
│   ├── internal/
│   │   ├── controllers/
│   │   │   └── pantry_controller.go   # Handles pantry-related requests
│   │   ├── models/
│   │   │   └── pantry_model.go        # Pantry data models
│   │   ├── routes/
│   │   │   └── api.go                 # API route definitions
│   │   └── services/
│   │       └── pantry_service.go      # Business logic for pantry
│   ├── migrations/                     # Database migration files
│   ├── scripts/                        # Helper scripts
│   ├── Dockerfile                      # Docker configuration for backend
│   ├── go.mod                          # Go module definition
│   └── go.sum                          # Go dependencies
├── frontend/                              # React frontend
│   ├── package.json                       # Node.js dependencies
│   ├── src/
│   │   ├── components/
│   │   │   ├── PantryForm.tsx             # Form for adding pantry items
│   │   │   └── PantryItemList.tsx         # List display for pantry items
│   │   ├── pages/
│   │   │   └── PantryPage.tsx              # Pantry management page
│   │   ├── services/
│   │   │   └── api.ts                      # API service for frontend
│   │   ├── styles/                          # CSS/SCSS styles
│   │   ├── utils/                           # Utility functions
│   │   ├── App.tsx                          # Main React component
│   │   └── index.tsx                        # React entry point
├── docker-compose.yml                      # Docker Compose configuration
├── .env.example                            # Example environment variables
├── .gitignore                              # Git ignore rules
└── README.md                               # Project documentation
```
