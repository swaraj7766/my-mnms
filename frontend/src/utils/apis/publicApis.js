import axios from "axios";

const baseURL =
  localStorage.getItem("nms-base-URL") === null
    ? "http://localhost:27182"
    : `${JSON.parse(localStorage.getItem("nms-base-URL"))}`;

export default axios.create({
  baseURL,
  headers: {
    accept: "application/json",
    "Content-Type": "application/json",
  },
});
