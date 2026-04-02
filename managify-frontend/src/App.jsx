import { HashRouter, BrowserRouter, Routes, Route } from "react-router-dom";

import Login from "./components/login/Login";
import GoogleCallback from "./components/login/GoogleCallback";
import Register from "./components/register/Register";
import Dashboard from "./components/dashboard/Dashboard";
import ManagifyLandingPage from "./components/main/home";
import CreateProject from "./components/project/CreateProject";
import PlanCards from "./components/plan/PlanCards";

import { AuthProvider } from "./content/AuthContent";
import ProtectedRoute from "./content/ProtectedRoute";
import PublicRoute from "./content/PublicRoute";

import { Toaster } from "react-hot-toast";
import ProjectDetail from "./components/project/DetailProject";
import Profile from "./components/main/Profile";
import VerifyEmail from "./components/verify/VerifyEmail";
import { ThemeProvider } from "./content/ThemeContent";

const isElectron = window?.process?.versions?.electron;

export default function App() {
  const RouterComponent = isElectron ? HashRouter : BrowserRouter;

  return (
    <ThemeProvider>
      <RouterComponent>
        <AuthProvider>
          <Routes>
            <Route path="/" element={<ManagifyLandingPage />} />
            <Route path="/verify" element={
              <PublicRoute>
                <VerifyEmail />
              </PublicRoute>
            } />
            <Route
              path="/register"
              element={
                <PublicRoute>
                  <Register />
                </PublicRoute>
              }
            />

            <Route
              path="/profile"
              element={
                <ProtectedRoute>
                  <Profile />
                </ProtectedRoute>
              }
            />

            <Route
              path="/projects/:id"
              element={
                <ProtectedRoute>
                  <ProjectDetail />
                </ProtectedRoute>
              }
            />
            <Route
              path="/login"
              element={
                <PublicRoute>
                  <Login />
                </PublicRoute>
              }
            />
            <Route
              path="/auth/google/callback"
              element={
                <PublicRoute>
                  <GoogleCallback />
                </PublicRoute>
              }
            />
            <Route
              path="/dashboard"
              element={
                <ProtectedRoute>
                  <Dashboard />
                </ProtectedRoute>
              }
            />
            <Route
              path="/create-project"
              element={
                <ProtectedRoute>
                  <CreateProject />
                </ProtectedRoute>
              }
            />
            <Route
              path="/plans"
              element={
                <ProtectedRoute>
                  <PlanCards />
                </ProtectedRoute>
              }
            />
          </Routes>

          <Toaster position="top-right" reverseOrder={false} />
        </AuthProvider>
      </RouterComponent>
    </ThemeProvider>
  );
}