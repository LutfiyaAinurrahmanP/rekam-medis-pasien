import { Fragment } from "react";
import { Route } from "react-router";
import RoleDashboardPage from "../pages/Roles/RoleDashboardPage";
import RoleReportsPage from "../pages/Roles/RoleReportsPage";
import UsersCreate from "../pages/Admin/Users/Create";
import UsersEdit from "../pages/Admin/Users/Edit";
import UsersIndex from "../pages/Admin/Users/Index";
import UsersShow from "../pages/Admin/Users/Show";
import DepartmentsIndex from "../pages/Admin/Departments/Index";
import DepartmentsCreate from "../pages/Admin/Departments/Create";
import DepartmentsEdit from "../pages/Admin/Departments/Edit";
import DepartmentsShow from "../pages/Admin/Departments/Show";

export default function AdminRoutes() {
  return (
    <Fragment>
      <Route path="/admin/dashboard" element={<RoleDashboardPage />} />
      <Route path="/admin/reports" element={<RoleReportsPage />} />
      <Route path="/admin/users" element={<UsersIndex />} />
      <Route path="/admin/users/create" element={<UsersCreate />} />
      <Route path="/admin/users/:id/edit" element={<UsersEdit />} />
      <Route path="/admin/users/:id" element={<UsersShow />} />
      <Route path="/admin/departments" element={<DepartmentsIndex />} />
      <Route path="/admin/departments/create" element={<DepartmentsCreate />} />
      <Route path="/admin/departments/:id/edit" element={<DepartmentsEdit />} />
      <Route path="/admin/departments/:id" element={<DepartmentsShow />} />
    </Fragment>
  );
}
