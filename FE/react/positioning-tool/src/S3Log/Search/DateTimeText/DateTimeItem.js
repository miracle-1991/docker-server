import React, {useState} from "react";
import { useReq } from "../ReqProvider";

export default function DateTimeItem({titlestr="", timetype="start"}) {
    const {req, onStartTimeChange, onEndTimeChange} = useReq();
    const [tstr, setTstr] = useState(req.starttime);
    return (
        <div className="col-xs-6">
            <span className="label">{titlestr}</span>
            <input type="datetime-local" step="1" className="form-control" value={tstr} onChange={
                event => {
                    let newT = event.target.value.replace("T", " ");
                    setTstr(newT);
                    timetype === "start" ? onStartTimeChange(newT) : onEndTimeChange(newT);
                }
            }/>
        </div>
    )
}