import React, {useState} from "react";

export default function ObsNavPathItem({name, initval, onValChange}) {
    const [v, setV] = useState(initval)
    return (
        <>
            <div className="row">
                <div className="col-xs-12">
                    <span className="label">{name}</span>
                    <input type="text" className="form-control" value={v} onChange={
                        event => {
                            let newV = event.target.value;
                            setV(newV);
                            onValChange(newV);
                        }
                    }/>
                </div>
            </div>
            <br></br>
        </>
    )
}