import React from "react";
import {useReq} from "../ReqProvider";
import PmodeItem from "./PmodeItem";

export default function Pmode() {
    const {pmodelist, pmodeCheckState, onPmodeCheckChange} = useReq();
    return (
        <>
            <span className="label">PMODE FOR RTK</span>
            <div>
                {
                    pmodelist.map(({name}, index) => {
                        return (
                            <PmodeItem key={index} name={name} index={index} checked={pmodeCheckState} onchange={onPmodeCheckChange}></PmodeItem>
                        )
                    })
                }
            </div>
        </>
    )
}