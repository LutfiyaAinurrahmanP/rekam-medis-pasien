import CreateDepartmentModal from "../../../layout/Admin/Departments/CreateDepartmentModal";
import { useNavigate } from "react-router";

export default function DepartmentsCreate() {
  const navigate = useNavigate();
  return (
    <CreateDepartmentModal
      isOpen={true}
      onClose={() => navigate("/admin/departments")}
    />
  );
}
