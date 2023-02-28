import EHG7512 from "../../assets/images/EHG7512.png";
import EHG2408 from "../../assets/images/EHG2408.png";

const TopologyImage = (model) => {
  switch (model) {
    case "EHG7512":
      return window.location.origin + EHG7512;
    case "EHG2408":
      return window.location.origin + EHG2408;
    default:
      return window.location.origin + EHG2408;
  }
};
export { TopologyImage };
