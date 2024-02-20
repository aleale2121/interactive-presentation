# Interactive Presentation Backend

## Project Overview

This project focuses on developing an Interactive Presentation Backend, designed to facilitate dynamic and engaging live presentations. The core concept involves presenters posing questions during the presentation, with each question corresponding to a dedicated slide. Attendees can participate by casting their votes during the display of each question, fostering real-time interaction.

### Database Schema

The Interactive Presentation API requires four tables:

- Presentations
- Polls
- Options
- Votes

## Key Objectives

- **Build Robust API Backend**: Develop a seamless REST API backend with reliable endpoints for efficient communication between frontend and backend components.
- **Database Integration**: Design an efficient database schema for storing server data related to questions, answers, and interactions.
- **Scalability and Performance**: Optimize the backend to handle occasional traffic spikes, ensuring smooth operations even during peak moments, such as 30,000 concurrent votes.
- **Data Integrity**: Ensure transactional integrity during voting and slide transitions to prevent data loss due to race conditions.
- **API Documentation**: Create clear documentation for the backend API, detailing functionalities and endpoints.
- **Clean Code and Testing**: Emphasize clean, maintainable code and implement testing for critical sections.

## Technologies

### Core Stacks

- **Golang**: Chosen for its efficiency, exceptional performance, seamless concurrency, and cloud-friendly nature.
- **PostgreSQL**: Provides a reliable foundation for data management, offering robustness and versatility.
- **Docker**: Utilized for containerization, simplifying the packaging of the application and its dependencies into a consistent and isolated environment.

### Frameworks

- **Gin**: A high-performance web framework in Golang, offering speed, productivity, and middleware support for flexible HTTP request handling.
- **SQLC**: Generates fully type-safe, idiomatic Go code from SQL queries, simplifying database interactions.
- **Golang Migrate**: CLI tool for efficient database migrations.
- **Viper**: Golang library for easy environmental variable loading.
- **Testify**: Golang library for easy unit testing.

For more details, refer to the [Medium blog](https://medium.com/@alefewyimer2/golang-hexagonal-architecture-sqlc-docker-gin-rest-api-interactive-presentation-787bb635080d) post.

## Dockerizing the Application

### Multistage Build

Multistage builds in Docker simplify the Dockerfile by breaking down complex build processes into smaller, focused stages. They offer advantages such as reduced image size, improved security, simplified build processes, and optimized build performance.

### Dockerfile Structure

The Dockerfile consists of two stages:

1. **Build Stage**: Compiles the Go application and prepares it for execution.
2. **Run Stage**: Sets up the runtime environment and executes the compiled application.

### Compose for Development Stacks

Docker Compose simplifies the management of multi-container Docker applications. It allows defining and running multi-container Docker applications using a YAML file.

Refer to the [Medium blog](https://medium.com/@alefewyimer2/golang-hexagonal-architecture-sqlc-docker-gin-rest-api-interactive-presentation-787bb635080d) post for detailed setup instructions.

## Iterating

After making changes to the source code:

- Use `docker-compose build` to rebuild container images.
- Use `docker-compose up` to restart the stack with the new images.

## Conclusion

By leveraging Golang, Hexagonal Architecture, SQLC, Docker, and REST API principles, this project aims to deliver a robust and interactive presentation backend, facilitating seamless communication between presenters and attendees. For more information, visit the [Medium blog](https://medium.com/@alefewyimer2/golang-hexagonal-architecture-sqlc-docker-gin-rest-api-interactive-presentation-787bb635080d) post.
