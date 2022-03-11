import React, { useState } from "react";
import { MoreVert } from "@mui/material";

import "./Table.css";

function ContextMenuButton({ menuComponent, menuProps, data, invisible }) {
  const [anchorEl, setAnchorEl] = useState(null);
  const handleClick = (event) => setAnchorEl(event.currentTarget);
  const handleClose = () => setAnchorEl(null);
  const MenuComponent = menuComponent;

  const button = !invisible ? (
    <div className="Table-row-menu-button" onClick={handleClick} aria-controls="simple-menu" aria-haspopup="true">
      <MoreVert fontSize="small" />
    </div>
  ) : (
    <div className="Table-row-menu-button-invisible" ria-controls="simple-menu" aria-haspopup="true">
      <MoreVert fontSize="small" />
    </div>
  );

  return (
    <React.Fragment>
      {!invisible ? (
        <div className="Table-row-menu-button" onClick={handleClick} aria-controls="simple-menu" aria-haspopup="true">
          <MoreVert fontSize="small" />
        </div>
      ) : (
        <div className="Table-row-menu-button-invisible" ria-controls="simple-menu" aria-haspopup="true">
          <MoreVert fontSize="small" />
        </div>
      )}
      {!invisible && (
        <MenuComponent
          data={data}
          anchorEl={anchorEl}
          keepMounted
          open={!!anchorEl}
          onClose={handleClose}
          {...menuProps}
        />
      )}
    </React.Fragment>
  );
}

function Row({ row, columns, isHeader, menuComponent, menuProps }) {
  const cells = columns.map((c, i) => {
    const { field, name, flex } = c;
    return (
      <div className={isHeader ? "Table-header-cell" : "Table-row-cell"} style={{ flex }} key={i}>
        {isHeader ? name : row[field]}
      </div>
    );
  });

  return (
    <div className={isHeader ? "Table-header" : "Table-row"}>
      {cells}
      {!!menuComponent && (
        <ContextMenuButton menuComponent={menuComponent} menuProps={menuProps} data={row} invisible={!!isHeader} />
      )}
    </div>
  );
}

export default function Table({ rows = [], columns, className, menuComponent, menuProps }) {
  const rowComponents = rows.map((r, i) => {
    return <Row row={r} columns={columns} menuComponent={menuComponent} menuProps={menuProps} key={i} />;
  });

  return (
    <div className={`Table ${className ? className : ""}`}>
      <Row isHeader columns={columns} menuComponent={menuComponent} menuProps={menuProps} />
      {rowComponents}
    </div>
  );
}
