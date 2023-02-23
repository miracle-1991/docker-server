import React from "react";
import {useReq} from "../ReqProvider";
import StdoutItem from "./StdoutItem";

export default function Stdout() {
    const {rtkData,rtkProcessingState} = useReq();
    return (
        <>
            <StdoutItem key="rtk" ProcessingState={rtkProcessingState} data={rtkData} name="RTK POSITIONING RESULT"></StdoutItem>
        </>
    )
}