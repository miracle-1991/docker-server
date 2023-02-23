import React from "react";

export default function StdoutItem({ProcessingState, data, name}) {
    if (ProcessingState === false) {
        return (
            <div className="panel">
                <div className="panel-heading">
                    {name}
                </div>
                <div className="panel-body" key="parselog">
                    <pre>{}</pre>
                </div>
            </div>
        )
    }else {
        return (
            <div className="panel">
                <div className="panel-heading">
                    {name}
                </div>
                <div className="panel-body" key="parselog">
                    <pre>{JSON.stringify(data, null, 2)}</pre>
                </div>
            </div>
        )
    }

}