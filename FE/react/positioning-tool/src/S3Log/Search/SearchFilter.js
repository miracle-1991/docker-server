import React, {useState} from "react";
import {useReq} from "./ReqProvider";

export default function SearchFilter() {
    const {req, onfilterChange} = useReq();
    const [filterlistStr, setFilterListStr] = useState(req.filter.join(","))
    return (
        <input type="text" className="form-control" value={filterlistStr} onChange={
            event => {
                let newV = event.target.value;
                setFilterListStr(newV);
                onfilterChange(newV)
            }
        }/>
    )
}