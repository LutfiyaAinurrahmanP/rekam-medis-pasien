import EditDepartmentModal from "../../../layout/Admin/Departments/EditDepartmentModal";
import { useParams } from "react-router";

export default function DepartmentsEdit() {
  const { id } = useParams();
  return <EditDepartmentModal isOpen={true} id={id} />;
}
