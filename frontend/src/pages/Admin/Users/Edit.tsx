import { useParams } from "react-router";
import EditUserModal from "../../../layout/Admin/Users/EditUserModal";

export default function UsersEdit() {
  const { role, id } = useParams();

  return <EditUserModal isOpen={true} id={id} role={role} />;
}
