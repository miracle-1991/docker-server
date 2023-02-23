import React from "react";
import Search from "./Search/Search";
import Menu from "../Menu/menu";

export function RTK() {
    return(
        <>
            <Menu choose="rtk"></Menu>
            <Search></Search>
        </>
    )
}