import { Fragment } from "react";
import { Route } from "react-router";
import RoleDashboardPage from "../pages/Roles/RoleDashboardPage";
import RoleReportsPage from "../pages/Roles/RoleReportsPage";

export default function SuperAdminRoutes() {
  return (
    <Fragment>
      <Route path="/super_admin/dashboard" element={<RoleDashboardPage />} />
      <Route path="/super_admin/reports" element={<RoleReportsPage />} />
    </Fragment>
  );
}
