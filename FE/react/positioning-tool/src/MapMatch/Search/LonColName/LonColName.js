import React, {useState} from "react";
import { useReq } from "../ReqProvider";

export default function LonColName() {
    const {req, onLongitudeColumnNameChange } = useReq();
    const [name, setName] = useState(req.longitude_column_name);
    return (
        <>
            <div className="row">
                <div className="col-xs-12">
                    <span className="label">NAME OF LON COLUMN</span>
                    <input type="text" className="form-control" value={name} onChange={
                        event => {
                            let newV = event.target.value;
                            setName(newV);
                            onLongitudeColumnNameChange(newV);
                        }
                    }/>
                </div>
            </div>
        </>
    )
}