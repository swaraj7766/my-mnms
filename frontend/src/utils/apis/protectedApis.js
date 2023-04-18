import axios from "axios";

const baseURL =
  localStorage.getItem("nms-base-URL") === null
    ? "http://localhost:27182"
    : `${JSON.parse(localStorage.getItem("nms-base-URL"))}`;

const instance = axios.create({
  baseURL: baseURL,
});

instance.defaults.headers.common[
  "Authorization"
] = `Bearer ${sessionStorage.getItem("nmstoken")}`;

export default instance;
