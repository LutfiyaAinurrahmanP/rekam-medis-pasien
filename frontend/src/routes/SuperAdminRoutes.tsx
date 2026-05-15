import { Fragment } from "react";
import { Route } from "react-router";
import RoleDashboardPage from "../pages/Roles/RoleDashboardPage";
import RoleReportsPage from "../pages/Roles/RoleReportsPage";
import UsersIndex from "../pages/Admin/Users/Index";

export default function SuperAdminRoutes() {
  return (
    <Fragment>
      <Route path="/super-admin/dashboard" element={<RoleDashboardPage />} />
      <Route path="/super-admin/reports" element={<RoleReportsPage />} />
      <Route path="/super-admin/users" element={<UsersIndex />} />
    </Fragment>
  );
}
