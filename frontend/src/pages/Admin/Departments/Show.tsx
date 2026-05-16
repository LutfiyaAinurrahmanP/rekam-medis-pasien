import ShowDepartmentModal from "../../../layout/Admin/Departments/ShowDepartmentModal";
import { useParams } from "react-router";

export default function DepartmentsShow() {
  const { id } = useParams();
  return <ShowDepartmentModal isOpen={true} id={id} />;
}
