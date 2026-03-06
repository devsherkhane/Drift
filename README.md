#  Trello Clone

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.0-4FC08D?style=flat&logo=vue.js)](https://vuejs.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A high-performance, real-time project management tool built with **Go (Gin)** and **Vue 3 (Vite)**. Featuring a modern "Vibrant Aurora" UI with high-energy glassmorphism, real-time WebSocket sync, and intuitive drag-and-drop mechanics.

---

## ✨ Key Features

- **🌈 Modern "Aurora" UI**: A fresh, high-energy interface featuring teal-indigo gradients, frosted glass panels (24px blur), and smooth spring animations.
- **🔄 Real-Time Synchronization**: Live updates across all clients via WebSockets—no refresh needed.
- **🖱️ Intuitive Drag & Drop**: Seamlessly move cards between lists and reorder tasks with fluid animations.
- **🔐 Secure & Robust Auth**: 
  - JWT-based authentication with secure cookie storage.
  - Full Password Recovery flow (Forgot/Reset).
- **📊 Board Management**:
  - Create, Archive, and Delete boards.
  - Multi-user collaboration with invite links.
  - Collaborative labels, comments, and file attachments.
  - Activity logs for every board and card change.
- **🌓 Adaptive Theming**: Native Dark Mode ("Midnight Aurora") and Light Mode support with perfect contrast.

---

## 🛠️ Tech Stack

### Backend (The Engine)
- **Go 1.18+ / Gin**: High-concurrency RESTful API.
- **MySQL**: Persistent relational data storage.
- **Gorilla WebSocket**: Low-latency bi-directional communication.
- **Swaggo**: Automated API documentation (Swagger/OpenAPI).

### Frontend (The Experience)
- **Vue 3 / Vite**: Fast, reactive component-based UI.
- **Pinia**: Centralized state management for boards and auth.
- **Tailwind CSS**: Utility-first styling with custom "Aurora" design tokens.
- **Vue Draggable (vuedraggable)**: Reliable drag-and-drop orchestration.
- **Lucide Icons**: Crisp, professional iconography.

---

## 🚀 Quick Start

### Prerequisites
- [Go](https://golang.org/dl/) (1.18+)
- [Node.js](https://nodejs.org/) (16+)
- [MySQL](https://www.mysql.com/downloads/)

### 1. Database Setup
Create a database named `trello_clone` and run the initial migration:
```bash
mysql -u root -p trello_clone < init.sql
```

### 2. Backend Installation
```bash
cd backend
go mod tidy
go run main.go
```
*API will be available at `http://localhost:8080`*

### 3. Frontend Installation
```bash
cd frontend
npm install
npm run dev
```
*App will be available at `http://localhost:5173`*

---
