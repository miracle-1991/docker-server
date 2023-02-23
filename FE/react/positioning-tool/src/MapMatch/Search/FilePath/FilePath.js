import React, {useState} from "react";
import { useReq } from "../ReqProvider";

export default function FilePath() {
    const {req, onFilepathChange } = useReq();
    const [dir, setDir] = useState(req.filepath);
    return (
        <>
            <div className="row">
                <div className="col-xs-12">
                    <span className="label">CSV FILE PATH</span>
                    <input type="text" className="form-control" value={dir} onChange={
                        event => {
                            let newV = event.target.value;
                            setDir(newV);
                            onFilepathChange(newV);
                        }
                    }/>
                </div>
            </div>
        </>
    )
}