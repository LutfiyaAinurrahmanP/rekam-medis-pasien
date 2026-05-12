import { useParams } from "react-router";
import ShowUserModal from "../../../layout/Admin/Users/ShowUserModal";

export default function UsersShow() {
  const { role, id } = useParams();

  return <ShowUserModal isOpen={true} id={id} role={role} />;
}
