# LIMS Frontend

This document provides guidance for developers working on the LIMS frontend application.

---

## Overview

The frontend is a single-page application (SPA) built with React that serves as the user interface for the LIMS. It communicates with the Go backend via a REST API to manage laboratory data.

---

## Tech Stack

- **Framework:** [React](https://reactjs.org/) with [Vite](https://vitejs.dev/)
- **Language:** [TypeScript](https://www.typescriptlang.org/)
- **Admin Framework:** [React Admin](https://marmelab.com/react-admin/)
- **Component Framework:** [Material UI](https://mui.com/)

- **HTTP Client:** [Axios](https://axios-http.com/)
- **Styling:** [Material UI](https://mui.com/system/styles/basics/) style guide (css-in-js, nested selectors, etc.)

---

## Project Structure

The `src` directory is organized as follows:

- **`src/`**: The root of the frontend application.
- **`src/admin/`**: Contains all modules and components related to the React Admin framework. Each subdirectory typically represents a data resource (e.g., `users`, `orders`).
- **`src/components/`**: Holds reusable, general-purpose React components that are shared across the application (e.g., custom buttons, layout elements).
- **`src/hooks/`**: Contains custom React hooks. For example, `useAxios.ts` provides a standardized way to handle API requests.
- **`src/helper/`**: Includes helper functions and utility scripts that can be used throughout the application.
- **`src/types/`**: Defines all TypeScript types and interfaces, ensuring data consistency across the app.

---

## Key Libraries & Conventions

### React Admin

The core of this application is built using **React Admin**. It is essential to be familiar with its architecture and conventions. When adding new features, please follow React Admin's patterns for creating resources, views (List, Edit, Create), and data providers.

### Data Fetching

API requests are handled using **Axios**. A custom hook, `useAxios`, is provided in the `src/hooks` directory to standardize data fetching logic, including handling loading states and errors. Please use this hook for all new API interactions.

### TypeScript

This is a TypeScript-first project. All new components, functions, and variables should be strongly typed. Define shared types in the `src/types` directory to ensure consistency and prevent duplication.

### Styling

Styling is done with Material UI. Global styles are located in `index.css`, while component-specific styles can be found in their respective `.tsx` files. Please keep styles scoped to their components to avoid conflicts.
