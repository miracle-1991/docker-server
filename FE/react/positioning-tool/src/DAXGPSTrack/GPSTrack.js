import React from "react";
import Menu from "../Menu/menu";
import Search from "./Search";

export default function GPSTrack() {
    return (
        <>
            <Menu choose="driver"></Menu>
            <Search></Search>
        </>
    )
}