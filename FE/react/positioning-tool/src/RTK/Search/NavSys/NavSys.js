import React from "react";
import {useReq} from "../ReqProvider";
import NavSysItem from "./NavSysItem";

export default function NavSys() {
    const {navsyslist, navsysCheckState, onNavSysCheckChange} = useReq();
    return(
        <>
            <span className="label">NAV SYSTEM</span>
            <div>
                {
                    navsyslist.map(({name}, index) => {
                        return (
                            <NavSysItem key={index} name={name} index={index} checked={navsysCheckState} onchange={onNavSysCheckChange}></NavSysItem>
                        )
                    })
                }
            </div>
        </>
    )
}