import React from "react";
import { useReq } from "../ReqProvider";
import ObsNavPathItem from "./ObsNavPathItem";

export default function ObsNavPath() {
    const { req, onRTKRoverObsPathChange,onRTKStationObsPathChange, onRTKStationNavPathChange } = useReq()
    return (
        <>
            <ObsNavPathItem name="roverobs" initval={req.rtk.obsnav.roverobs} onValChange={onRTKRoverObsPathChange}></ObsNavPathItem>
            <ObsNavPathItem name="roverobs" initval={req.rtk.obsnav.stationobs} onValChange={onRTKStationObsPathChange}></ObsNavPathItem>
            <ObsNavPathItem name="roverobs" initval={req.rtk.obsnav.stationnav} onValChange={onRTKStationNavPathChange}></ObsNavPathItem>
        </>
    )
}