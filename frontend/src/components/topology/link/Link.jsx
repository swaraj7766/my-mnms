import React from "react";
import injectStyle from "../injectStyle";

/**
 * Link component is responsible for encapsulating link render.
 * @example
 * const onClickLink = function(source, target) {
 *      window.alert(`Clicked link between ${source} and ${target}`);
 * };
 *
 * const onRightClickLink = function(source, target) {
 *      window.alert(`Right clicked link between ${source} and ${target}`);
 * };
 *
 * const onMouseOverLink = function(source, target) {
 *      window.alert(`Mouse over in link between ${source} and ${target}`);
 * };
 *
 * const onMouseOutLink = function(source, target) {
 *      window.alert(`Mouse out link between ${source} and ${target}`);
 * };
 *
 * <Link
 *     d="M1..."
 *     source="idSourceNode"
 *     target="idTargetNode"
 *     markerId="marker-small"
 *     strokeWidth=1.5
 *     stroke="green"
 *     strokeDasharray="5 1"
 *     strokeDashoffset="3"
 *     strokeLinecap="round"
 *     className="link"
 *     opacity=1
 *     mouseCursor="pointer"
 *     onClickLink={onClickLink}
 *     onRightClickLink={onRightClickLink}
 *     onMouseOverLink={onMouseOverLink}
 *     onMouseOutLink={onMouseOutLink} />
 */
export default class Link extends React.Component {
  constructor(props) {
    super(props);

    this.selectorRef = React.createRef(null);
    const keyframesStyle = `
      @-webkit-keyframes dash {
        from {
          stroke-dashoffset: 40;
        }
      }
    `;
    injectStyle(keyframesStyle);
  }
  /**
   * Handle link click event.
   * @returns {undefined}
   */
  handleOnClickLink = () =>
    this.props.onClickLink &&
    this.props.onClickLink(this.props.source, this.props.target);

  /**
   * Handle link right click event.
   * @param {Object} event - native event.
   * @returns {undefined}
   */
  handleOnRightClickLink = (event) =>
    this.props.onRightClickLink &&
    this.props.onRightClickLink(event, this.props.source, this.props.target);

  /**
   * Handle mouse over link event.
   * @returns {undefined}
   */
  handleOnMouseOverLink = () =>
    this.props.onMouseOverLink &&
    this.props.onMouseOverLink(this.props.source, this.props.target);

  /**
   * Handle mouse out link event.
   * @returns {undefined}
   */
  handleOnMouseOutLink = () =>
    this.props.onMouseOutLink &&
    this.props.onMouseOutLink(this.props.source, this.props.target);

  linkTransform = () => {
    if (this.props.x2 < this.props.x1) {
      let rx = this.props.x2 + (this.props.x1 - this.props.x2) / 2;
      let ry = this.props.y1 + (this.props.y2 - this.props.y1) / 2;
      return "rotate(180 " + rx + " " + ry + ")";
    } else {
      return "rotate(0)";
    }
  };

  render() {
    const lineStyle = {
      strokeWidth: this.props.strokeWidth,
      stroke: this.props.stroke,
      opacity: this.props.opacity,
      fill: "none",
      cursor: this.props.mouseCursor,
      strokeDasharray: this.props.strokeDasharray,
      strokeDashoffset: this.props.strokeDasharray,
      strokeLinecap: this.props.strokeLinecap,
    };

    const linkFlowStyle = {
      strokeWidth: this.props.strokeWidth,
      stroke: "yellow",
      strokeDasharray: [1, 20],
      strokeLinecap: "square",
      animation: "dash 1.25s linear infinite alternate",
    };

    const lineProps = {
      className: this.props.className,
      d: this.props.d,
      onClick: this.handleOnClickLink,
      onContextMenu: this.handleOnRightClickLink,
      onMouseOut: this.handleOnMouseOutLink,
      onMouseOver: this.handleOnMouseOverLink,
      style: lineStyle,
    };

    if (this.props.markerId) {
      lineProps.markerEnd = `url(#${this.props.markerId})`;
    }

    const { label, id, linkFlow, x1, x2, y1, y2 } = this.props;
    const textProps = {
      dy: 0,
      dx: 0,

      style: {
        fill: this.props.fontColor,
        fontSize: this.props.fontSize,
        fontWeight: this.props.fontWeight,
      },
    };

    //console.log(this.props.x1);
    //console.log(this.props.x2);

    return (
      <g>
        <path {...lineProps} id={id} />
        {linkFlow && <path d={this.props.d} style={linkFlowStyle} />}
        {label && (
          <g
            transform={`translate(${x2 + (x1 - x2) / 2} ${y1 + (y2 - y1) / 2})`}
          >
            <text {...textProps} textAnchor="middle">
              {label}
            </text>
          </g>
        )}
        {/* {label && (
          <text
            {...textProps}
            transform={this.linkTransform()}
            textAnchor="middle"
          >
            <textPath href={`#${id}`} textAnchor="middle" startOffset="50%">
              {label}
            </textPath>
          </text>
        )} */}
      </g>
    );
  }
}
