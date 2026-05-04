import PageMeta from "../../components/common/PageMeta";
import AuthLayout from "./AuthPageLayout";
import LoginForm from "../../components/auth/LoginForm";

export default function Login() {
  return (
    <>
      <PageMeta
        title="Login | Rekam Medis Pasien"
        description="Login page to access the dashboard according to user roles."
      />
      <AuthLayout>
        <LoginForm />
      </AuthLayout>
    </>
  );
}
