import { useParams } from "react-router";
import { useNavigate } from "react-router";
import ShowPatientModal from "../../../layout/Admin/Patients/ShowPatientModal";
import { useAuth } from "../../../context/AuthContext";
import { getRolePatientsPath } from "../../Roles/shared/role-routing";

export default function PatientsShow() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();

  return (
    <ShowPatientModal
      isOpen={true}
      id={id}
      onClose={() => navigate(getRolePatientsPath(user?.role))}
    />
  );
}
