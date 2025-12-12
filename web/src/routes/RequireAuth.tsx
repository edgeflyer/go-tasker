import { Navigate, useLocation } from "react-router-dom";
import { isAuthed } from "../auth/auth";

export default function RequireAuth({
  children,
}: {
  children: React.ReactNode;
}) {
  const loc = useLocation();
  if (!isAuthed()) {
    return <Navigate to="/login" replace state={{ from: loc.pathname }} />;
  }
  return <>{children}</>;
}
