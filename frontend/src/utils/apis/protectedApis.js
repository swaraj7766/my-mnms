import axios from "axios";

const baseURL =
  process.env.NODE_ENV === "production"
    ? process.env.REACT_APP_BASE_URL
    : process.env.REACT_APP_BASE_URL;

const instance = axios.create({
  baseURL: baseURL,
});

instance.defaults.headers.common[
  "Authorization"
] = `Bearer ${sessionStorage.getItem("nmstoken")}`;

export default instance;
