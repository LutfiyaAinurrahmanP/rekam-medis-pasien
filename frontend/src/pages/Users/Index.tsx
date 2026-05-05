import { useState } from "react";
import { useNavigate, useParams } from "react-router";
import UserTable from "../../components/Users/UserTable";
import Button from "../../components/ui/button/Button";
import { User } from "../../hooks/Users/useUsers";
import PageBreadcrumb from "../../components/common/PageBreadCrumb";

export default function UsersIndex() {
  const navigate = useNavigate();
  const { role } = useParams();
  const [selectedUser, setSelectedUser] = useState<User | null>(null);

  const handleEdit = (user: User) => {
    setSelectedUser(user);
    navigate(`/${role}/users/${user.id}/edit`);
  };

  const handleDelete = (user: User) => {
    setSelectedUser(user);
    // TODO: Implement delete functionality with confirmation modal
    console.log("Delete user:", user);
  };

  return (
    <div className="space-y-4">
      {/* Header */}
     <PageBreadcrumb pageTitle="Users" />

      {/* User Table */}
      <UserTable onEdit={handleEdit} onDelete={handleDelete} />
    </div>
  );
}
