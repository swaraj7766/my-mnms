import { Navigate } from "react-router-dom";
import { isExpired } from "react-jwt";

export default function PrivateRoute({ children }) {
  const token = sessionStorage.getItem("nmstoken");
  let istokenExpired = isExpired(token);
  return token && !istokenExpired ? (
    children
  ) : (
    <Navigate to="/login" replace={true} />
  );
}
