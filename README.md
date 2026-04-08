# Drift
 
A high-performance, real-time project management application featuring a **Go (Gin)** backend and a reactive **Vue 3** frontend. This platform provides a seamless task management experience through a professional glassmorphism interface, integrated drag-and-drop functionality, and live synchronization powered by WebSockets.
 
---   
           
## Key Features            
  
* **Real-Time Synchronization:** Utilizes WebSockets for instant state updates across all connected clients. Card movements and content edits are reflected globally in real-time.
* **Drag-and-Drop Interface:** Fluid task organization and list reordering for intuitive workflow management.
* **Secure Authentication:** Comprehensive security suite including JWT-based authentication, password hashing, secure session management, and password recovery flows.
* **Professional UI/UX:** Features "Vibrant Aurora" (Dark Mode) and "Clean Minimalist" (Light Mode) themes. The design employs custom CSS variables, deep slate card aesthetics, and micro-animations for refined user interaction.
* **Collaboration Management:** Tools for board sharing, granular access control, and an invitation system for team expansion.
* **Advanced Filtering:** High-speed search functionality to locate specific cards and tasks across multiple workspaces.
* **User Profiles:** Individualized account management including avatar customization and security settings.

---
  
## Technical Architecture
 
### Frontend
* **Vue 3:** Built with the Composition API and `<script setup>` for optimal performance and maintainability.
* **Vite:** Used as the build tool for rapid development and optimized production bundling.
* **Vue Router:** Manages client-side navigation and route protection.
* **Vanilla CSS:** Implements custom design tokens and flexible layouts without the overhead of heavy CSS frameworks.

### Backend
* **Go (Golang):** Leveraged for high concurrency, low latency, and robust backend performance.
* **Gin Gonic:** High-performance HTTP web framework for routing and middleware management.
* **GORM:** Provides an abstraction layer for database interactions and schema migrations.
* **Gorilla WebSockets:** Manages bidirectional communication for live updates.
* **MySQL:** Reliable relational storage for boards, cards, and user data.

---
 
