import { Fragment } from "react";
import { Route } from "react-router";
import RoleDashboardPage from "../pages/Roles/RoleDashboardPage";
import RoleReportsPage from "../pages/Roles/RoleReportsPage";
import UsersCreate from "../pages/Admin/Users/Create";
import UsersEdit from "../pages/Admin/Users/Edit";
import UsersIndex from "../pages/Admin/Users/Index";
import UsersShow from "../pages/Admin/Users/Show";

export default function AdminRoutes() {
  return (
    <Fragment>
      <Route path="/admin/dashboard" element={<RoleDashboardPage />} />
      <Route path="/admin/reports" element={<RoleReportsPage />} />
      <Route path="/admin/users" element={<UsersIndex />} />
      <Route path="/admin/users/create" element={<UsersCreate />} />
      <Route path="/admin/users/:id/edit" element={<UsersEdit />} />
      <Route path="/admin/users/:id" element={<UsersShow />} />
    </Fragment>
  );
}
