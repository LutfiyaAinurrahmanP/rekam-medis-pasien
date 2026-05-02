import { Navigate, Outlet } from "react-router";
import authService from "../../services/auth";

export default function GuestOnly() {
  if (authService.isAuthenticated()) {
    return <Navigate to="/dashboard" replace />;
  }

  return <Outlet />;
}
