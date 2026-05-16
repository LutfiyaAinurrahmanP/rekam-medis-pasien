import CreateDepartmentModal from "../../../layout/Admin/Departments/CreateDepartmentModal";
import { useNavigate } from "react-router";
import { useAuth } from "../../../context/AuthContext";
import { getRoleDepartmentsPath } from "../../../pages/Roles/shared/role-routing";

export default function DepartmentsCreate() {
  const navigate = useNavigate();
  const { user } = useAuth();
  return (
    <CreateDepartmentModal
      isOpen={true}
      onClose={() => navigate(getRoleDepartmentsPath(user?.role))}
    />
  );
}
