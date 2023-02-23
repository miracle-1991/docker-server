import React, {useState} from "react";
import { useReq } from "../ReqProvider";

export default function FileName() {
    const {req, onFilenameChange } = useReq();
    const [file, setFile] = useState(req.filename);
    return (
        <>
            <div className="row">
                <div className="col-xs-12">
                    <span className="label">CSV FILE NAME</span>
                    <input type="text" className="form-control" value={file} onChange={
                        event => {
                            let newV = event.target.value;
                            setFile(newV);
                            onFilenameChange(newV);
                        }
                    }/>
                </div>
            </div>
        </>
    )
}