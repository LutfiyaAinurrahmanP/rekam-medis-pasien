import { Fragment } from "react";
import { Route } from "react-router";
import RoleDashboardPage from "../pages/Roles/RoleDashboardPage";
import RoleReportsPage from "../pages/Roles/RoleReportsPage";

export default function DoctorRoutes() {
  return (
    <Fragment>
      <Route path="/doctor/dashboard" element={<RoleDashboardPage />} />
      <Route path="/doctor/reports" element={<RoleReportsPage />} />
    </Fragment>
  );
}
