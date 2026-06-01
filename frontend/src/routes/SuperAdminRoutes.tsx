import { Fragment } from "react";
import { Route } from "react-router";
import RoleDashboardPage from "../pages/Roles/RoleDashboardPage";
import RoleReportsPage from "../pages/Roles/RoleReportsPage";
import UsersIndex from "../pages/Admin/Users/Index";
import DepartmentsIndex from "../pages/Admin/Departments/Index";
import PatientsIndex from "../pages/Admin/Patients/Index";

export default function SuperAdminRoutes() {
  return (
    <Fragment>
      <Route path="/super-admin/dashboard" element={<RoleDashboardPage />} />
      <Route path="/super-admin/reports" element={<RoleReportsPage />} />

      <Route path="/super-admin/users" element={<UsersIndex />} />
      <Route path="/super-admin/departments" element={<DepartmentsIndex />} />
      <Route path="/super-admin/patients" element={<PatientsIndex />} />
    </Fragment>
  );
}
