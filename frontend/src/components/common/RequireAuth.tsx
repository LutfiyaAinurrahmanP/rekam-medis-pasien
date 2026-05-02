import { Navigate, Outlet, useLocation } from "react-router";
import authService from "../../services/auth";

export default function RequireAuth() {
  const location = useLocation();

  if (!authService.isAuthenticated()) {
    return <Navigate to="/auth/login" replace state={{ from: location }} />;
  }

  return <Outlet />;
}
