import { FileExcelOutlined, FilePdfOutlined } from "@ant-design/icons";
import { Button, Dropdown } from "antd";
import jsPDF from "jspdf";
import "jspdf-autotable";
import { CSVLink } from "react-csv";
import React, { useEffect, useRef, useState } from "react";
const items = [
  {
    label: "PDF",
    key: "pdf",
    icon: <FilePdfOutlined style={{ fontSize: 15 }} />,
  },
  {
    label: "CSV",
    key: "csv",
    icon: <FileExcelOutlined style={{ fontSize: 15 }} />,
  },
];
const ExportData = ({ Columns, DataSource, title }) => {
  //const [allowDownloadCsv, setAllowDownloadCsv] = useState(false);
  const csvLink = useRef();
  const [csvHeader, setCsvHeader] = useState([]);
  const [csvData, setCsvData] = useState([]);
  const handlePdfDownload = () => {
    const headers = [Columns.map((col) => col.title)];
    const data = DataSource.map((item) => {
      return Columns.map((col) => item[col.key]);
    });
    exportPDF(data, headers, title);
  };

  const handleCsvDownload = () => {
    const headers = Columns.map((col) => col.title);
    const data = DataSource.map((item) => {
      return Columns.map((col) => item[col.key]);
    });
    setCsvHeader((prev) => headers);
    setCsvData((prev) => data);
  };

  const handleMenuClick = (e) => {
    if (e.key === "pdf") {
      console.log("do pdf oeration");
      handlePdfDownload();
    } else {
      console.log("do csv operation");
      handleCsvDownload();
    }
  };

  useEffect(() => {
    if (csvData.length > 0) csvLink.current.link.click();
  }, [csvData]);

  return (
    <>
      <Dropdown
        menu={{ items, onClick: handleMenuClick }}
        placement="bottomLeft"
      >
        <Button type="primary">Export</Button>
      </Dropdown>
      {/* {allowDownloadCsv && (
        <CSVDownload
          headers={csvHeader}
          data={csvData}
          filename={getFilename(title)}
        />
      )} */}
      <CSVLink
        headers={csvHeader}
        data={csvData}
        filename={getFilename(title)}
        className="hidden"
        ref={csvLink}
        target="_blank"
      />
    </>
  );
};

export default ExportData;

const exportPDF = (data, headers, title) => {
  const unit = "pt";
  const size = "A4"; // Use A1, A2, A3 or A4
  const orientation = "landscape"; // portrait or landscape

  const marginLeft = 40;
  const doc = new jsPDF(orientation, unit, size);

  doc.setFontSize(15);

  let content = {
    startY: 50,
    head: headers,
    body: data,
  };

  doc.text(title, marginLeft, 40);
  doc.autoTable(content);
  doc.save(getFilename(title));
};

export const getFilename = (title) => {
  // For todays date;
  // eslint-disable-next-line no-extend-native
  Date.prototype.today = function () {
    return (
      (this.getDate() < 10 ? "0" : "") +
      this.getDate() +
      (this.getMonth() + 1 < 10 ? "0" : "") +
      (this.getMonth() + 1) +
      this.getFullYear()
    );
  };

  // For the time now
  // eslint-disable-next-line no-extend-native
  Date.prototype.timeNow = function () {
    return (
      (this.getHours() < 10 ? "0" : "") +
      this.getHours() +
      (this.getMinutes() < 10 ? "0" : "") +
      this.getMinutes() +
      (this.getSeconds() < 10 ? "0" : "") +
      this.getSeconds()
    );
  };
  const currentdate = new Date();

  return `${title}_${currentdate.today()}_${currentdate.timeNow()}`;
};
